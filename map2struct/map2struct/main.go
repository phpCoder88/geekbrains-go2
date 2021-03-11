package map2struct

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrNotPinter        = errors.New("not a pointer")
	ErrNotStruct        = errors.New("not a struct type")
	ErrIncompatibleType = errors.New("incompatible type")
	ErrCannotSetValue   = errors.New("cannot set value")
)

func MapToStruct(inStruct interface{}, mapStructValues map[string]interface{}) error {
	inStructValue := reflect.ValueOf(inStruct)
	if inStructValue.Kind() != reflect.Ptr {
		return ErrNotPinter
	}

	inStructValue = inStructValue.Elem()
	if inStructValue.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	inStructType := inStructValue.Type()

	for i := 0; i < inStructType.NumField(); i++ {
		structItemField := inStructType.Field(i)
		structItemValue := inStructValue.Field(i)

		if val, ok := mapStructValues[structItemField.Name]; ok {
			mapItemValue := reflect.ValueOf(val)

			if mapItemValue.Type() != structItemValue.Type() {
				return fmt.Errorf("%w for field %s", ErrIncompatibleType, structItemField.Name)
			}

			if !structItemValue.CanSet() {
				return fmt.Errorf("%w for field %s", ErrCannotSetValue, structItemField.Name)
			}

			if structItemValue.Kind() != reflect.Struct {
				structItemValue.Set(mapItemValue)
				continue
			}

			err := manageInnerStruct(structItemValue.Addr().Interface(), mapItemValue.Interface())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func manageInnerStruct(in interface{}, mapStruct interface{}) error {
	reflectValue := reflect.ValueOf(in)

	if reflectValue.Kind() != reflect.Ptr {
		return ErrNotPinter
	}
	reflectValue = reflectValue.Elem()
	if reflectValue.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	reflectType := reflect.TypeOf(in).Elem()

	reflectMapValue := reflect.ValueOf(mapStruct)
	reflectMapType := reflect.TypeOf(mapStruct)
	if reflectMapValue.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	for i := 0; i < reflectType.NumField(); i++ {
		structItemField := reflectType.Field(i)
		structItemValue := reflectValue.Field(i)

		if mapStructField, ok := reflectMapType.FieldByName(structItemField.Name); ok {

			if mapStructField.Type != structItemValue.Type() {
				return fmt.Errorf("%w for field %s", ErrIncompatibleType, structItemField.Name)
			}

			if !structItemValue.CanSet() {
				return fmt.Errorf("%w for field %s", ErrCannotSetValue, structItemField.Name)
			}

			mapItemValue := reflectMapValue.FieldByName(structItemField.Name)
			if structItemValue.Kind() != reflect.Struct {
				structItemValue.Set(mapItemValue)
				continue
			}

			err := manageInnerStruct(structItemValue.Addr().Interface(), mapItemValue.Interface())
			if err != nil {
				return err
			}

		}
	}

	return nil
}
