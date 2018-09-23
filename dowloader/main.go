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
	out, err := os.OpenFile(os.Args[3], os.O_RDWR|os.O_CREATE|os.O_EXCL, os.ModePerm)
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

	fmt.Println(contentLength)

	chunkSize := contentLength / workers

	wg := sync.WaitGroup{}
	wg.Add(int(workers))

	pd := make([]*partialDownload, 0)

	for i := int64(0); i < workers; i++ {
		chunkStart, chunkEnd := chunkSize*i, (chunkSize*i)+chunkSize-1
		if i+1 == workers {
			chunkEnd = contentLength
		}

		rangeHeader := fmt.Sprintf("bytes=%d-%d", chunkStart, chunkEnd)
		_ = rangeHeader

		fmt.Println(rangeHeader)
		tmp := createPartialDownload(resourceURL, rangeHeader, i, &wg)
		pd = append(pd, tmp)
		go tmp.Download()
	}

	_ = workers
	_ = resourceURL
	_ = out
	_ = chunkSize

	wg.Wait()

	for _, v := range pd {
		io.Copy(out, v.out)
		defer func(f *os.File) {
			f.Close()
			os.Remove(f.Name())
		}(v.out)
	}
}

type partialDownload struct {
	resourceURL *url.URL
	rangeHeader string
	i           int64
	wg          *sync.WaitGroup
	out         *os.File
}

func createPartialDownload(resourceURL *url.URL, rangeHeader string, i int64, wg *sync.WaitGroup) *partialDownload {
	return &partialDownload{resourceURL, rangeHeader, i, wg, nil}
}

func (p *partialDownload) Download() {
	defer p.wg.Done()
	req, _ := http.NewRequest("GET", p.resourceURL.String(), nil)
	req.Header.Add("Range", p.rangeHeader)
	var client http.Client
	resp, _ := client.Do(req)
	fileName := fmt.Sprintf("%d-output.bin", p.i)
	os.Remove(fileName)
	out, _ := os.Create(fileName)
	p.out = out
	defer func(fileName string) {
		out.Close()
		p.out, _ = os.Open(fileName)
	}(fileName)

	_, _ = io.Copy(out, resp.Body)
}
