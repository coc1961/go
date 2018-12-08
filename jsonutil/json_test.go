package jsonutil

import (
	"testing"
)

const sjson string = `
{
	"id": "xy23",
	"name": "Nombre",
	"amount": 100.45,
	"age": 55,
	"creationDate": "2018-11-24T01:10:22Z",
	"hijo": {
	  "idAtt": "cc23",
	  "nameAtt": "AttrHijo"
	},
	"soyArray": [
	  "elem1",
	  "elem2"
	],
	"arrayOnject": [
	  {
		"id": 123,
		"name": "Nombre1"
	  },
	  {
		"id": 456,
		"name": "Nombre2"
	  }
	]
  }	
 `

func TestJSON_Get(t *testing.T) {
	t.Run("Test Simple Get", func(t *testing.T) {
		json := NewFromString(sjson)
		if got := json.Get("hijo").Get("idAtt").Value().(string); got != "cc23" {
			t.Errorf("JSON.Get() = %v, want %v", got, "cc23")
		}
	})
	t.Run("Test Array of String Get", func(t *testing.T) {
		json := NewFromString(sjson)
		if got := json.Get("soyArray").ValueAsArray(); len(got) != 2 {
			t.Errorf("JSON.Get() = %v, want %v", len(got), 2)
		}
	})
	t.Run("Test Array of Object Get", func(t *testing.T) {
		json := NewFromString(sjson)
		if got := json.Get("arrayOnject").ValueAsArray(); len(got) != 2 {
			t.Errorf("JSON.Get() = %v, want %v", len(got), 2)
		}
	})
	t.Run("Test Array of Object Get ", func(t *testing.T) {
		json := NewFromString(sjson)
		tmp := json.Get("arrayOnject").ValueAsArray()[0]
		tmp1 := tmp.(map[string]interface{})
		json1 := NewFromMap(&tmp1)
		if got := json1.Get("id").Value().(float64); got != float64(123) {
			t.Errorf("JSON.Get() = %v, want %v", got, 123)
		}
	})
}

func TestJSON_GetArray(t *testing.T) {
	t.Run("Test Array of Object Get ", func(t *testing.T) {
		json := NewFromString(sjson)
		tmp := json.Get("arrayOnject").ValueAsArrayOfObjects()
		if got := tmp[0].Get("id").Value().(float64); got != float64(123) {
			t.Errorf("JSON.Get() = %v, want %v", got, 123)
		}
	})
}
