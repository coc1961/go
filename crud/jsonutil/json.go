package jsonutil

import (
	"encoding/json"
)

// MJson Reprernta un atributo
type MJson struct {
	rootValue *map[string]interface{}
	path      []string
}

// New creo un objeto MJson vacio
func New() *MJson {
	pt := make([]string, 0)
	entity := make(map[string]interface{})
	return &MJson{&entity, pt}
}

// NewFromMap creo un objeto MJson
func NewFromMap(rootValue *map[string]interface{}) *MJson {
	pt := make([]string, 0)
	return &MJson{rootValue, pt}
}

// NewFromString creo un objeto MJson
func NewFromString(sjson string) *MJson {
	pt := make([]string, 0)
	var entity map[string]interface{}
	err := json.Unmarshal([]byte(sjson), &entity)
	if err != nil {
		return nil
	}
	return &MJson{&entity, pt}
}

// GetRoot retorna el map interno
func (e *MJson) GetRoot() *map[string]interface{} {
	return e.rootValue
}

// SetRootValue setea el valor del root
func (e *MJson) SetRootValue(rootValue *map[string]interface{}) *MJson {
	e.rootValue = rootValue
	return e
}

func nullMJson() *MJson {
	tmp := make(map[string]interface{})
	pt := make([]string, 0)
	return &MJson{&tmp, pt}
}

// Get get attribute value
func (e *MJson) Get(attName string) *MJson {
	tmpPath := e.path
	tmpPath = append(tmpPath, attName)
	if len(e.path) == 0 {
		return &MJson{e.rootValue, tmpPath}
	} else {
		tmp := e.internalGet()
		if tmp != nil {
			_, ok := (*tmp).(map[string]interface{})
			if ok {
				return &MJson{e.rootValue, tmpPath}
			}
		}
	}
	return nullMJson()
}

func (e *MJson) internalGet() *interface{} {
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
			return nil
		}
		json = &tmpMap
	}
	ret := (*json)[lastPt]
	return &ret
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
	return &MJson{e.rootValue, tmpPath}
}

// Value get attribute value
func (e *MJson) Value() interface{} {
	ret := e.internalGet()
	if ret == nil {
		return nil
	}
	return *ret
}

// ValueAsArray get attribute value
func (e *MJson) ValueAsArray() []interface{} {
	tmpValue := e.Value()
	if tmpValue == nil {
		return nil
	}
	ret, ok := (tmpValue).([]interface{})
	if ok {
		return ret
	}
	return nil
}
