package entities

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

// Add add attribute value
func (e *Attribute) Add(attName string) *Attribute {
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

/*
// Get get json value
func (e *Entity) Get(path string) interface{} {
	json := e.json
	pt := strings.Split(path, "/")
	var ok bool
	var tmp map[string]interface{}
	var ret interface{}
	for _, p := range pt {
		if p == "" {
			continue
		}
		ret = nil

		tmp, ok = json[p].(map[string]interface{})
		if !ok {
			ret, ok = json[p].(interface{})
		} else {
			json = tmp
			ret = tmp
		}
	}
	return ret
}

// Set set json value
func (e *Entity) Set(path string, value interface{}) {
	json := &e.json
	pt := strings.Split(path, "/")
	cont := len(pt)
	var lastPt string
	for _, p := range pt {
		cont--
		if p == "" {
			continue
		}
		if cont == 0 {
			lastPt = p
			break
		}
		tt := (*json)[p]
		t, ok := tt.(map[string]interface{})
		if !ok {
			return
		}
		json = &t
	}
	(*json)[lastPt] = value
}
*/
