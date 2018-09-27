package lib

// Download File, permite descargar un archivo pedazos con varios hilos simultaneos
//
// Utilzando la funcion
//
// func DownloadFile(resourceURL *url.URL, workers int64, out *os.File, listener func(status []*ProgressReader)) {
//
// recibe la url, la cantidad de hilos, el archhivo de salida, y una funcion que recibe el status con el progreso
// de la descarga

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

// Estructura de descarga parcial
type partialDownload struct {
	resourceURL *url.URL // Url de descarga
	rangeHeader string   // bytes a descargar
	out         *os.File // archivo de salida
	err         error    // error ?
	chunkSize   int64    // bytes a descargar
	chunkStart  int64    // posicion de inicio
}

// Creo partialDownload
func createPartialDownload(resourceURL *url.URL, chunkStart int64, chunkEnd int64, out *os.File) *partialDownload {
	rangeHeader := fmt.Sprintf("bytes=%d-%d", chunkStart, chunkEnd)
	return &partialDownload{resourceURL, rangeHeader, out, nil, chunkEnd - chunkStart, chunkStart}
}

// Descarga Parcial
func (p *partialDownload) Download(progressArray *[]*ProgressReader, wg *sync.WaitGroup) {
	defer wg.Done()
	// Request
	req, error := http.NewRequest("GET", p.resourceURL.String(), nil)
	if error != nil {
		p.err = error
		return
	}
	// Seteo rango de descarga
	req.Header.Add("Range", p.rangeHeader)
	var client http.Client
	resp, error := client.Do(req)
	if error != nil {
		p.err = error
		return
	}

	// Archivo de descarga parcial
	out, err := os.OpenFile(p.out.Name(), os.O_RDWR, os.ModePerm)
	if err != nil {
		p.err = err
		return
	}
	p.out = out

	// Me posiciono en la posicion en donde debo descargar
	p.out.Seek(p.chunkStart, os.SEEK_SET)

	// Creo un emboltorio de reader para que setee valores
	// de progreso de descarga para alimentar la barra de %
	wrapReader := createProgressReader(&resp.Body, p.chunkSize)

	// Agrego al array para desplegar el porc de descarga
	*progressArray = append(*progressArray, wrapReader)

	// Descargo!
	_, error = io.Copy(p.out, wrapReader)
	if error != nil {
		p.err = error
	}
	return
}

// Creo el objeto para alimentar la barra de progreso
func createProgressReader(reader *io.ReadCloser, len int64) *ProgressReader {
	return &ProgressReader{reader, len, 0}
}

//ProgressReader  Enboltorio de reader que guarda el % de descarga
type ProgressReader struct {
	reader *io.ReadCloser // Reader original
	len    int64          // total a descargar
	pos    int64          // bytes procesados
}

//Progress Retorna al porcentaje de la descarga realizada
func (r *ProgressReader) Progress() int64 {
	return int64(r.pos / (r.len / 100))
}

// Actualizo % de descarga y realizo la lectura real
func (r *ProgressReader) Read(p []byte) (n int, err error) {
	rr := *(r.reader)
	lei, err := rr.Read(p)
	r.pos += int64(lei)
	return lei, err
}

// DownloadFile descarga un archivo en partes procesadas en procesos simultaneos
func DownloadFile(resourceURL *url.URL, workers int64, out *os.File, listener func(status []*ProgressReader)) {

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

	// Inicializo Variables
	progressBarArray := make([]*ProgressReader, 0)
	partialDownloadArray := make([]*partialDownload, 0)

	// Reservo Espacio en el Archivo de Salida
	out.Truncate(contentLength)

	// Calculo el tamaño del chunk tamaño / hilos
	chunkSize := contentLength / workers

	wg := sync.WaitGroup{}
	wg.Add(int(workers))

	// Creo los downloaders parciales
	for i := int64(0); i < workers; i++ {
		chunkStart, chunkEnd := chunkSize*i, (chunkSize*i)+chunkSize-1
		if i+1 == workers {
			chunkEnd = contentLength
		}

		tmp := createPartialDownload(resourceURL, chunkStart, chunkEnd, out)
		partialDownloadArray = append(partialDownloadArray, tmp)

		go func() {
			// Comienzo Descarga
			tmp.Download(&progressBarArray, &wg)
		}()
	}

	// Muestro Barra de Progreso
	if listener != nil {
		go func() {
			for {
				listener(progressBarArray)
				time.Sleep(time.Millisecond * 10)
			}
		}()
	}

	wg.Wait()

	// Cierro archivo de salida
	out.Close()

	// Cierro archivos parciales
	for _, v := range partialDownloadArray {
		if v.out != nil {
			v.out.Close()
		}
	}
	// Hay Errores ?
	for _, v := range partialDownloadArray {
		if v.err != nil {
			go os.Remove(out.Name())
			log.Fatal(v.err)
		}
	}
	time.Sleep(time.Millisecond * 10)

}
