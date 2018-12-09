package jsonutil

import (
	"encoding/json"
)

// JSON Reprernta un Json
type JSON struct {
	rootValue *map[string]interface{}
	path      []string
	isNil     bool
}

/******************
** Factory
*******************/

// New creo un objeto MJson vacio
func New() *JSON {
	pt := make([]string, 0)
	entity := make(map[string]interface{})
	return internalNew(&entity, pt, false)
}

// NewFromMap creo un objeto MJson
func NewFromMap(rootValue *map[string]interface{}) *JSON {
	pt := make([]string, 0)
	return internalNew(rootValue, pt, false)
}

// NewFromString creo un objeto MJson
func NewFromString(sjson string) *JSON {
	pt := make([]string, 0)
	var entity map[string]interface{}
	err := json.Unmarshal([]byte(sjson), &entity)
	if err != nil {
		return nil
	}
	return internalNew(&entity, pt, false)
}

func internalNew(rootValue *map[string]interface{}, path []string, isNil bool) *JSON {
	return &JSON{rootValue, path, isNil}
}

/******************
** Getter y Setter
*******************/

// GetRoot retorna el map interno
func (e *JSON) GetRoot() *map[string]interface{} {
	return e.rootValue
}

// SetRootValue setea el valor del root
func (e *JSON) SetRootValue(rootValue *map[string]interface{}) *JSON {
	e.rootValue = rootValue
	return e
}

// JSON return the json
func (e *JSON) JSON() string {
	b, err := json.Marshal(e.GetRoot())
	if err != nil || string(b) == "null" {
		return ""
	}
	return string(b)
}

/******************
** Invalid Object
*******************/

func nullMJson() *JSON {
	tmp := make(map[string]interface{})
	pt := make([]string, 0)
	return internalNew(&tmp, pt, true)
}

/**********************
** Get and Set Values
***********************/

// Get get attribute value
func (e *JSON) Get(attName string) *JSON {
	tmpPath := e.path
	tmpPath = append(tmpPath, attName)
	if len(e.path) == 0 {
		if _, ok := (*e.rootValue)[attName]; ok {
			return internalNew(e.rootValue, tmpPath, false)
		}
	}
	tmp := e.internalValue()
	if tmp != nil {
		_, ok := (*tmp).(map[string]interface{})
		if ok {
			return internalNew(e.rootValue, tmpPath, false)
		}
	}

	return nullMJson()
}

// Set set attribute value
func (e *JSON) Set(value interface{}) *JSON {
	tmp, ok := value.(*JSON)
	if ok {
		value = *tmp.rootValue
	}
	json, lastPt := e.parentPath()
	if len(e.path) == 0 {
		tmpArr, okArr := value.(map[string]interface{})
		if okArr {
			for k, v := range tmpArr {
				(*json)[k] = v
			}
		} else {
			tmpArr, okArr := value.(*map[string]interface{})
			if okArr {
				for k, v := range *tmpArr {
					(*json)[k] = v
				}
			}
		}
	} else {
		(*json)[lastPt] = value
	}
	return e
}

// Add add attribute value
func (e *JSON) Add(attName string) *JSON {
	tmpObject := make(map[string]interface{})
	tmpObject[attName] = ""
	e.Set(tmpObject)
	tmpPath := e.path
	tmpPath = append(tmpPath, attName)
	return internalNew(e.rootValue, tmpPath, false)
}

/******************
** Object Value
*******************/

// IsNil return if nil object
func (e *JSON) IsNil() bool {
	return e.isNil
}

// Value get attribute value
func (e *JSON) Value() interface{} {
	ret := e.internalValue()
	return *ret
}

// ValueAsArray get attribute value
func (e *JSON) ValueAsArray() []interface{} {
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

/**********************
** Internal Functions
***********************/

// Valor del Objeto
func (e *JSON) internalValue() *interface{} {
	tmp, lastPt := e.parentPath()
	tmp1 := (*tmp)[lastPt]
	return &tmp1
}

// map del parent y nombre del atributo del objeto actual
func (e *JSON) parentPath() (*map[string]interface{}, string) {
	json := e.rootValue
	path := e.path
	pathLen := len(path)
	var lastPt string
	for _, p := range path {
		pathLen--
		if pathLen == 0 {
			lastPt = p
			break
		}
		tmpInterface := (*json)[p]
		tmpMap, ok := tmpInterface.(map[string]interface{})
		if !ok {
			return nil, ""
		}
		json = &tmpMap
	}
	return json, lastPt
}
