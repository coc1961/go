package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"
)

// BUG(carlos): Manage Errors!
func main() {
	if len(os.Args) != 4 {
		log.Printf("usage: %s [concurrency] [url] [output]", os.Args[0])
		return
	}

	workers, err := strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		log.Fatalf("error parsing concurrency param: %v", err)
	}

	resourceURL, err := url.Parse(os.Args[2])
	if err != nil {
		log.Fatalf("error parsing url param: %v", err)
	}

	os.Remove(os.Args[3])
	out, err := os.OpenFile(os.Args[3], os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("error opening file for writing: %v", err)
	}

	res, err := http.Head(resourceURL.String())
	if err != nil {
		log.Fatalf("error requesting HEAD of file: %v", err)
	}

	if res.StatusCode >= 300 {
		log.Fatalf("unexpected status code received from server: %s", res.Status)
	}

	acceptRange := res.Header.Get("accept-ranges")
	if acceptRange != "bytes" && workers > 1 {
		log.Fatalf("remote server does not accept range downloads")
	}

	contentLength, err := strconv.ParseInt(res.Header.Get("content-length"), 10, 64)
	if err != nil || contentLength == 0 {
		log.Fatalf("remote server content-length is invalid")
	}

	pba = make([]*progressReader, 0)

	out.Truncate(contentLength)

	chunkSize := contentLength / workers

	wg := sync.WaitGroup{}
	wg.Add(int(workers))

	pd := make([]*partialDownload, 0)

	for i := int64(0); i < workers; i++ {
		chunkStart, chunkEnd := chunkSize*i, (chunkSize*i)+chunkSize-1
		if i+1 == workers {
			chunkEnd = contentLength
		}

		tmp := createPartialDownload(resourceURL, &wg, chunkStart, chunkEnd, out)
		pd = append(pd, tmp)
		go tmp.Download()
	}

	go func() {
		for {
			fmt.Print("\rProgress [ ")
			for _, v := range pba {
				fmt.Print(v.pos / (v.len / 100))
				fmt.Print("% ")
			}
			fmt.Print("]")
			time.Sleep(time.Millisecond * 10)
		}
	}()

	wg.Wait()

	out.Close()

	for _, v := range pd {
		if v.out != nil {
			v.out.Close()
		}
	}
	for _, v := range pd {
		if v.err != nil {
			go os.Remove(out.Name())
			log.Fatal(v.err)
		}
	}
	time.Sleep(time.Millisecond * 10)
	fmt.Println("\nEnd...")

}

type partialDownload struct {
	resourceURL *url.URL
	rangeHeader string
	wg          *sync.WaitGroup
	out         *os.File
	err         error
	len         int64
	pos         int64
}

func createPartialDownload(resourceURL *url.URL, wg *sync.WaitGroup, chunkStart int64, chunkEnd int64, out *os.File) *partialDownload {
	rangeHeader := fmt.Sprintf("bytes=%d-%d", chunkStart, chunkEnd)
	return &partialDownload{resourceURL, rangeHeader, wg, out, nil, chunkEnd - chunkStart, chunkStart}
}

func (p *partialDownload) Download() {
	defer p.wg.Done()
	req, error := http.NewRequest("GET", p.resourceURL.String(), nil)
	if error != nil {
		p.err = error
		return
	}
	req.Header.Add("Range", p.rangeHeader)
	var client http.Client
	resp, error := client.Do(req)
	if error != nil {
		p.err = error
		return
	}

	out, err := os.OpenFile(p.out.Name(), os.O_RDWR, os.ModePerm)
	if err != nil {
		p.err = err
		return
	}
	p.out = out

	p.out.Seek(p.pos, os.SEEK_SET)

	wrapReader := createProgressReader(&resp.Body, p.len)

	pba = append(pba, wrapReader)
	_, error = io.Copy(p.out, wrapReader)
	if error != nil {
		p.err = error
	}
}

var pba []*progressReader

func createProgressReader(reader *io.ReadCloser, len int64) *progressReader {
	return &progressReader{reader, len, 0}
}

type progressReader struct {
	reader *io.ReadCloser
	len    int64
	pos    int64
}

func (r *progressReader) Read(p []byte) (n int, err error) {
	rr := *(r.reader)
	lei, err := rr.Read(p)
	r.pos += int64(lei)
	return lei, err
}
