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

			if len(errors) == 0 {
				return
			}

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

func TestUrl(t *testing.T) {
	type Request struct {
		Field1 string  `validate:"url"`
		Field2 *string `validate:"url"`
	}

	type BadRequest struct {
		Field1 int      `validate:"url"`
		Field2 *float64 `validate:"url"`
	}

	urlOk := "http://renxzen.com"
	urlNotOk := "renxzen.com"
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
				Field1: "https://renxzen.com",
				Field2: &urlOk,
			},
			output: 0,
		},
		{
			name: "NotOk",
			input: Request{
				Field1: "www.renxzen",
				Field2: &urlNotOk,
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

			if len(errors) == 0 {
				return
			}

			message := ""
			if tt.badInput.Field1 != 0 {
				message = validator.NOTSTRING_ERROR
			} else {
				message = validator.URL_ERROR
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

func TestMin(t *testing.T) {
	type Request[T int | float64 | string | bool] struct {
		Field1 T  `validate:"min=5"`
		Field2 *T `validate:"min=5"`
	}

	minIntOk := 5
	minIntNotOk := 4
	minFloatOk := 5.001
	minFloatNotOk := 4.999
	minStringOk := "abcde"
	minStringNotOk := "abcd"
	invalidType := true

	testTable := []struct {
		name         string
		intInput     Request[int]
		floatInput   Request[float64]
		stringInput  Request[string]
		invalidInput Request[bool]
		output       int
	}{
		{
			name: "IntOk",
			intInput: Request[int]{
				Field1: 6,
				Field2: &minIntOk,
			},
			output: 0,
		},
		{
			name: "IntNotOk",
			intInput: Request[int]{
				Field1: 4,
				Field2: &minIntNotOk,
			},
			output: 2,
		},
		{
			name: "FloatOk",
			floatInput: Request[float64]{
				Field1: 6,
				Field2: &minFloatOk,
			},
			output: 0,
		},
		{
			name: "FloatNotOk",
			floatInput: Request[float64]{
				Field1: 4,
				Field2: &minFloatNotOk,
			},
			output: 2,
		},
		{
			name: "StringOk",
			stringInput: Request[string]{
				Field1: "abcdef",
				Field2: &minStringOk,
			},
			output: 0,
		},
		{
			name: "StringNotOk",
			stringInput: Request[string]{
				Field1: "abcd",
				Field2: &minStringNotOk,
			},
			output: 2,
		},
		{
			name: "InvalidType",
			invalidInput: Request[bool]{
				Field1: false,
				Field2: &invalidType,
			},
			output: 2,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			var errors []ValidationError
			var err error
			isInvalidTest := false

			if tt.intInput.Field1 != 0 {
				errors, err = Validate(tt.intInput)
			} else if tt.floatInput.Field1 != 0 {
				errors, err = Validate(tt.floatInput)
			} else if tt.stringInput.Field1 != "" {
				errors, err = Validate(tt.stringInput)
			} else {
				errors, err = Validate(tt.invalidInput)
				isInvalidTest = true
			}

			if err != nil {
				t.Fatal(err)
			}

			LogErrorsJson(t, errors)

			if len(errors) != tt.output {
				t.Errorf("\nExpected: %v.\nResult: %v.", 0, len(errors))
			}

			if len(errors) == 0 {
				return
			}

			invalidTypeMessage := validator.NOTSTRINGORNUMERIC_ERROR
			pattern := `^Must (have|be) more or equal than`
			re := regexp.MustCompile(pattern)
			for _, validationError := range errors {
				for _, error := range validationError.Errors {
					if isInvalidTest && error != invalidTypeMessage {
						t.Errorf("Expected: %v. Result: %v", invalidTypeMessage, error)
					}

					if !isInvalidTest && !re.MatchString(error) {
						t.Errorf("Expected to match: %v. Result: %v", pattern, error)
					}
				}
			}
		})
	}
}

func TestMax(t *testing.T) {
	type Request[T int | float64 | string | bool] struct {
		Field1 T  `validate:"max=5"`
		Field2 *T `validate:"max=5"`
	}

	maxIntOk := 4
	maxIntNotOk := 6
	maxFloatOk := 4.999
	maxFloatNotOk := 5.001
	maxStringOk := "abcd"
	maxStringNotOk := "abcdef"
	invalidType := true

	testTable := []struct {
		name         string
		intInput     Request[int]
		floatInput   Request[float64]
		stringInput  Request[string]
		invalidInput Request[bool]
		output       int
	}{
		{
			name: "IntOk",
			intInput: Request[int]{
				Field1: 4,
				Field2: &maxIntOk,
			},
			output: 0,
		},
		{
			name: "IntNotOk",
			intInput: Request[int]{
				Field1: 6,
				Field2: &maxIntNotOk,
			},
			output: 2,
		},
		{
			name: "FloatOk",
			floatInput: Request[float64]{
				Field1: 4.0,
				Field2: &maxFloatOk,
			},
			output: 0,
		},
		{
			name: "FloatNotOk",
			floatInput: Request[float64]{
				Field1: 6.0,
				Field2: &maxFloatNotOk,
			},
			output: 2,
		},
		{
			name: "StringOk",
			stringInput: Request[string]{
				Field1: "abcd",
				Field2: &maxStringOk,
			},
			output: 0,
		},
		{
			name: "StringNotOk",
			stringInput: Request[string]{
				Field1: "abcdef",
				Field2: &maxStringNotOk,
			},
			output: 2,
		},
		{
			name: "InvalidType",
			invalidInput: Request[bool]{
				Field1: false,
				Field2: &invalidType,
			},
			output: 2,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			var errors []ValidationError
			var err error
			isInvalidTest := false

			if tt.intInput.Field1 != 0 {
				errors, err = Validate(tt.intInput)
			} else if tt.floatInput.Field1 != 0 {
				errors, err = Validate(tt.floatInput)
			} else if tt.stringInput.Field1 != "" {
				errors, err = Validate(tt.stringInput)
			} else {
				errors, err = Validate(tt.invalidInput)
				isInvalidTest = true
			}

			if err != nil {
				t.Fatal(err)
			}

			LogErrorsJson(t, errors)

			if len(errors) != tt.output {
				t.Errorf("\nExpected: %v.\nResult: %v.", 0, len(errors))
			}

			if len(errors) == 0 {
				return
			}

			invalidTypeMessage := validator.NOTSTRINGORNUMERIC_ERROR
			pattern := `^Must (have|be) less or equal than`
			re := regexp.MustCompile(pattern)
			for _, validationError := range errors {
				for _, error := range validationError.Errors {
					if isInvalidTest && error != invalidTypeMessage {
						t.Errorf("Expected: %v. Result: %v", invalidTypeMessage, error)
					}

					if !isInvalidTest && !re.MatchString(error) {
						t.Errorf("Expected to match: %v. Result: %v", pattern, error)
					}
				}
			}
		})
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
