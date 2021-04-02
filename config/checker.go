package config

import (
	"fmt"
	"reflect"

	"github.com/shitpostingio/autopostingbot/config/structs"
)

// checkMandatoryFields uses reflection to see if there are
// mandatory fields with zero value.
// When using with isReload flag, only checks for fields with reloadable:"true" tag
func checkMandatoryFields(isReload bool, config structs.Config) error {
	resultingError := ""
	for _, err := range checkStruct(isReload, "root", reflect.TypeOf(config), reflect.ValueOf(config)) {
		resultingError = fmt.Sprintf("%s\n%s", resultingError, err.Error())
	}

	if resultingError == "" {
		return nil
	}

	return fmt.Errorf("%s", resultingError)
}

// checkStruct explores structures recursively and checks if
// struct fields have a zero value.
func checkStruct(isReload bool, parent string, typeToCheck reflect.Type, valueToCheck reflect.Value) (errors []error) {
	errors = make([]error, 0)
	for i := 0; i < typeToCheck.NumField(); i++ {
		currentField := typeToCheck.Field(i)
		currentValue := valueToCheck.Field(i)
		switch currentField.Type.Kind() {
		case
			reflect.Uintptr,
			reflect.UnsafePointer,
			reflect.Ptr:
			continue

		case reflect.Struct:
			if errs := checkStruct(isReload, currentField.Name, currentField.Type, currentValue); errs != nil {
				errors = append(errors, errs...)
			}

		case reflect.Slice:
			if errs := checkSlice(isReload, fmt.Sprintf("%s in %s", currentField.Name, parent), currentField, currentValue); errs != nil {
				errors = append(errors, errs...)
			}

		default:
			if err := checkField(isReload, parent, currentField, currentValue); err != nil {
				errors = append(errors, err)
			}
		}
	}

	if len(errors) != 0 {
		return errors
	}

	return nil
}

// checkSlice explores slices recursively and checks they have a zero value.
func checkSlice(isReload bool, parent string, typeToCheck reflect.StructField, sliceToCheck reflect.Value) []error {

	//only check reloadable fields if isReload is true
	if isReload {

		reloadableTagValue := typeToCheck.Tag.Get("reloadable")
		if reloadableTagValue != "true" {
			return nil
		}

	}

	typeTagValue := typeToCheck.Tag.Get("type")
	if typeTagValue == "optional" {
		return nil
	}

	if sliceToCheck.Len() == 0 {
		return []error{fmt.Errorf("non optional slice field %s in %s had zero length", typeToCheck.Name, parent)}
	}

	errors := make([]error, 0)
	for i := 0; i < sliceToCheck.Len(); i++ {
		item := sliceToCheck.Index(i)
		if item.Kind() == reflect.Struct {
			if errs := checkStruct(isReload, parent, reflect.TypeOf(item.Interface()), reflect.ValueOf(item.Interface())); errs != nil {
				errors = append(errors, errs...)
			}
			continue
		}

		zeroValue := reflect.Zero(item.Type())
		if item.Interface() == zeroValue.Interface() {
			errors = append(errors, fmt.Errorf("non optional field %s had zero value at index %d", typeToCheck.Name, i))
		}
	}

	if len(errors) != 0 {
		return errors
	}

	return nil
}

// checkField checks if a field is optional or a webhook field
// if it isn't, it checks if the field has a zero value.
func checkField(isReload bool, section string, typeToCheck reflect.StructField, valueToCheck reflect.Value) error {
	//only check reloadable fields if isReload is true
	if isReload {

		reloadableTagValue := typeToCheck.Tag.Get("reloadable")
		if reloadableTagValue != "true" {
			return nil
		}

	}

	typeTagValue := typeToCheck.Tag.Get("type")

	if typeTagValue == "optional" || typeTagValue == "webhook" {
		return nil
	}

	zeroValue := reflect.Zero(typeToCheck.Type)

	if valueToCheck.Interface() == zeroValue.Interface() {
		return fmt.Errorf("non optional field %s in section %s had zero value", typeToCheck.Name, section)
	}

	return nil
}
