package utils

import (
	"fmt"
	"reflect"
)

/*Any traverses the src slice, invoking the predicate function on all all elements.

If any invocation returns true, return the index of first occurance.
Else, return -1 (and error is applicable).*/
func Any(src, testFunc interface{}) (int, error) {
	err := matchFuncType(reflect.TypeOf(src), reflect.TypeOf(testFunc))
	if err != nil {
		fmt.Println(err)
	} else {
		srcVal := reflect.ValueOf(src)
		funcVal := reflect.ValueOf(testFunc)

		var input [1]reflect.Value
		for i := 0; i < srcVal.Len(); i++ {
			input[0] = srcVal.Index(i)
			b := funcVal.Call(input[:])[0]
			if b.Bool() {
				return i, nil
			}
		}
	}
	return -1, err
}

/*Every traverses the src slice, invoking the predicate function on all all elements.

If all invocations return true, return true.
Else, return false (and error is applicable).*/
func Every(src, testFunc interface{}) (bool, error) {
	err := matchFuncType(reflect.TypeOf(src), reflect.TypeOf(testFunc))
	if err != nil {
		fmt.Println(err)
	} else {
		srcVal := reflect.ValueOf(src)
		funcVal := reflect.ValueOf(testFunc)

		var input [1]reflect.Value
		for i := 0; i < srcVal.Len(); i++ {
			input[0] = srcVal.Index(i)
			b := funcVal.Call(input[:])[0]
			if !b.Bool() {
				return false, nil
			}
		}
	}
	return true, err
}
