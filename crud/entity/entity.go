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
	Fields      []Field
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
	return fmt.Sprintf("panic: %s => %v", p.message, p.reason)
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

	// Seteo valores leidos
	f.Type = fld["type"].(string)
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

	fields := ojson["schema"].(map[string]interface{})

	for key, value := range fields {
		message = "Field '" + key + "'"

		fld := value.(map[string]interface{})
		e.Fields = append(e.Fields, *newField(key, fld))
	}

	return nil
}
