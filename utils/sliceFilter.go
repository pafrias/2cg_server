package utils

import (
	"errors"
	"fmt"
	"reflect"
)

/*Filter accepts an interface slice and predicate function, calls the predicate function on every element, and
returns a new slice for all elements for which the predicate returned true.

Bad functions and mismatching values/parameters return nil, empty arrays return a new empty array*/
func Filter(src, testFunc interface{}) ([]interface{}, error) {
	var results = []interface{}{}

	err := matchFuncType(reflect.TypeOf(src), reflect.TypeOf(testFunc))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	srcVal := reflect.ValueOf(src)
	funcVal := reflect.ValueOf(testFunc)

	var input [1]reflect.Value
	for i := 0; i < srcVal.Len(); i++ {
		input[0] = srcVal.Index(i)
		b := funcVal.Call(input[:])[0]
		if b.Bool() {
			in := input[0].Interface()
			results = append(results, in)
		}
	}
	return results, nil
}

/*matchFuncType tests for matching elements in a slice and in a functions input, for use as
a predicate in slice utility functions.

It also tests if that function returns bool*/
func matchFuncType(src, test reflect.Type) error {
	var errorStr string

	sourceKind := src.Kind()
	testKind := test.Kind()
	if sourceKind != reflect.Slice {
		errorStr = fmt.Sprintf("Expected type (slice) but received src of type (%v)\n", sourceKind)
	} else if testKind != reflect.Func {
		errorStr = fmt.Sprintf("Expected type (func) but received testFunc of type (%v)\n", testKind)
	}

	testOutput := test.Out(0).Kind()
	testInput := test.In(0)
	elemType := src.Elem()
	if test.NumIn() != 1 {
		errorStr = "Expected testFunc to have 1 input\n"
	} else if testInput != elemType {
		errorStr = fmt.Sprintf("Received slice of type (%v), but testFunc expects type (%v)", elemType, testInput)
	} else if test.NumOut() != 1 || testOutput != reflect.Bool {
		errorStr = "Expected func to have 1 return of type (bool)\n"
	}

	if errorStr != "" {
		return errors.New(errorStr)
	}

	return nil
}

//FilterFunc persists a given  for reuse
type FilterFunc func(arr []interface{}) ([]interface{}, error)

/*NewFilterFunc takes a predicate function and returns a filter that applies that function.

Simply helps readability of code, no performative gains measured.*/
func NewFilterFunc(test interface{}) (FilterFunc, error) {

	return func(arr []interface{}) ([]interface{}, error) {
		return Filter(arr, test)
	}, nil
}
