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

func TestEmailNoErrors(t *testing.T) {
	type Request struct {
		Field1 string  `validate:"email"`
		Field2 *string `validate:"email"`
	}

	field2 := "renxzen@gmail.com"
	input := Request{
		Field1: "renxzen@gmail.com",
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

func TestEmailInvalidValueError(t *testing.T) {
	type Request struct {
		Field1 string  `validate:"email"`
		Field2 *string `validate:"email"`
	}

	field2 := "renxzen@g.n"
	input := Request{
		Field1: "renxzen@gmail",
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

	message := validator.EMAIL_ERROR
	for _, validationError := range errors {
		for _, error := range validationError.Errors {
			if error != message {
				t.Errorf("Expected: %v. Result: %v", message, error)
			}
		}
	}
}

func TestEmailInvalidTypeError(t *testing.T) {
	type Request struct {
		Field1 int  `validate:"email"`
		Field2 *float64 `validate:"email"`
	}

	field2 := 420.0
	input := Request{
		Field1: 69.0,
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

	message := validator.NOTSTRING_ERROR
	for _, validationError := range errors {
		for _, error := range validationError.Errors {
			if error != message {
				t.Errorf("Expected: %v. Result: %v", message, error)
			}
		}
	}
}

func TestUrlNoErrors(t *testing.T) {
	type Request struct {
		Field1 string  `validate:"url"`
		Field2 *string `validate:"url"`
	}

	field2 := "http://renxzen.com"
	input := Request{
		Field1: "https://renxzen.com",
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

func TestUrlInvalidValueError(t *testing.T) {
	type Request struct {
		Field1 string  `validate:"url"`
		Field2 *string `validate:"url"`
	}

	field2 := "renxzen.com"
	input := Request{
		Field1: "www.renxzen",
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

	message := validator.URL_ERROR
	for _, validationError := range errors {
		for _, error := range validationError.Errors {
			if error != message {
				t.Errorf("Expected: %v. Result: %v", message, error)
			}
		}
	}
}

func TestUrlInvalidTypeError(t *testing.T) {
	type Request struct {
		Field1 int  `validate:"url"`
		Field2 *float64 `validate:"url"`
	}

	field2 := 420.0
	input := Request{
		Field1: 69.0,
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

	message := validator.NOTSTRING_ERROR
	for _, validationError := range errors {
		for _, error := range validationError.Errors {
			if error != message {
				t.Errorf("Expected: %v. Result: %v", message, error)
			}
		}
	}
}
