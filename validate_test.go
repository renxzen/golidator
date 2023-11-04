package golidator

import (
	"encoding/json"
	"testing"

	"github.com/renxzen/golidator/cmd/validator"
)

func TestNotBlankNoErrors(t *testing.T) {
	type Request struct {
		Field1 string  `validate:"notblank"`
		Field2 *string `validate:"notblank"`
	}

	field2 := "testing"
	input := Request{
		Field1: "testing",
		Field2: &field2,
	}

	errors, err := Validate(input)
	if err != nil {
		t.Fatal(err)
	}

	b, err := json.MarshalIndent(errors, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))

	if len(errors) != 0 {
		t.Errorf("\nExpected: %v.\nResult: %v.", 0, len(errors))
	}
}

func TestNotBlankWithErrors(t *testing.T) {
	type Request struct {
		Field1 string  `validate:"notblank"`
		Field2 *string `validate:"notblank"`
	}

	field2 := ""
	input := Request{
		Field1: "",
		Field2: &field2,
	}

	errors, err := Validate(input)
	if err != nil {
		t.Fatal(err)
	}

	b, err := json.MarshalIndent(errors, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))

	if len(errors) != 2 {
		t.Errorf("\nExpected: %v.\nResult: %v.", 2, len(errors))
	}

	message := validator.NOTBLANK_ERROR
	for _, validationError := range errors {
		for _, error := range validationError.Errors {
			if error != message {
				t.Errorf("Expected: %v. Result: %v", message, error)
			}
		}
	}
}
