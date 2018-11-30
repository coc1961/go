package crudframework

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/coc1961/go/crud/entity"
)

// Test Prueba
func Test() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	var e entity.Entity
	err := e.Load(dir, "prueba")
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(e)

	sjson := `{
			"id": "123",
			"name": "Nombre",
			"amount": 100.45,
			"age": 55,
			"creationDate": "2018-11-24T01:10:22"
			}
			`
	var ojson map[string]interface{}
	err = json.Unmarshal([]byte(sjson), &ojson)
	err = e.Validate(ojson)

	fmt.Println(err)
}
