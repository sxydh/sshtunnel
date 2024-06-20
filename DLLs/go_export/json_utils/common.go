package json_utils

import (
	"encoding/json"
	"log"
	"reflect"
)

func ToJsonStr(p interface{}) string {
	empty := "{}"
	if reflect.TypeOf(p).Kind() != reflect.Ptr {
		return empty
	}
	v := reflect.ValueOf(p).Elem()
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return empty
	}
	bytes, err := json.Marshal(v.Interface())
	if err != nil {
		log.Printf("Marshal traverses error: err=%v", err)
		return empty
	}
	return string(bytes)
}
