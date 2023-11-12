package golidator

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/renxzen/golidator/cmd/validator"
)

func TestNotBlankOk(t *testing.T) {
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

func TestEmailOk(t *testing.T) {
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
		Field1 int      `validate:"email"`
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

func TestUrlOk(t *testing.T) {
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
		Field1 int      `validate:"url"`
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

func TestMinIntOk(t *testing.T) {
	type Request struct {
		Field1 int  `validate:"min=5"`
		Field2 *int `validate:"min=5"`
	}

	field2 := 5
	input := Request{
		Field1: 6,
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

func TestMinFloatOk(t *testing.T) {
	type Request struct {
		Field1 float32  `validate:"min=5"`
		Field2 *float64 `validate:"min=5"`
	}

	field2 := 5.1
	input := Request{
		Field1: 5.01,
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

func TestMinStringOk(t *testing.T) {
	type Request struct {
		Field1 string  `validate:"min=5"`
		Field2 *string `validate:"min=5"`
	}

	field2 := "abcde"
	input := Request{
		Field1: "abcdef",
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

func TestMinStringError(t *testing.T) {
	type Request struct {
		Field1 int      `validate:"min=5"`
		Field2 *int     `validate:"min=5"`
		Field3 float64  `validate:"min=5"`
		Field4 *float64 `validate:"min=5"`
		Field5 string   `validate:"min=5"`
		Field6 *string  `validate:"min=5"`
	}

	field2 := 4
	field4 := 4.99
	field6 := "abcd"
	input := Request{
		Field1: 4,
		Field2: &field2,
		Field3: 4.99,
		Field4: &field4,
		Field5: "abdc",
		Field6: &field6,
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

	if len(errors) != 6 {
		t.Errorf("\nExpected: %v.\nResult: %v.", 2, len(errors))
	}

	pattern := `^Must (have|be) more or equal than`
	re := regexp.MustCompile(pattern)
	for _, validationError := range errors {
		for _, error := range validationError.Errors {
			if !re.MatchString(error) {
				t.Errorf("Expected to match: %v. Result: %v", pattern, error)
			}
		}
	}
}

func TestMinInvalidTypeError(t *testing.T) {
	type Request struct {
		Field1 bool `validate:"min"`
		Field2 *bool `validate:"min"`
	}

	field2 := true
	input := Request{
		Field1: false,
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

	message := validator.NOTSTRINGORNUMERIC_ERROR
	for _, validationError := range errors {
		for _, error := range validationError.Errors {
			if error != message {
				t.Errorf("Expected: %v. Result: %v", message, error)
			}
		}
	}
}
