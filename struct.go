package populate_struct

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

// StructToMap converts a struct to a map[string]string using the struct's JSON tags
func StructToMap(data any, prefix ...string) map[string]string {
	result := make(map[string]string)

	value := reflect.ValueOf(data)

	// Handle pointers by dereferencing
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	// Ensure the value is a struct
	if value.Kind() != reflect.Struct {
		panic(fmt.Sprintf("StructToMap only accepts structs or pointers to structs; got %T", data))
	}

	t := value.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := value.Field(i)

		// Get the json tag name
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name
		} else if len(strings.Split(jsonTag, ",")) > 1 {
			jsonTag = strings.Split(jsonTag, ",")[0]
		}

		// Construct the full key using the prefix and json tag
		fullKey := jsonTag
		if len(prefix) > 0 {
			fullKey = prefix[0] + "." + jsonTag
		}

		switch fieldValue.Kind() {
		case reflect.Struct:
			// Recursively process nested structs
			nestedMap := StructToMap(fieldValue.Interface(), fullKey)
			for k, v := range nestedMap {
				result[k] = v
			}
		case reflect.Ptr:
			// Handle pointer to a struct or a value
			if !fieldValue.IsNil() {
				dereferencedValue := fieldValue.Elem()
				if dereferencedValue.Kind() == reflect.Struct {
					nestedMap := StructToMap(dereferencedValue.Interface(), fullKey)
					for k, v := range nestedMap {
						result[k] = v
					}
				} else {
					result[fullKey] = fmt.Sprintf("%v", dereferencedValue.Interface())
				}
			}
		case reflect.Interface:
			// Handle interface fields by processing their concrete value
			if !fieldValue.IsNil() {
				interfaceValue := fieldValue.Interface()

				// Handle the case where the interface holds a struct (like helmValueRedis)
				nestedMap := StructToMap(interfaceValue, fullKey)
				for k, v := range nestedMap {
					result[k] = v
				}
			}
		case reflect.Slice:
			// Handle slices/arrays
			result[fullKey] = fmt.Sprintf("%q", fieldValue)
			// for j := 0; j < fieldValue.Len(); j++ {
			// 	item := fieldValue.Index(j)
			// 	result[fmt.Sprintf("%s[%d]", fullKey, j)] = fmt.Sprintf("%v", item.Interface())
			// }

		case reflect.Bool:
			// Convert boolean to string
			result[fullKey] = strconv.FormatBool(fieldValue.Bool())
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			// Convert int to string
			result[fullKey] = strconv.FormatInt(fieldValue.Int(), 10)
		case reflect.String:
			// Convert string to string (directly assign)
			result[fullKey] = fieldValue.String()
		default:
			// Handle other kinds of fields by converting them to string
			log.Println("WARN: StructToMap unsupported type")
			result[fullKey] = fmt.Sprintf("%v", fieldValue.Interface())
		}
	}

	return result
}