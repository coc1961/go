package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/coc1961/go/downloader/lib"
)

// BUG(carlos): Manage Errors!
func main() {

	// Inicio de operacion
	start := time.Now()

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

	lib.DownloadFile(resourceURL, workers, out, progress)

	//End!
	elapsed := time.Since(start)
	p := message.NewPrinter(language.English)
	p.Printf("\nProcess %d Bytes in %d seconds\n", contentLength, int(elapsed.Seconds()))
}

func progress(status []*lib.ProgressReader) {
	fmt.Print("\rProgress [ ")
	for _, v := range status {
		fmt.Print(v.Progress())
		fmt.Print("% ")
	}
	fmt.Print("]")

}

// fmt.Print("\rProgress [ ")
// for _, v := range progressBarArray {
// 	fmt.Print(v.pos / (v.len / 100))
// 	fmt.Print("% ")
// }
// fmt.Print("]")
