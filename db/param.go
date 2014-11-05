// Copyright 2014 li. All rights reserved.

package db

import (
	"fmt"
	"reflect"
)

func procParam(names []string, param interface{}) []interface{} {
	fmt.Println(param)
	fmt.Println(names)
	if param == nil {
		return make([]interface{}, len(names))
	}

	v := reflect.ValueOf(param)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Struct:
		params := make([]interface{}, len(names))
		for i, name := range names {
			// Compatible if name is lower case.
			params[i] = v.FieldByName(toUpperCase(name)).Interface()
		}
		return params

	case reflect.Map:
		params := make([]interface{}, len(names))
		for i, name := range names {
			params[i] = v.MapIndex(reflect.ValueOf(name)).Interface()
		}
		return params

	default:
		return []interface{}{param}
	}
}
