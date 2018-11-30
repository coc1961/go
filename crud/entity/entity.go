package entity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Entity Representa una Entidad a administrar con el crud
type Entity struct {
	Name        string
	Description string
	Fields      map[string]*Field
}

// Field Representa un atributo de una Entity
type Field struct {
	Name     string
	Type     string
	ID       bool
	Optional bool
	Default  interface{}
}

// Error Error en entidad
type Error struct {
	message string
	reason  interface{}
}

func (p Error) Error() string {
	if err, ok := p.reason.(error); ok {
		return fmt.Sprintf("%s => %v", p.message, err.Error())
	}
	return fmt.Sprintf("%s => %v", p.message, p.reason)
}

// Creo un objeto Field en base al json
func newField(name string, fld map[string]interface{}) *Field {
	var tmp interface{}
	var f Field

	// Valores Default
	f.Name = name
	f.ID = false
	f.Optional = false
	f.Default = nil
	f.Type = fld["type"].(string)

	// Seteo valores leidos
	tmp = fld["id"]
	if tmp != nil {
		f.ID = tmp.(bool)
	}
	tmp = fld["optional"]
	if tmp != nil {
		f.Optional = tmp.(bool)
	}
	tmp = fld["default"]
	if tmp != nil {
		f.Default = tmp
	}
	return &f
}

//Load carga una entidad desde un archivos de configuracion (json)
func (e *Entity) Load(home, name string) (err error) {
	var message string

	defer func() {
		if p := recover(); p != nil {
			err = Error{"Entity '" + name + "' " + message, p}
		}
	}()

	dat, err := ioutil.ReadFile(home + "/data/" + name + ".json")

	if err != nil {
		return err
	}
	message = "Reading Data"
	var body = string(dat)

	message = "Convert to Json"
	var ojson map[string]interface{}
	err = json.Unmarshal([]byte(body), &ojson)

	if err != nil {
		return err
	}

	message = "Setting name and description"
	e.Name = name
	e.Description = ojson["description"].(string)
	e.Fields = make(map[string]*Field)

	fields := ojson["schema"].(map[string]interface{})

	for key, value := range fields {
		message = "Field '" + key + "'"

		fld := value.(map[string]interface{})
		f := newField(key, fld)
		e.Fields[key] = f
	}

	return nil
}

// Validate Valida un json en base al esquema
func (e *Entity) Validate(json map[string]interface{}) error {
	for key, value := range json {
		fld := e.Fields[key]
		if fld == nil {
			return Error{"Entity '" + e.Name + "' ", "Field '" + key + "' not valid"}
		}
		valid := validateValue(fld, value)
		if !valid {
			return Error{"Entity '" + e.Name + "' ", "Field '" + key + "' Invalid Value"}
		}

		fmt.Println(value)

	}
	return nil
}

func validateValue(fld *Field, value interface{}) bool {
	if value == nil {
		if fld.Optional == true {
			return true
		}
		return false
	}
	switch fld.Type {
	case "string":
		_, ok := getString(value)
		return ok
	case "float":
		_, ok := getFloat(value)
		return ok
	case "integer":
		_, ok := getInt(value)
		return ok
	case "date":
		_, ok := getDate(value)
		return ok
	case "bool":
		_, ok := getBool(value)
		return ok
	}
	return false
}

// Convierto interface a tipo de valor
func getInt(value interface{}) (ret int64, ok bool) {
	ret, ok = value.(int64)
	return
}
func getFloat(value interface{}) (ret float64, ok bool) {
	ret, ok = value.(float64)
	return
}
func getString(value interface{}) (ret string, ok bool) {
	ret, ok = value.(string)
	return
}
func getBool(value interface{}) (ret bool, ok bool) {
	ret, ok = value.(bool)
	return
}
func getDate(value interface{}) (ret string, ok bool) {
	ret, ok = value.(string)
	//TODO COC Verifcar Formato
	return
}
