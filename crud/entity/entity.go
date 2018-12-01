package entity

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

// Definition Representa una Entidad a administrar con el crud
type Definition struct {
	Name              string
	Description       string
	FieldsRefinitions map[string]*FieldDefinition
}

// FieldDefinition Representa un atributo de una Entity
type FieldDefinition struct {
	Name     string
	Type     string
	ID       bool
	Optional bool
	Default  interface{}
	Child    map[string]*FieldDefinition
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
func newField(name string, fld map[string]interface{}) *FieldDefinition {
	var tmp interface{}
	var f FieldDefinition

	// Valores Default
	f.Name = name
	f.ID = false
	f.Optional = false
	f.Default = nil
	f.Type = fld["type"].(string)
	f.Child = nil

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

	if f.Type == "object" {
		arrList := fld["attributes"].(map[string]interface{})
		f.Child = make(map[string]*FieldDefinition)

		for k, v := range arrList {
			f.Child[k] = newField(k, v.(map[string]interface{}))
			fmt.Println(k)
		}

	}
	return &f
}

//Load carga una entidad desde un archivos de configuracion (json)
func (e *Definition) Load(home, name string) (err error) {
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
	e.FieldsRefinitions = make(map[string]*FieldDefinition)

	fields := ojson["schema"].(map[string]interface{})

	for key, value := range fields {
		message = "Field '" + key + "'"

		fld := value.(map[string]interface{})
		f := newField(key, fld)
		e.FieldsRefinitions[key] = f
	}

	return nil
}

// Parse Valida un json en base al esquema y retorna un objeto Entity
func (e *Definition) Parse(json map[string]interface{}) (*Entity, error) {
	ent := Entity{e.Name, make(map[string]*Attribute)}
	var err error
	ent.Atributes, err = e.fields(json)
	if err != nil {
		return nil, err
	}
	return &ent, err
}

func (e *Definition) fields(json map[string]interface{}) (map[string]*Attribute, error) {
	attributes := make(map[string]*Attribute)
	for key, value := range json {
		fld := e.FieldsRefinitions[key]
		if fld == nil {
			return nil, Error{"Entity '" + e.Name + "' ", "Field '" + key + "' not valid"}
		}
		valid := validateValue(fld, value)
		if !valid {
			return nil, Error{"Entity '" + e.Name + "' ", "Field '" + key + "' Invalid Value"}
		}
		attr := Attribute{fld, value, nil}
		attributes[key] = &attr
		if fld.Type == "object" {
			var err error
			attr.Child, err = e.fields(value.(map[string]interface{}))
			if err != nil {
				return nil, err
			}
		}
	}
	return attributes, nil
}

// Valida un valor en base a la configuración
func validateValue(fld *FieldDefinition, value interface{}) bool {
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
	ok = false
	flo, ok1 := getFloat(value)
	ok = ok1
	if ok1 {
		ret = int64(flo)
	}
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
func getDate(value interface{}) (ret time.Time, ok bool) {
	var err error
	var rets string
	rets, ok = getString(value)
	if ok {
		ret, err = time.Parse(time.RFC3339, rets)
		if err != nil {
			ok = false
		}
	}
	return
}
