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
	var e entity.Definition
	err := e.Load(dir, "prueba")
	if err != nil {
		fmt.Println(err.Error())
	}

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
			}
		}
			`
	var ojson map[string]interface{}
	err = json.Unmarshal([]byte(sjson), &ojson)

	var ent *entity.Entity
	ent, err = e.Parse(ojson)

	fmt.Println("******* " + ent.Name + " *******")

	for _, e := range ent.Atributes {
		print(e, "")
	}
	fmt.Println("=======================================")
}

func print(e *entity.Attribute, space string) {
	fmt.Println(space + "=======================================")
	fmt.Println(space + e.FieldDefinition.Name)
	fmt.Println(space + e.FieldDefinition.Type)
	fmt.Print(space)
	fmt.Println(e.Value)
	if e.Child != nil {
		for _, e1 := range e.Child {
			print(e1, space+"\t")
		}
	}
}
