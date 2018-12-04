package util

import "encoding/json"

// MJson Reprernta un atributo
type MJson struct {
	localValue interface{}
	rootValue  *map[string]interface{}
	path       []string
}

// New creo un objeto MJson
func New(sjson string) *MJson {
	pt := make([]string, 1)
	var entity map[string]interface{}
	err := json.Unmarshal([]byte(sjson), &entity)
	if err != nil {
		return nil
	}
	return &MJson{nil, &entity, pt}
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
	(*json)[lastPt] = value
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
