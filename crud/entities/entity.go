package entities

import (
	"encoding/json"
)

// Entity representa una entidad del crud
type Entity struct {
	definition *Definition
	json       map[string]interface{}
}

// Get get attribute value
func (e *Entity) Get(attName string) *Attribute {
	tmp := e.json[attName]
	if tmp == nil {
		return nil
	}
	pt := make([]string, 1)
	pt = append(pt, attName)
	return &Attribute{tmp, e, pt}
}

// JSON return the json
func (e *Entity) JSON() string {
	b, err := json.Marshal(e.json)
	if err != nil {
		return ""
	}
	return string(b)
}

// Attribute Reprernta un atributo
type Attribute struct {
	json   interface{}
	entity *Entity
	path   []string
}

// Get get attribute value
func (e *Attribute) Get(attName string) *Attribute {
	pt1 := e.path
	pt1 = append(pt1, attName)
	ret, ok := e.json.(map[string]interface{})
	if ok {
		ret1 := ret[attName]
		return &Attribute{ret1, e.entity, pt1}
	}
	return nil
}

// Set set attribute value
func (e *Attribute) Set(value interface{}) {
	json := &e.entity.json
	path := e.path
	pathLen := len(path)
	var lastPt string
	for _, p := range path {
		pathLen--
		if p == "" {
			continue
		}
		if pathLen == 0 {
			lastPt = p
			break
		}
		tmpInterface := (*json)[p]
		tmpMap, ok := tmpInterface.(map[string]interface{})
		if !ok {
			return
		}
		json = &tmpMap
	}
	(*json)[lastPt] = value
}

// AddObject add attribute value
func (e *Attribute) AddObject(attName string) *Attribute {
	tmp := make(map[string]interface{})
	tmp[attName] = ""
	e.Set(tmp)
	pt := e.path

	pt = append(pt, attName)
	return &Attribute{tmp[attName], e.entity, pt}
}

// Value get attribute value
func (e *Attribute) Value() interface{} {
	return e.json
}
