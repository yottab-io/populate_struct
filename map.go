package populate_struct

import (
	"fmt"
	"reflect"
	"strings"
)

// MapToStruct function to populate struct fields from map values
func MapToStruct(resultPointer any, data map[string]string, escapePath string) error {
	var (
		esc    []string
		escLen int
	)

	if len(escapePath) > 0 {
		esc = strings.Split(escapePath, SPLIT_CHAR)
		escLen = len(esc)
	}

	for key, value := range data {
		fields := strings.Split(key, SPLIT_CHAR)

		// if escapePath && 'key' not have escape path: not process 'key'
		if len(escapePath) > 0 && !haveEscapePath(fields, esc) {
			continue
		}

		if err := fromPathAndValue(resultPointer, fields[escLen:], value); err != nil {
			return err
		}
	}

	return nil
}

func haveEscapePath(fields, esc []string) bool {
	if len(esc) >= len(fields) {
		return false
	}

	for i, e := range esc {
		if fields[i] != e {
			return false
		}
	}
	return true
}

// Function to set a field value using a path of JSON tags
func fromPathAndValue(obj any, fields []string, value string) error {
	// Get the final field to set the value
	finalField, err := findField(obj, fields)
	if err != nil {
		return err
	}

	// Set the value of the final field
	if !finalField.CanSet() {
		return fmt.Errorf("cannot set field with tag: %s", fields[len(fields)-1])
	}

	val, err := convertStringToType(value, finalField.Type())
	if err != nil {
		return err
	}

	finalField.Set(val)
	return nil
}

func findField(obj any, path []string) (reflect.Value, error) {
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem() // Dereference pointer
	}

	for _, field := range path {
		val = val.FieldByName(field)

		if !val.IsValid() {
			return val, fmt.Errorf("field %s not found", field)
		}

		if val.Kind() == reflect.Ptr {
			val = val.Elem() // Dereference pointer if the field is a pointer
		}
	}

	return val, nil
}

//////////////////////////////////////////////////////////////////////////////////////

func GetFieldValue(obj any, path string) (any, error) {
	// Split the path into fields
	fields := strings.Split(path, SPLIT_CHAR)

	val, err := findField(obj, fields)
	if err != nil {
		return nil, err
	}

	return val.Interface(), nil
}
