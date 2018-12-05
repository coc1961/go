package util

import (
	"encoding/json"
)

// MJson Reprernta un atributo
type MJson struct {
	localValue interface{}
	rootValue  *map[string]interface{}
	path       []string
}

// NewFromMap creo un objeto MJson
func NewFromMap(rootValue *map[string]interface{}) *MJson {
	pt := make([]string, 0)
	return &MJson{nil, rootValue, pt}
}

// NewFromString creo un objeto MJson
func NewFromString(sjson string) *MJson {
	pt := make([]string, 0)
	var entity map[string]interface{}
	err := json.Unmarshal([]byte(sjson), &entity)
	if err != nil {
		return nil
	}
	return &MJson{nil, &entity, pt}
}

// GetRoot retorna el map interno
func (e *MJson) GetRoot() *map[string]interface{} {
	return e.rootValue
}

// SetValue setea el valor local
func (e *MJson) SetValue(localValue interface{}) *MJson {
	e.localValue = localValue
	return e
}

// SetRootValue setea el valor del root
func (e *MJson) SetRootValue(rootValue *map[string]interface{}) *MJson {
	e.rootValue = rootValue
	return e
}

// SetPath setea el valor del path
func (e *MJson) SetPath(path []string) *MJson {
	e.path = path
	return e
}

// Get get attribute value
func (e *MJson) Get(attName string) *MJson {
	tmpPath := e.path
	tmpPath = append(tmpPath, attName)
	ret, ok := e.localValue.(map[string]interface{})
	if ok {
		localValue := ret[attName]
		return &MJson{localValue, e.rootValue, tmpPath}
	}
	return nil
}

// Set set attribute value
func (e *MJson) Set(value interface{}) {
	tmp, ok := value.(*MJson)
	if ok {
		value = tmp.rootValue
	}
	json := e.rootValue
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
	if len(path) == 0 {
		for k, v := range value.(map[string]interface{}) {
			(*json)[k] = v
		}
	} else {
		(*json)[lastPt] = value
	}
}

// AddObject add attribute value
func (e *MJson) AddObject(attName string) *MJson {
	tmpObject := make(map[string]interface{})
	tmpObject[attName] = ""
	e.Set(tmpObject)
	tmpPath := e.path
	tmpPath = append(tmpPath, attName)
	return &MJson{tmpObject[attName], e.rootValue, tmpPath}
}

// Value get attribute value
func (e *MJson) Value() interface{} {
	return e.localValue
}

// ValueAsArray get attribute value
func (e *MJson) ValueAsArray() []interface{} {
	ret, ok := e.localValue.([]interface{})
	if ok {
		return ret
	}
	return nil
}
