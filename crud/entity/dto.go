package entity

import "time"

// Entity representa una entidad
type Entity struct {
	Name      string
	Atributes map[string]*Attribute
}

// Attribute repesenta un atributo
type Attribute struct {
	FieldDefinition *FieldDefinition
	Value           interface{}
	Child           map[string]*Attribute
}

// Int valor como int
func (a *Attribute) Int() (int64, bool) {
	if a.Value == nil {
		return 0, false
	}
	return getInt(a.Value)
}

// Float valor como float
func (a *Attribute) Float() (float64, bool) {
	if a.Value == nil {
		return 0, false
	}
	return getFloat(a.Value)
}

// String valor como string
func (a *Attribute) String() (string, bool) {
	if a.Value == nil {
		return "", false
	}
	return getString(a.Value)
}

// Bool valor como boll
func (a *Attribute) Bool() (bool, bool) {
	if a.Value == nil {
		return false, false
	}
	return getBool(a.Value)
}

// Date valor como Time
func (a *Attribute) Date() (ret time.Time, ok bool) {
	if a.Value == nil {
		ret, _ = time.Parse(time.RFC3339, "1900-01-01T00:00:00Z")
		return ret, false
	}
	return getDate(a.Value)
}
