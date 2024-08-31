package populate_struct

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

const (
	SPLIT_CHAR      = "."
	ClearEscapePath = ""
)

var (
	ErrParameterNotAccessToSet = errors.New("error Parameter Not Access To Set")
	ErrNotHaveAccessKeys       = errors.New("error Not Have Access Keys")
	ErrFieldNotFound           = errors.New("error field not found")
)

// Helper function to convert value types
func convertStringToType(val string, targetType reflect.Type) (reflect.Value, error) {
	switch targetType.Kind() {
	case reflect.String:
		return reflect.ValueOf(val), nil
	case reflect.Int:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("cannot convert %s to int: %v", val, err)
		}
		return reflect.ValueOf(int(floatVal)), nil
	case reflect.Int32:
		floatVal, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("cannot convert %s to int32: %v", val, err)
		}
		return reflect.ValueOf(int32(floatVal)), nil
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
		// Assume comma-separated values for slice of strings
		jsonSlice := make([]string, 0)
		if err := json.Unmarshal([]byte(val), &jsonSlice); err != nil {
			return reflect.Value{}, err
		}
		return reflect.ValueOf(jsonSlice), nil

	// Add more cases as needed
	default:
		return reflect.Value{}, fmt.Errorf("unsupported type: %s", targetType.Kind())
	}
}

func JsonReplaceInterface(s, d any) error {
	if bytes, err := json.Marshal(s); err != nil {
		return err
	} else if err = json.Unmarshal(bytes, d); err != nil {
		return err
	}

	return nil
}

func ConvertStructToMapStringAny(obj any) (map[string]any, error) {
	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}

	mapStrToAny := make(map[string]any, 0)
	if err := json.Unmarshal(jsonObj, &mapStrToAny); err != nil {
		return nil, err
	}

	return mapStrToAny, nil
}
