package entities

import (
	"github.com/coc1961/go/crud/jsonutil"
)

// Entity representa una entidad del crud
type Entity struct {
	definition *Definition
	data       *jsonutil.MJson
}

// Get get attribute value
func (e *Entity) Get(attName string) *jsonutil.MJson {
	return e.data.Get(attName)
}

// Add add attribute value
func (e *Entity) Add(attName string) *jsonutil.MJson {
	return e.data.Add(attName)
}

// JSON return the json
func (e *Entity) JSON() string {
	return e.data.JSON()
}
