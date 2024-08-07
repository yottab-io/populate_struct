package populate_struct

import (
	"reflect"
	"testing"
)

// Define StructA with nested structures and a slice field
type StructA struct {
	Field0  bool         `json:"field0"`
	Field1  int          `json:"field1"`
	Field2  string       `json:"field2"`
	Field3  float64      `json:"field3"`
	Nested  NestedStruct `json:"nested"`
	Strings []string     `json:"strings"`
}
type NestedStruct struct {
	NestedField int `json:"field0"`
}

func TestSetBoolField(t *testing.T) {
	instance := &StructA{}
	err := FromMap(instance, map[string]any{"field0": "true"}, "field0")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := true
	if instance.Field0 != expected {
		t.Errorf("Expected Flag to be %v, but got %v", expected, instance.Field0)
	}

	//////////////////////////////////////////////////////////////////////////////////

	instance = &StructA{}
	err = FromMapString(instance, map[string]string{"field0": "true"}, "field0")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if instance.Field0 != expected {
		t.Errorf("Expected Flag to be %v, but got %v", expected, instance.Field0)
	}

	//////////////////////////////////////////////////////////////////////////////////

	instance = &StructA{}
	err = FromMap(instance, map[string]any{"field0": true}, "field0")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if instance.Field0 != expected {
		t.Errorf("Expected Flag to be %v, but got %v", expected, instance.Field0)
	}
}

func TestSetIntField(t *testing.T) {
	instance := &StructA{}

	err := FromMap(instance, map[string]any{"nested.field0": 2}, "nested")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := 2
	if instance.Nested.NestedField != expected {
		t.Errorf("Expected Nested.NestedField to be %d, but got %d", expected, instance.Nested.NestedField)
	}

	//////////////////////////////////////////////////////////////////////////////////
	instance = &StructA{}
	err = FromMap(instance, map[string]any{"nested.field0": "2"}, "nested")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if instance.Nested.NestedField != expected {
		t.Errorf("Expected Nested.NestedField to be %d, but got %d", expected, instance.Nested.NestedField)
	}

	//////////////////////////////////////////////////////////////////////////////////
	instance = &StructA{}
	err = FromMap(instance, map[string]any{"nested.field0": "-2.2"}, "nested")
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

	err := FromMap(instance, map[string]any{"nested.field0": 2})
	if err == nil {
		t.Errorf("Expected an error for not have access key, but got none")
	}

	err = FromMap(instance, map[string]any{"nested.field0": 2}, "nonexistent")
	if err == nil {
		t.Errorf("Expected an error for not have access key, but got none")
	}

	err = FromMap(instance, map[string]any{"nested.nonexistent": 2}, "nested")
	if err == nil {
		t.Errorf("Expected an error for not have access key, but got none")
	}
}

// Test for setting string fields
func TestSetStringField(t *testing.T) {
	instance := &StructA{}

	err := FromMap(instance, map[string]any{"field2": "updated"}, "field2")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := "updated"
	if instance.Field2 != expected {
		t.Errorf("Expected Field2 to be %s, but got %s", expected, instance.Field2)
	}
}

// Test for setting slice fields
func TestSetSliceField(t *testing.T) {
	instance := &StructA{}

	err := FromMap(instance, map[string]any{"strings": `["a", "b", "c"]`}, "strings")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(instance.Strings, expected) {
		t.Errorf("Expected Strings to be %v, but got %v", expected, instance.Strings)
	}
}

func TestSetFloatField(t *testing.T) {
	instance := &StructA{}
	err := FromMap(instance, map[string]any{"field3": "12.34"}, "field3")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := 12.34
	if instance.Field3 != expected {
		t.Errorf("Expected Amount to be %v, but got %v", expected, instance.Field3)
	}
}

//////////////////////////////////////////////////////////////////////////////////////

func TestMapStrAnyToMapStrStr(t *testing.T) {
	input := map[string]interface{}{
		"name":      "Alice",
		"age":       25,
		"number":    8.2,
		"isStudent": true,
		"hobbies":   []string{"reading", "swimming"},
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
