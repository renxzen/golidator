package golidator

import (
	"encoding/json"
	"regexp"
	"testing"

	"github.com/renxzen/golidator/cmd/validator"
)

var LogErrors = false

func LogErrorsJson(t *testing.T, errors []ValidationError) {
	if !LogErrors {
		return
	}

	b, err := json.MarshalIndent(errors, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}

func TestNotBlank(t *testing.T) {
	type Request struct {
		Field1 string  `validate:"notblank"`
		Field2 *string `validate:"notblank"`
	}

	blankField := ""
	notBlankField := "testing"
	testTable := []struct {
		name   string
		input  Request
		output int
	}{
		{
			name: "Ok",
			input: Request{
				Field1: "testing",
				Field2: &notBlankField,
			},
			output: 0,
		},
		{
			name: "NotOk",
			input: Request{
				Field1: "",
				Field2: &blankField,
			},
			output: 2,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			errors, err := Validate(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			LogErrorsJson(t, errors)

			if len(errors) != tt.output {
				t.Errorf("\nExpected: %v.\nResult: %v.", 0, len(errors))
			}
		})
	}
}

func TestEmail(t *testing.T) {
	type Request struct {
		Field1 string  `validate:"email"`
		Field2 *string `validate:"email"`
	}

	type BadRequest struct {
		Field1 int      `validate:"email"`
		Field2 *float64 `validate:"email"`
	}

	emailOk := "renxzen@gmail.com"
	emailNotOk := "renxzen@g.n"
	invalidType := 420.0
	testTable := []struct {
		name     string
		input    Request
		badInput BadRequest
		output   int
	}{
		{
			name: "Ok",
			input: Request{
				Field1: "renxzen@gmail.com",
				Field2: &emailOk,
			},
			output: 0,
		},
		{
			name: "NotOk",
			input: Request{
				Field1: "renxzen@gmail",
				Field2: &emailNotOk,
			},
			output: 2,
		},
		{
			name: "InvalidType",
			badInput: BadRequest{
				Field1: 69.0,
				Field2: &invalidType,
			},
			output: 2,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			var errors []ValidationError
			var err error
			if tt.badInput.Field1 != 0 {
				errors, err = Validate(tt.badInput)
			} else {
				errors, err = Validate(tt.input)
			}

			if err != nil {
				t.Fatal(err)
			}

			LogErrorsJson(t, errors)

			if len(errors) != tt.output {
				t.Errorf("\nExpected: %v.\nResult: %v.", 0, len(errors))
			}

			// has errors
			message := ""
			if tt.badInput.Field1 != 0 {
				message = validator.NOTSTRING_ERROR
			} else {
				message = validator.EMAIL_ERROR
			}

			for _, error := range errors {
				for _, error := range error.Errors {
					if error != message {
						t.Errorf("Expected: %v. Result: %v", message, error)
					}
				}
			}
		})
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

	LogErrorsJson(t, errors)

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

	LogErrorsJson(t, errors)

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

	LogErrorsJson(t, errors)

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

	LogErrorsJson(t, errors)

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

	LogErrorsJson(t, errors)

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

	LogErrorsJson(t, errors)

	if len(errors) != 0 {
		t.Errorf("\nExpected: %v.\nResult: %v.", 0, len(errors))
	}
}

func TestMinError(t *testing.T) {
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

	LogErrorsJson(t, errors)

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
		Field1 bool  `validate:"min"`
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

	LogErrorsJson(t, errors)

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

func TestMaxIntOk(t *testing.T) {
	type Request struct {
		Field1 int  `validate:"max=5"`
		Field2 *int `validate:"max=5"`
	}

	field2 := 5
	input := Request{
		Field1: 4,
		Field2: &field2,
	}

	errors, err := Validate(input)
	if err != nil {
		t.Fatal(err)
	}

	LogErrorsJson(t, errors)

	if len(errors) != 0 {
		t.Errorf("\nExpected: %v.\nResult: %v.", 0, len(errors))
	}
}

func TestMaxFloatOk(t *testing.T) {
	type Request struct {
		Field1 float32  `validate:"max=5"`
		Field2 *float64 `validate:"max=5"`
	}

	field2 := 4.0
	input := Request{
		Field1: 4.99,
		Field2: &field2,
	}

	errors, err := Validate(input)
	if err != nil {
		t.Fatal(err)
	}

	LogErrorsJson(t, errors)

	if len(errors) != 0 {
		t.Errorf("\nExpected: %v.\nResult: %v.", 0, len(errors))
	}
}

func TestMaxStringOk(t *testing.T) {
	type Request struct {
		Field1 string  `validate:"max=5"`
		Field2 *string `validate:"max=5"`
	}

	field2 := "abcde"
	input := Request{
		Field1: "abcd",
		Field2: &field2,
	}

	errors, err := Validate(input)
	if err != nil {
		t.Fatal(err)
	}

	LogErrorsJson(t, errors)

	if len(errors) != 0 {
		t.Errorf("\nExpected: %v.\nResult: %v.", 0, len(errors))
	}
}

func TestMaxError(t *testing.T) {
	type Request struct {
		Field1 int      `validate:"max=5"`
		Field2 *int     `validate:"max=5"`
		Field3 float64  `validate:"max=5"`
		Field4 *float64 `validate:"max=5"`
		Field5 string   `validate:"max=5"`
		Field6 *string  `validate:"max=5"`
	}

	field2 := 6
	field4 := 5.99
	field6 := "abcdef"
	input := Request{
		Field1: 6,
		Field2: &field2,
		Field3: 5.01,
		Field4: &field4,
		Field5: "abdcef",
		Field6: &field6,
	}

	errors, err := Validate(input)
	if err != nil {
		t.Fatal(err)
	}

	LogErrorsJson(t, errors)

	if len(errors) != 6 {
		t.Errorf("\nExpected: %v.\nResult: %v.", 2, len(errors))
	}

	pattern := `^Must (have|be) less or equal than`
	re := regexp.MustCompile(pattern)
	for _, validationError := range errors {
		for _, error := range validationError.Errors {
			if !re.MatchString(error) {
				t.Errorf("Expected to match: %v. Result: %v", pattern, error)
			}
		}
	}
}

func TestMaxInvalidTypeError(t *testing.T) {
	type Request struct {
		Field1 bool  `validate:"max"`
		Field2 *bool `validate:"max"`
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

	LogErrorsJson(t, errors)

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

func TestNotemptyOk(t *testing.T) {
	type Request struct {
		Field1 []int  `validate:"notempty"`
		Field2 []*int `validate:"notempty"`
		Field3 *[]int `validate:"notempty"`
	}

	field2 := 1
	field3 := []int{1}
	input := Request{
		Field1: []int{1},
		Field2: []*int{&field2},
		Field3: &field3,
	}

	errors, err := Validate(input)
	if err != nil {
		t.Fatal(err)
	}

	LogErrorsJson(t, errors)

	if len(errors) != 0 {
		t.Errorf("\nExpected: %v.\nResult: %v.", 0, len(errors))
	}
}

func TestNotemptyError(t *testing.T) {
	type Request struct {
		Field1 []int  `validate:"notempty"`
		Field2 []*int `validate:"notempty"`
		Field3 *[]int `validate:"notempty"`
	}

	field3 := []int{}
	input := Request{
		Field1: []int{},
		Field2: []*int{},
		Field3: &field3,
	}

	errors, err := Validate(input)
	if err != nil {
		t.Fatal(err)
	}

	LogErrorsJson(t, errors)

	if len(errors) != 3 {
		t.Errorf("\nExpected: %v.\nResult: %v.", 3, len(errors))
	}

	message := validator.NOTEMPTY_ERROR
	for _, validationError := range errors {
		for _, error := range validationError.Errors {
			if error != message {
				t.Errorf("Expected: %v. Result: %v", message, error)
			}
		}
	}
}

func TestNotemptyInvalidTypeError(t *testing.T) {
	type Request struct {
		Field1 int      `validate:"notempty"`
		Field2 *float64 `validate:"notempty"`
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

	LogErrorsJson(t, errors)

	if len(errors) != 2 {
		t.Errorf("\nExpected: %v.\nResult: %v.", 2, len(errors))
	}

	message := validator.NOTARRAY_ERROR
	for _, validationError := range errors {
		for _, error := range validationError.Errors {
			if error != message {
				t.Errorf("Expected: %v. Result: %v", message, error)
			}
		}
	}
}

func TestArrayOk(t *testing.T) {
	type SubRequest struct {
		Field1 string `validate:"notblank,email"`
		Field2 string `validate:"notblank,max=20"`
	}

	type Request struct {
		Array []SubRequest `validate:"isarray"`
	}

	input := Request{
		Array: []SubRequest{
			{
				Field1: "renzo@gmail.com",
				Field2: "Renzo Mondragón",
			},
			{
				Field1: "renato@gmail.com",
				Field2: "Renato Mondragón",
			},
			{
				Field1: "mya@gmail.com",
				Field2: "Mya García",
			},
		},
	}

	errors, err := Validate(input)
	if err != nil {
		t.Fatal(err)
	}

	LogErrorsJson(t, errors)

	if len(errors) != 0 {
		t.Errorf("\nExpected: %v.\nResult: %v.", 0, len(errors))
	}
}

func TestArrayErrors(t *testing.T) {
	type SubRequest struct {
		Field1 string `validate:"notblank,email"`
		Field2 string `validate:"notblank,max=20"`
	}

	type Request struct {
		Array []SubRequest `validate:"isarray"`
	}

	input := Request{
		Array: []SubRequest{
			{
				Field1: "renzo@gmail",
				Field2: "Renzo Mondragón Arango",
			},
			{
				Field1: "renatogmail.com",
				Field2: "Renato Mondragón",
			},
			{
				Field1: "@gmail.com",
				Field2: "",
			},
		},
	}

	errors, err := Validate(input)
	if err != nil {
		t.Fatal(err)
	}

	LogErrorsJson(t, errors)

	if len(errors) != 5 {
		t.Errorf("\nExpected: %v.\nResult: %v.", 0, len(errors))
	}
}
