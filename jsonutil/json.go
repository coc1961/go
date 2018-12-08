package jsonutil

import (
	"encoding/json"
)

// JSON Reprernta un Json
type JSON struct {
	rootValue *map[string]interface{}
	path      []string
	index     int
}

/******************
** Factory
*******************/

// New creo un objeto MJson vacio
func New() *JSON {
	pt := make([]string, 0)
	entity := make(map[string]interface{})
	return &JSON{&entity, pt, -1}
}

// NewFromMap creo un objeto MJson
func NewFromMap(rootValue *map[string]interface{}) *JSON {
	pt := make([]string, 0)
	return &JSON{rootValue, pt, -1}
}

// NewFromString creo un objeto MJson
func NewFromString(sjson string) *JSON {
	pt := make([]string, 0)
	var entity map[string]interface{}
	err := json.Unmarshal([]byte(sjson), &entity)
	if err != nil {
		return nil
	}
	return &JSON{&entity, pt, -1}
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
	if err != nil {
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
	return &JSON{&tmp, pt, -1}
}

/**********************
** Get and Set Values
***********************/

// Get get attribute value
func (e *JSON) Get(attName string) *JSON {
	tmpPath := e.path
	tmpPath = append(tmpPath, attName)
	if len(e.path) == 0 {
		return &JSON{e.rootValue, tmpPath, e.index}
	}
	tmp := e.internalValue()
	if tmp != nil {
		_, ok := (*tmp).(map[string]interface{})
		if ok {
			return &JSON{e.rootValue, tmpPath, e.index}
		}
	}
	if tmp != nil {
		_, ok := (*tmp).([]interface{})
		if ok {
			return &JSON{e.rootValue, tmpPath, e.index}
		}
	}

	return nullMJson()
}

// Set set attribute value
func (e *JSON) Set(value interface{}) *JSON {
	tmp, ok := value.(*JSON)
	if ok {
		value = tmp.rootValue
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
	return &JSON{e.rootValue, tmpPath, -1}
}

/******************
** Object Value
*******************/

// Value get attribute value
func (e *JSON) Value() interface{} {
	ret := e.internalValue()
	if ret == nil {
		return nil
	}
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

// ValueAsArrayOfObjects get attribute value
func (e *JSON) ValueAsArrayOfObjects() []JSON {
	tmpValue := e.Value()
	if tmpValue == nil {
		return nil
	}
	ret, ok := (tmpValue).([]interface{})
	if ok {
		arr := make([]JSON, 0)
		for it, i1 := range ret {
			_ = i1
			arr = append(arr, JSON{e.rootValue, e.path, it})
		}
		return arr
	}
	return nil
}

/**********************
** Internal Functions
***********************/

// Valor del Objeto
func (e *JSON) internalValue() *interface{} {
	tmp, lastPt := e.parentPath()
	if tmp != nil {
		tmp1 := (*tmp)[lastPt]
		return &tmp1
	}
	tmp1, lastPt1 := e.arrayParentPath()
	if tmp1 != nil {
		tmp2 := (*tmp1)[e.index]
		tmp3 := tmp2.(map[string]interface{})
		tmp4 := tmp3[lastPt1]
		return &tmp4
	}
	return nil
}

// map del parent y nombre del atributo del objeto actual
func (e *JSON) parentPath() (*map[string]interface{}, string) {
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
			return nil, ""
		}
		json = &tmpMap
	}
	return json, lastPt
}
func (e *JSON) arrayParentPath() (*[]interface{}, string) {
	json := e.rootValue
	path := e.path
	pathLen := len(path)
	var lastPt string
	var retArray []interface{}
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
			tmpArr, ok := tmpInterface.([]interface{})
			if ok {
				retArray = tmpArr
			}
		}
		json = &tmpMap
	}
	return &retArray, lastPt
}
