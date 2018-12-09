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

func TestJSON_New(t *testing.T) {
	t.Run("Test New New", func(t *testing.T) {
		json := New()
		if got := json.GetRoot(); len(*got) != 0 {
			t.Errorf("JSON.Get() = %v, want %v", got, 0)
		}
		if got := json.JSON(); got != "{}" {
			t.Errorf("JSON.Get() = %v, want %v", got, "{}")
		}
	})

	t.Run("Test New New", func(t *testing.T) {
		json := New()
		mp := make(map[string]interface{}, 0)
		mp["id"] = 1
		json.SetRootValue(&mp)
		if got := json.GetRoot(); len(*got) != 1 {
			t.Errorf("JSON.Get() = %v, want %v", len(*got), 1)
		}
	})

	t.Run("Test NewFromMap New", func(t *testing.T) {
		mp := make(map[string]interface{})
		mp["id"] = "1"
		json := NewFromMap(&mp)
		if got := json.GetRoot(); len(*got) != 1 {
			t.Errorf("JSON.Get() = %v, want %v", got, 1)
		}
		if got := json.Get("id").Value().(string); got != "1" {
			t.Errorf("JSON.Get() = %v, want %v", got, "1")
		}
	})
	t.Run("Test NewFromString New", func(t *testing.T) {
		json := NewFromString(`{"id": 1}`)
		if got := json.GetRoot(); len(*got) != 1 {
			t.Errorf("JSON.Get() = %v, want %v", got, 1)
		}
		if got := json.Get("id").Value().(float64); got != float64(1) {
			t.Errorf("JSON.Get() = %v, want %v", got, 1)
		}

		json = NewFromString(`{"id: 1}`)
		if got := json; got != nil {
			t.Errorf("JSON.Get() = %v, want %v", got, nil)
		}

		json = NewFromString(`{"id": 1}`)
		if got := json.Get("id1"); got.IsNil() != true {
			t.Errorf("JSON.Get() = %v, want %v", got.IsNil(), true)
		}
		if got := json.Get("id"); got.IsNil() != false {
			t.Errorf("JSON.Get() = %v, want %v", got.IsNil(), false)
		}

	})
}

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
	t.Run("Test Array of Object Get 1", func(t *testing.T) {
		json := NewFromString(sjson)
		tmp := json.Get("arrayOnject").ValueAsArray()[0]
		tmp1 := tmp.(map[string]interface{})
		json1 := NewFromMap(&tmp1)
		if got := json1.Get("id").Value().(float64); got != float64(123) {
			t.Errorf("JSON.Get() = %v, want %v", got, 123)
		}
	})
	t.Run("Test Array of Nil Get", func(t *testing.T) {
		json := NewFromString(sjson)
		if got := json.Get("soyArray1").ValueAsArray(); got != nil {
			t.Errorf("JSON.Get() = %v, want %v", got, nil)
		}
	})
	t.Run("Test Null Value Get", func(t *testing.T) {
		json := NewFromString(sjson)
		if got := json.Get("id").Get("NoEsta").Get("NoEsta2").Value(); got != nil {
			t.Errorf("JSON.Get() = %v, want %v", got, nil)
		}
	})

	t.Run("Test Object Of nill Get", func(t *testing.T) {
		json := NewFromString(sjson)
		json.Get("arrayOnject").Set(nil)
		if got := json.Get("arrayOnject").Value(); got != nil {
			t.Errorf("JSON.Get() = %v, want %v", got, nil)
		}
		if got := json.Get("arrayOnject").Get("id").Value(); got != nil {
			t.Errorf("JSON.Get() = %v, want %v", got, nil)
		}
	})

}

func TestJSON_Set(t *testing.T) {
	t.Run("Add Array Object Set", func(t *testing.T) {
		json := NewFromString(sjson)

		ent2 := New()
		ent2.Add("id").Set(788)
		ent2.Add("name").Set("Nombre2")

		tmp := json.Get("arrayOnject").ValueAsArray()
		tmp = append(tmp, ent2.GetRoot())
		json.Get("arrayOnject").Set(tmp)

		if got := json.Get("arrayOnject").ValueAsArray(); len(got) != 3 {
			t.Errorf("JSON.Get() = %v, want %v", got, 3)
		}
		element2 := json.Get("arrayOnject").ValueAsArray()[2]
		mapElement2 := element2.(*map[string]interface{})
		if got := (*mapElement2)["name"].(string); got != "Nombre2" {
			t.Errorf("JSON.Get() = %v, want %v", got, "Nombre2")
		}
		if got := (*mapElement2)["id"].(int); got != int(788) {
			t.Errorf("JSON.Get() = %v, want %v", got, 788)
		}
	})

	t.Run("Add Map Set", func(t *testing.T) {
		json := NewFromString(sjson)
		mp := make(map[string]interface{})
		mp["id"] = "1"
		json.Add("pri").Set(mp)

		if got := json.Get("pri").Get("id").Value(); got != "1" {
			t.Errorf("JSON.Get() = %v, want %v", got, "1")
		}

		json = NewFromString(sjson)
		json.Set(&mp)

		if got := json.Get("id").Value(); got != "1" {
			t.Errorf("JSON.Get() = %v, want %v", got, "1")
		}
	})

}
func TestJSON_Add(t *testing.T) {
	t.Run("Add String Value Add", func(t *testing.T) {
		json := NewFromString(sjson)
		json.Get("hijo").Get("idAtt").Add("Otro").Set("Prueba1")

		if got := json.Get("hijo").Get("idAtt").Get("Otro").Value().(string); got != "Prueba1" {
			t.Errorf("JSON.Get() = %v, want %v", got, "Prueba1")
		}

		json.Get("hijo").Get("idAtt").Get("Otro").Set("Prueba2")

		if got := json.Get("hijo").Get("idAtt").Get("Otro").Value().(string); got != "Prueba2" {
			t.Errorf("JSON.Get() = %v, want %v", got, "Prueba2")
		}
	})
	t.Run("Add Object Value Add", func(t *testing.T) {
		json := NewFromString(sjson)

		ent2 := New()
		ent2.Add("id").Set(788)
		ent2.Add("name").Set("Nombre2")

		json.Get("hijo").Get("idAtt").Add("Otro").Set(ent2)

		tmp := json.Get("hijo").Get("idAtt").Get("Otro").Value()
		entRead := tmp.(map[string]interface{})

		if got := entRead["id"].(int); got != 788 {
			t.Errorf("JSON.Get() = %v, want %v", got, 788)
		}
		if got := entRead["name"].(string); got != "Nombre2" {
			t.Errorf("JSON.Get() = %v, want %v", got, "Nombre2")
		}

		tmp2 := json.Get("hijo").Get("idAtt").Get("Otro").Get("id").Value()
		_ = tmp2

		if got := json.Get("hijo").Get("idAtt").Get("Otro").Get("id").Value().(int); got != 788 {
			t.Errorf("JSON.Get() = %v, want %v", got, 788)
		}

		json.Get("hijo").Get("idAtt").Get("Otro").Get("id").Set(1000)

		if got := json.Get("hijo").Get("idAtt").Get("Otro").Get("id").Value().(int); got != 1000 {
			t.Errorf("JSON.Get() = %v, want %v", got, 1000)
		}
	})

}
