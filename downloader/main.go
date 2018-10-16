package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/coc1961/go/downloader/download"
)

// Gestor de Descarga Descarga en paralelo de archivo
// Soporta hasta n hilos, pasados por parametro

func main() {
	var pointerVerbose = flag.Bool("v", false, "show progress (Optional)")
	var pointerWorkers = flag.Int64("n", 2, "number of concurent downloads (Optional)")
	var pointerSUrl = flag.String("url", "", "download file url (required)")
	var pointerOutputFile = flag.String("o", "", "output file (required)")

	flag.Parse()

	if *pointerWorkers < 1 || *pointerSUrl == "" || *pointerOutputFile == "" {
		fmt.Fprintf(flag.CommandLine.Output(), "\n\nDownload Manager (%s) \n\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
		return
	}

	verbose := *pointerVerbose
	workers := *pointerWorkers
	surl := *pointerSUrl
	outputFile := *pointerOutputFile

	resourceURL, err := url.Parse(surl)
	if err != nil {
		log.Fatalf("error parsing url param: %v", err)
	}

	// Inicio de operacion
	start := time.Now()

	os.Remove(outputFile)
	out, err := os.OpenFile(outputFile, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalf("error opening file for writing: %v", err)
	}
	defer out.Close()

	client := download.CreateClient()
	res, err := client.Head(resourceURL.String())
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

	if verbose {
		download.File(resourceURL, workers, out, progress)
	} else {
		download.File(resourceURL, workers, out, nil)
	}

	//End!
	elapsed := time.Since(start)
	p := message.NewPrinter(language.English)
	if verbose {
		p.Printf("\nProcess %d Bytes in %d seconds\n", contentLength, int(elapsed.Seconds()))
	}
}

// Funcion que muestra el de Progreso de la descarga
func progress(status []download.Status) {
	fmt.Print("\rProgress [ ")
	for _, v := range status {
		fmt.Print(v.Progress())
		fmt.Print("% ")
	}
	fmt.Print("]")

}
