package main

import "testing"

func TestValidateJson(t *testing.T) {
	someString := `{"name": "Mark", "age": 12}`

	want := validateJson(someString)

	if want != true {
		t.Errorf("Invalid JSON passed to this funciton")
	}
}
