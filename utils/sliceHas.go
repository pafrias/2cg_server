package utils

import (
	"errors"
	"reflect"
)

func SliceHas(src, target interface{}) (bool, error) {

	_, ok := matchSliceType(reflect.TypeOf(src), reflect.TypeOf(target))
	if !ok {
		return false, errors.New("Cannot search for element in slice of wrong type")
	}

	srcVal := reflect.ValueOf(src)
	tarVal := reflect.ValueOf(target)

	for i := 0; i < srcVal.Len(); i++ {
		if tarVal == srcVal.Index(i) {
			return true, nil
		}
	}
	return false, nil
}

func matchSliceType(src, test reflect.Type) (reflect.Type, bool) {
	sk := src.Kind()
	if sk != reflect.Slice {
		return nil, false
	}
	elemType := src.Elem()
	return elemType, elemType == test
}
