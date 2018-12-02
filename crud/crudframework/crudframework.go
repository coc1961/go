package crudframework

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/coc1961/go/crud/entities"
)

// Test Prueba
func Test() {
	sjson := `
		{
			"id": "xy23",
			"name": "Nombre",
			"amount": 100.45,
			"age": 55,
			"creationDate": "2018-11-24T01:10:22Z",
			"hijo": {
				"idAtt": "cc23",
                "nameAtt": "AttrHijo"
			},
			"soyArray": [
				"elem1",
				"elem2"
			]
		}
			`
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	var e entities.Definition

	dir = "/home/carlos/gopath/src/github.com/coc1961/go/crud"
	err := e.Load(dir, "/data/prueba")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = e.Validate(sjson)
	if err != nil {
		fmt.Println(err)
		return
	}

	var ent *entities.Entity
	ent, err = e.New(sjson)

	fmt.Println(ent.Get("hijo").Get("idAtt").Value())

	ent.Get("hijo").Get("idAtt").Set("Prueba")

	fmt.Println(ent.Get("hijo").Get("idAtt").Value())

	ent.Get("hijo").Get("idAtt").Add("Otro").Set("Prueba1")

	fmt.Println(ent.Get("hijo").Get("idAtt").Get("Otro").Value())

	ent.Get("hijo").Get("idAtt").Get("Otro").Set("Prueba2")

	fmt.Println(ent.Get("hijo").Get("idAtt").Get("Otro").Value())

	fmt.Println(ent.JSON())
}
