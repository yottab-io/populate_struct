package populate_struct

import (
	"testing"
)

func TestStruct_BoolField(t *testing.T) {
	instance := &StructA{}
	instance.FieldBool = true
	data := StructToMap(instance)

	expected := "true"
	expectedKey := "StructX.StructY.StructZ.FieldBool"
	if out, exist := data[expectedKey]; !exist {
		t.Errorf("Expected key %s, but not exist", expectedKey)
	} else if out != expected {
		t.Errorf("Expected Flag to be %s, but got %s", out, expected)
	}
}
