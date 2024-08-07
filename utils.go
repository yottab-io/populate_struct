package populate_struct

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrParameterNotAccessToSet = errors.New("")
	ErrNotHaveAccessKeys       = errors.New("")
)

// populateStructFromMap function to populate struct fields from map values
func FromMap(resultPointer interface{}, data map[string]interface{}, accessKey ...string) error {
	if len(accessKey) == 0 {
		return ErrNotHaveAccessKeys
	}
	accessMap := makeKeysMap(accessKey)

	for key, value := range data {
		path := strings.Split(key, ".")

		if _, exist := accessMap[path[0]]; !exist {
			return ErrParameterNotAccessToSet
		}

		if err := fromPathAndValue(resultPointer, path, value); err != nil {
			return err
		}
	}

	return nil
}

// Function to set a field value using a path of JSON tags
func fromPathAndValue(obj interface{}, path []string, value interface{}) (err error) {
	v := reflect.ValueOf(obj).Elem() // Get the value of the passed object

	for _, tag := range path[:len(path)-1] {
		v = findFieldByTag(v, tag)
		if !v.IsValid() {
			return fmt.Errorf("invalid field with tag: %s", tag)
		}
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
	}

	// Get the final field to set the value
	finalField := findFieldByTag(v, path[len(path)-1])
	if !finalField.IsValid() {
		return fmt.Errorf("invalid field with tag: %s", path[len(path)-1])
	}

	// Set the value of the final field
	if !finalField.CanSet() {
		return fmt.Errorf("cannot set field with tag: %s", path[len(path)-1])
	}

	val := reflect.ValueOf(value)
	if finalField.Type() != val.Type() {
		val, err = ConvertStringToType(val.String(), finalField.Type())
		if err != nil {
			return err
		}
	}

	finalField.Set(val)
	return nil
}

// Helper function to find a field by its JSON tag
func findFieldByTag(v reflect.Value, tag string) reflect.Value {
	v = reflect.Indirect(v)
	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)

		fieldTag := field.Tag.Get("json")
		if len(strings.Split(fieldTag, ",")) > 1 {
			fieldTag = strings.Split(fieldTag, ",")[0]
		}

		if fieldTag == tag {
			return v.Field(i)
		}
	}
	return reflect.Value{}
}

// Helper function to convert value types
func ConvertStringToType(val string, targetType reflect.Type) (reflect.Value, error) {
	switch targetType.Kind() {
	case reflect.Int:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("cannot convert %s to int: %v", val, err)
		}
		return reflect.ValueOf(int(floatVal)), nil
	case reflect.String:
		return reflect.ValueOf(val), nil
	case reflect.Float64:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("cannot convert %s to float64: %v", val, err)
		}
		return reflect.ValueOf(floatVal), nil
	case reflect.Bool:
		boolVal, err := strconv.ParseBool(val)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("cannot convert %s to bool: %v", val, err)
		}
		return reflect.ValueOf(boolVal), nil
	case reflect.Slice:
		if targetType.Elem().Kind() == reflect.String {
			// Assume comma-separated values for slice of strings
			jsonSlice := make([]string, 0)
			if err := json.Unmarshal([]byte(val), &jsonSlice); err != nil {
				return reflect.Value{}, err
			}
			return reflect.ValueOf(jsonSlice), nil
		}
	// Add more cases as needed
	default:
		return reflect.Value{}, fmt.Errorf("unsupported type: %s", targetType.Kind())
	}
	return reflect.Value{}, fmt.Errorf("unsupported type: %s", targetType.Kind())
}

func makeKeysMap(keys []string) map[string]any {
	keysMap := make(map[string]any, len(keys))
	for _, key := range keys {
		keysMap[key] = nil
	}
	return keysMap
}