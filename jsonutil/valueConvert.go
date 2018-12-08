package jsonutil

import "time"

/***************************************************************
 ** Funciones de Conversion
 ***************************************************************/

// Convierto interface a tipo de valor
func getInt(value interface{}) (ret int64, ok bool) {
	ok = false
	flo, ok1 := getFloat(value)
	ok = ok1
	if ok1 {
		ret = int64(flo)
	}
	return
}
func getFloat(value interface{}) (ret float64, ok bool) {
	ret, ok = value.(float64)
	return
}
func getString(value interface{}) (ret string, ok bool) {
	ret, ok = value.(string)
	return
}
func getBool(value interface{}) (ret bool, ok bool) {
	ret, ok = value.(bool)
	return
}

func getDate(value interface{}) (ret time.Time, ok bool) {
	var err error
	var rets string
	rets, ok = getString(value)
	if ok {
		ret, err = time.Parse(time.RFC3339, rets)
		if err != nil {
			ok = false
		}
	}
	return
}

func getMap(value interface{}) (ret map[string]interface{}, ok bool) {
	ret, ok = value.(map[string]interface{})
	return
}
