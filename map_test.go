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

	StructP struct {
		FieldP *StructZ
	}

	StructZ struct {
		FieldBool        bool
		FieldString      string
		FieldFloat       float64
		Nested           NestedStruct
		FieldArrayString []string
	}

	NestedStruct struct {
		NestedField int32
	}
)

func TestSetBoolField(t *testing.T) {
	instance := &StructA{}
	err := MapToStruct(instance, map[string]string{"FieldBool": "true"}, ClearEscapePath)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := true
	if instance.FieldBool != expected {
		t.Errorf("Expected Flag to be %v, but got %v", expected, instance.FieldBool)
	}
}

func TestSetPointerField(t *testing.T) {
	instance := &StructP{}
	instance.FieldP = new(StructZ)
	err := MapToStruct(instance, map[string]string{"FieldP.FieldBool": "true"}, ClearEscapePath)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := true
	if instance.FieldP.FieldBool != expected {
		t.Errorf("Expected Flag to be %v, but got %v", expected, instance.FieldP.FieldBool)
	}
}

func TestSetIntField(t *testing.T) {
	instance := &StructA{}

	err := MapToStruct(instance, map[string]string{"Nested.NestedField": "2"}, ClearEscapePath)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	var expected int32 = 2
	if instance.Nested.NestedField != expected {
		t.Errorf("Expected Nested.NestedField to be %d, but got %d", expected, instance.Nested.NestedField)
	}

	//////////////////////////////////////////////////////////////////////////////////

	instance = &StructA{}
	err = MapToStruct(instance, map[string]string{"esc.Nested.NestedField": "2"}, "esc")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected = 2
	if instance.Nested.NestedField != expected {
		t.Errorf("Expected Nested.NestedField to be %d, but got %d", expected, instance.Nested.NestedField)
	}

	//////////////////////////////////////////////////////////////////////////////////

	instance = &StructA{}
	err = MapToStruct(instance, map[string]string{"Nested.NestedField": "-2.2"}, ClearEscapePath)
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

	err := MapToStruct(instance, map[string]string{"test.esc.Nested.NestedField": "2"}, "test.esc")
	if err != nil {
		t.Errorf("Error: %v", err)
	}
	var expected int32 = 2
	if instance.Nested.NestedField != expected {
		t.Errorf("Expected Nested.NestedField to be %d, but got %d", expected, instance.Nested.NestedField)
	}

	// test no set
	instance = &StructA{}
	err = MapToStruct(instance, map[string]string{"Nested.NestedField": "2"}, "test.esc")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected = 0
	if instance.Nested.NestedField != expected {
		t.Errorf("Expected Nested.NestedField to be %d, but got %d", expected, instance.Nested.NestedField)
	}

	// test set to Non Existent Field must be scape
	instance = &StructA{}
	err = MapToStruct(instance, map[string]string{"Nested.NonExistent": "2"}, ClearEscapePath)
	if err != nil {
		t.Errorf("Expected no error for not have access key, but got err: %s", err.Error())
	}
}

// Test for setting string fields
func TestSetStringField(t *testing.T) {
	instance := &StructA{}

	err := MapToStruct(instance, map[string]string{"FieldString": "updated"}, ClearEscapePath)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := "updated"
	if instance.FieldString != expected {
		t.Errorf("Expected Field2 to be %s, but got %s", expected, instance.FieldString)
	}
}

// Test for setting slice fields
func TestSetSliceField(t *testing.T) {
	instance := &StructA{}

	err := MapToStruct(instance, map[string]string{"FieldArrayString": `["a", "b", "c"]`}, ClearEscapePath)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := []string{"a", "b", "c"}
	if !reflect.DeepEqual(instance.FieldArrayString, expected) {
		t.Errorf("Expected Strings to be %v, but got %v", expected, instance.FieldArrayString)
	}

	////////////////////

	err = MapToStruct(instance, map[string]string{"FieldArrayString": `["a, "b", "c"]`}, ClearEscapePath)
	if err == nil {
		t.Errorf("Expected an error for json format, but got none")
	}
}

func TestSetFloatField(t *testing.T) {
	instance := &StructA{}
	err := MapToStruct(instance, map[string]string{"FieldFloat": "12.34"}, ClearEscapePath)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	expected := 12.34
	if instance.FieldFloat != expected {
		t.Errorf("Expected Amount to be %v, but got %v", expected, instance.FieldFloat)
	}
}

//////////////////////////////////////////////////////////////////////////////////////

func TestMap_GetFloatField(t *testing.T) {
	instance := &StructA{}
	instance.FieldFloat = 12.34
	val, err := GetFieldValue(instance, "FieldFloat")
	if err != nil {
		t.Errorf("Error: %s", err.Error())
	}

	expected := 12.34
	if val.(float64) != expected {
		t.Errorf("Expected Amount to be %f, but got %v", expected, val)
	}
}
