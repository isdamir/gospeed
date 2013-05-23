package utils

import (
	"fmt"
	"reflect"
	"strings"
)

//用于方便的处理Struct的一些信息
type Struct struct {
	I interface{}
}

//获取类型信息
func (s Struct) GetTypeName() string {
	var typestr string
	typ := reflect.TypeOf(s.I)
	typestr = typ.String()

	lastDotIndex := strings.LastIndex(typestr, ".")
	if lastDotIndex != -1 {
		typestr = typestr[lastDotIndex+1:]
	}

	return typestr
}

func (s Struct) StructName() string {
	v := reflect.TypeOf(s.I)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v.Name()
}

// convert struct to map
// s must to be struct, can not be a pointer
func (s Struct) rawStructToMap(snakeCasedKey bool) Map {
	v := reflect.ValueOf(s.I)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic(fmt.Sprintf("param s must be struct, but got %s", s.I))
	}

	m := Map{}
	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Name
		if snakeCasedKey {
			key = Strings(key).SnakeCasedName()
		}
		val := v.Field(i).Interface()

		m[key] = val
	}
	return m
}

// convert struct to map
func (s Struct) StructToMap() Map {
	return s.rawStructToMap(false)
}

// convert struct to map
// but struct's field name to snake cased map key
func (s Struct) StructToSnakeKeyMap() Map {
	return s.rawStructToMap(true)
}
