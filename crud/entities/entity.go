package entities

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/coc1961/go/jsonutil"
	schema "github.com/lestrrat/go-jsschema"
	"github.com/lestrrat/go-jsschema/validator"
)

// Definition Define una Entidad
type Definition struct {
	schema    *schema.Schema
	validator *validator.Validator
	json      *jsonutil.JSON
	byteRead  int
}

// NewEntityDefinition  New Entity Definition
func NewEntityDefinition() *Definition {
	return &Definition{nil, nil, nil, 0}
}

// IsValid IsValid
func (e *Definition) IsValid() bool {
	if e == nil || e.schema == nil || e.validator == nil || e.json == nil {
		return false
	}
	return true
}

// Load Lectura de la definicion de una entidad
func (e *Definition) Load(path, entity string) error {

	fullPath := filepath.Join(path, entity)

	b, err := ioutil.ReadFile(fullPath)
	if err != nil {
		return err
	}
	content := string(b)
	e.json = jsonutil.NewFromString(content)
	e.byteRead = 0

	s, err := schema.Read(e)
	if err != nil {
		log.Printf("failed to read schema: %s", err)
		return err
	}

	e.schema = s
	v := validator.New(s)
	e.validator = v
	return nil
}

func (e *Definition) Read(p []byte) (n int, err error) {
	str := e.json.Get("schema").JSON()
	bt := []byte(str)
	cont := len(p)
	ind := 0
	for i := e.byteRead; i < len(bt); i++ {
		if ind >= cont {
			break
		}
		p[ind] = bt[i]
		ind++
	}
	e.byteRead += ind
	n = ind
	err = nil
	return
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
	data := jsonutil.New().Set(&ojson)
	return &Entity{e, data}, nil
}

// Name get entity Name
func (e *Definition) Name() string {
	str := e.json.Get("name").Value()
	if str == nil {
		return ""
	}
	return str.(string)
}

/***************************
* Entity
****************************/

// Entity representa una entidad del crud
type Entity struct {
	definition *Definition
	data       *jsonutil.JSON
}

// Set set attribute value
func (e *Entity) Set(attName string) *jsonutil.JSON {
	return e.data.Set(attName)
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
