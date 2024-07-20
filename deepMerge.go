package main

import (
	"reflect"
)

// isZeroValue checks if the given value is a zero value
func isZeroValue(v reflect.Value) bool {
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

// deepMerge recursively merges src into dest
func DeepMerge(dest, src interface{}) {
	destVal := reflect.ValueOf(dest).Elem()
	srcVal := reflect.ValueOf(src).Elem()

	for i := 0; i < destVal.NumField(); i++ {
		destField := destVal.Field(i)
		srcField := srcVal.Field(i)

		if !srcField.IsValid() || !destField.CanSet() || isZeroValue(srcField) {
			continue
		}

		switch destField.Kind() {
		case reflect.Struct:
			// Recursively merge nested structs
			DeepMerge(destField.Addr().Interface(), srcField.Addr().Interface())
		default:
			// Overwrite dest field with src field value
			destField.Set(srcField)
		}
	}
}
