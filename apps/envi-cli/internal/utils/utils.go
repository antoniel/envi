package utils

import (
	"fmt"
	"reflect"

	F "github.com/IBM/fp-go/function"
)

func generalSetter[T any, S any](fieldName string, fieldValue T, state S) S {
	val := reflect.ValueOf(state)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Make sure we're dealing with a struct
	if val.Kind() != reflect.Struct {
		panic("❌ State must be a struct or a pointer to struct")
	}

	// Make a copy of the struct to avoid mutating the original
	newState := reflect.New(val.Type()).Elem()
	newState.Set(val)

	// Set the field value
	fld := newState.FieldByName(fieldName)
	if fld.IsValid() && fld.CanSet() {
		fld.Set(reflect.ValueOf(fieldValue))
	} else {
		panic(fmt.Sprintf("🚨 Field `%s` not found or not settable, error at `GeneralSetter`", fieldName))
	}

	return newState.Interface().(S)
}

func Setter[T, S any](fieldName string) func(T) func(S) S {
	type TypeOfGeneralSetter = func(string, T, S) S
	generalSetterBound := F.Bind1of3[TypeOfGeneralSetter](generalSetter)(fieldName)
	type TypeOfGeneralSetterBound = func(T, S) S
	return F.Curry2[TypeOfGeneralSetterBound](generalSetterBound)
}
