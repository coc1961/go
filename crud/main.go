package main

import (
	"os"
	"path/filepath"

	"github.com/coc1961/go/crud/crudframework"
)

func main() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir = "/home/carlos/gopath/src/github.com/coc1961/go/crud"

	cf := crudframework.New(dir)
	cf.Load()

}
