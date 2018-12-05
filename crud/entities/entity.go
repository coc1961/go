package entities

import (
	"encoding/json"

	"github.com/coc1961/go/crud/jsonutil"
)

// Entity representa una entidad del crud
type Entity struct {
	definition *Definition
	json       map[string]interface{}
}

// Get get attribute value
func (e *Entity) Get(attName string) *jsonutil.MJson {
	tmp := e.json[attName]
	if tmp == nil {
		return nil
	}
	pt := make([]string, 1)
	pt = append(pt, attName)
	return jsonutil.New().SetValue(tmp).SetRootValue(&e.json).SetPath(pt)
}

// AddObject add attribute value
func (e *Entity) AddObject(attName string) *jsonutil.MJson {
	e.json[attName] = ""
	tmpPath := make([]string, 1)
	tmpPath = append(tmpPath, attName)
	return jsonutil.New().SetValue(e.json[attName]).SetRootValue(&e.json).SetPath(tmpPath)
}

// JSON return the json
func (e *Entity) JSON() string {
	b, err := json.Marshal(e.json)
	if err != nil {
		return ""
	}
	return string(b)
}
