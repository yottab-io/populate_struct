package populate_struct

import (
	"reflect"
	"testing"
)

// Define StructA with nested structures and a slice field
type (
	StructA struct {
		StructX
	}

	StructX struct {
		StructY
	}

	StructY struct {
		StructZ
	}

	StructZ struct {
		Field0  bool         `json:"field0"`
		Field1  string       `json:"field1"`
		Field2  float64      `json:"field2"`
		Nested  NestedStruct `json:"nested"`
		Strings []any        `json:"strings"`
	}

	NestedStruct struct {
		NestedField int32 `json:"field0"`
	}
)

func TestSetBoolField(t *testing.T) {
	instance := &StructA{}
	err := FromMap(instance, ColonSplitChar, map[string]any{"field0": "true"}, "field0")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := true
	if instance.Field0 != expected {
		t.Errorf("Expected Flag to be %v, but got %v", expected, instance.Field0)
	}

	//////////////////////////////////////////////////////////////////////////////////

	instance = &StructA{}
	err = FromMapString(instance, ColonSplitChar, map[string]string{"field0": "true"}, "field0")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if instance.Field0 != expected {
		t.Errorf("Expected Flag to be %v, but got %v", expected, instance.Field0)
	}

	//////////////////////////////////////////////////////////////////////////////////

	instance = &StructA{}
	err = FromMap(instance, ColonSplitChar, map[string]any{"field0": true}, "field0")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if instance.Field0 != expected {
		t.Errorf("Expected Flag to be %v, but got %v", expected, instance.Field0)
	}
}

func TestSetIntField(t *testing.T) {
	instance := &StructA{}

	err := FromMap(instance, ColonSplitChar, map[string]any{"nested:field0": 2}, "nested")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	var expected int32 = 2
	if instance.Nested.NestedField != expected {
		t.Errorf("Expected Nested.NestedField to be %d, but got %d", expected, instance.Nested.NestedField)
	}

	//////////////////////////////////////////////////////////////////////////////////
	instance = &StructA{}
	err = FromMap(instance, ColonSplitChar, map[string]any{"nested:field0": "2"}, "nested")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if instance.Nested.NestedField != expected {
		t.Errorf("Expected Nested.NestedField to be %d, but got %d", expected, instance.Nested.NestedField)
	}

	//////////////////////////////////////////////////////////////////////////////////
	instance = &StructA{}
	err = FromMap(instance, ColonSplitChar, map[string]any{"nested:field0": "-2.2"}, "nested")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected = -2
	if instance.Nested.NestedField != expected {
		t.Errorf("Expected Nested.NestedField to be %d, but got %d", expected, instance.Nested.NestedField)
	}
}

func TestPathField(t *testing.T) {
	instance := &StructA{}

	err := FromMap(instance, ColonSplitChar, map[string]any{"nested:field0": 2})
	if err == nil {
		t.Errorf("Expected an error for not have access key, but got none")
	}

	err = FromMap(instance, ColonSplitChar, map[string]any{"nested:field0": 2}, "nonexistent")
	if err == nil {
		t.Errorf("Expected an error for not have access key, but got none")
	}

	err = FromMap(instance, ColonSplitChar, map[string]any{"nested:nonexistent": 2}, "nested")
	if err == nil {
		t.Errorf("Expected an error for not have access key, but got none")
	}
}

// Test for setting string fields
func TestSetStringField(t *testing.T) {
	instance := &StructA{}

	err := FromMap(instance, ColonSplitChar, map[string]any{"field1": "updated"}, "field1")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := "updated"
	if instance.Field1 != expected {
		t.Errorf("Expected Field2 to be %s, but got %s", expected, instance.Field1)
	}
}

// Test for setting slice fields
func TestSetSliceField(t *testing.T) {
	instance := &StructA{}

	err := FromMap(instance, ColonSplitChar, map[string]any{"strings": `["a", "b", "c"]`}, "strings")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := []any{"a", "b", "c"}
	if !reflect.DeepEqual(instance.Strings, expected) {
		t.Errorf("Expected Strings to be %v, but got %v", expected, instance.Strings)
	}
}

func TestSetFloatField(t *testing.T) {
	instance := &StructA{}
	err := FromMap(instance, ColonSplitChar, map[string]any{"field2": "12.34"}, "field2")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := 12.34
	if instance.Field2 != expected {
		t.Errorf("Expected Amount to be %v, but got %v", expected, instance.Field2)
	}
}

//////////////////////////////////////////////////////////////////////////////////////

func TestMapStrAnyToMapStrStr(t *testing.T) {
	input := map[string]interface{}{
		"name":      "Alice",
		"age":       25,
		"number":    8.2,
		"isStudent": true,
		"hobbies":   []any{"reading", "swimming"},
	}

	output := map[string]string{
		"name":      "Alice",
		"age":       "25",
		"number":    "8.2",
		"isStudent": "true",
		"hobbies":   `["reading","swimming"]`,
	}

	result := MapStrAnyToMapStrStr(input)
	if !reflect.DeepEqual(result, output) {
		t.Errorf("convertMap(%v) = %v; want %v", input, result, output)
	}
}
