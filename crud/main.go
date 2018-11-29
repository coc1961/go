package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/coc1961/go/crud/entity"
)

func main() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	var e entity.Entity
	err := e.Load(dir, "prueba")
	fmt.Println(err.Error())

	fmt.Println(e)
}
