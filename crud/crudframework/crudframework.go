package crudframework

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/coc1961/go/crud/entity"

	"github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-jsschema/validator"
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
			},
			"soyArray": [
				"elem1",
				"elem2"
			]
		}
			`
	var ojson map[string]interface{}
	err = json.Unmarshal([]byte(sjson), &ojson)

	s, err := schema.ReadFile(dir + "/data/pruebaschema.json")
	if err != nil {
		log.Printf("failed to read schema: %s", err)
		return
	}

	for name, pdef := range s.Properties {
		// Do what you will with `pdef`, which contain
		// Schema information for `name` property
		fmt.Println(name)
		fmt.Println(pdef)
	}

	// You can also validate an arbitrary piece of data

	v := validator.New(s)
	if err := v.Validate(ojson); err != nil {
		log.Printf("failed to validate data: %s", err)
	}

	fmt.Println(parse(ojson, "/hijo/idAtt"))
	fmt.Println(parse(ojson, "/soyArray"))

	/*
		var ent *entity.Entity
		ent, err = e.Parse(ojson)

		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("******* " + ent.Name + " *******")

		for _, e := range ent.Atributes {
			print(e, "")
		}
		fmt.Println("=======================================")
	*/
}

func parse(json map[string]interface{}, path string) interface{} {
	pt := strings.Split(path, "/")
	var ok bool
	var tmp map[string]interface{}
	var ret interface{}
	for _, p := range pt {
		if p == "" {
			continue
		}
		ret = nil

		tmp, ok = json[p].(map[string]interface{})
		if !ok {
			ret, ok = json[p].(interface{})
		} else {
			json = tmp
			ret = tmp
		}
	}
	return ret
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
