package crudframework

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/coc1961/go/crud/entities"
	"github.com/coc1961/go/jsonutil"
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
		],
		"arrayOnject": [
		  {
			"id": 123,
			"name": "Nombre1"
		  },
		  {
			"id": 456,
			"name": "Nombre2"
		  }
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

	ent1 := jsonutil.New()

	ent1.Add("Prueba").Set("Hola Men")

	ent.Add("Nuevo").Set(ent1)

	fmt.Println(ent.JSON())

	ent2 := jsonutil.New()
	ent2.Add("id").Set("788")
	ent2.Add("name").Set("Nombre2")

	fmt.Println(ent2.Get("id").Value())

	fmt.Println(ent2.Get("id").Get("NoEsta").Get("NoEsta2").Value())

	tmp := ent.Get("arrayOnject").ValueAsArray()
	tmp = append(tmp, ent2.GetRoot())
	ent.Get("arrayOnject").Set(tmp)
	fmt.Println(ent.JSON())
	tmp = ent.Get("arrayOnject").ValueAsArray()

	fmt.Println(ent.Get("hijo").Get("idAtt").Value())

	fmt.Println(ent.Get("soyArray").Value())

	fmt.Println(ent.Get("hijo").Get("idAtt").Value())

	ent.Get("hijo").Get("idAtt").Set("Prueba")

	fmt.Println(ent.Get("hijo").Get("idAtt").Value())

	ent.Get("hijo").Get("idAtt").Add("Otro").Set("Prueba1")

	fmt.Println(ent.Get("hijo").Get("idAtt").Get("Otro").Value())

	ent.Get("hijo").Get("idAtt").Get("Otro").Set("Prueba2")

	fmt.Println(ent.Get("hijo").Get("idAtt").Get("Otro").Value())

	fmt.Println(ent.JSON())

	fmt.Println("=============================================================")

	fmt.Println(ent.Get("hijo").Get("idAtt").Value())

	ent.Get("hijo").Get("idAtt").Set("Prueba")

	fmt.Println(ent.Get("hijo").Get("idAtt").Value())

	ent.Get("hijo").Get("idAtt").Add("Otro").Set("Prueba1")

	fmt.Println(ent.Get("hijo").Get("idAtt").Get("Otro").Value())

	ent.Get("hijo").Get("idAtt").Get("Otro").Set("Prueba2")

	fmt.Println(ent.Get("hijo").Get("idAtt").Get("Otro").Value())

	fmt.Println(ent.JSON())

	ent.Get("name").Set("Soy Nombre")

	fmt.Println(ent.Get("name").Value())

	fmt.Println(ent.JSON())

}

// CrudFramework crud framework
type CrudFramework struct {
	configPath  string
	definitions map[string]entities.Definition
}

// New new Crud Framework
func New(configPath string) *CrudFramework {
	return &CrudFramework{configPath, make(map[string]entities.Definition, 0)}
}

// Load Config Files
func (e *CrudFramework) Load() error {
	path := filepath.Join(e.configPath, "data", "")
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		def := entities.NewEntityDefinition()
		err := def.Load(filepath.Join(e.configPath, "data"), f.Name())
		if err != nil {
			return err
		}
		e.definitions[def.Name()] = *def
	}
	return nil
}
