package entities

import (
	"encoding/json"
	"log"
	"path/filepath"

	"github.com/coc1961/go/crud/jsonutil"
	schema "github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-jsschema/validator"
)

// Definition Define de una Entidad
type Definition struct {
	schema    *schema.Schema
	validator *validator.Validator
}

// Load Lectura de la definicion de una entidad
func (e *Definition) Load(path, entity string) error {

	fullPath := filepath.Join(path, entity+".json")
	s, err := schema.ReadFile(fullPath)
	if err != nil {
		log.Printf("failed to read schema: %s", err)
		return err
	}

	e.schema = s
	v := validator.New(s)
	e.validator = v
	return nil
}

// Validate Valida un json en base al schema
func (e *Definition) Validate(sjson string) (map[string]interface{}, error) {
	var ojson map[string]interface{}
	err := json.Unmarshal([]byte(sjson), &ojson)
	if err != nil {
		return nil, err
	}
	err = e.validator.Validate(ojson)
	return ojson, err
}

// New nueva entidad en base a un json
func (e *Definition) New(sjson string) (*Entity, error) {
	ojson, err := e.Validate(sjson)
	if err != nil {
		return nil, err
	}
	data := jsonutil.New().Set(ojson)
	return &Entity{e, data}, nil
}

/***************************
* Entity
****************************/

// Entity representa una entidad del crud
type Entity struct {
	definition *Definition
	data       *jsonutil.JSON
}

// Get get attribute value
func (e *Entity) Get(attName string) *jsonutil.JSON {
	return e.data.Get(attName)
}

// Add add attribute value
func (e *Entity) Add(attName string) *jsonutil.JSON {
	return e.data.Add(attName)
}

// JSON return the json
func (e *Entity) JSON() string {
	return e.data.JSON()
}
