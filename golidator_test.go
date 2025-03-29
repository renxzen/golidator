package golidator_test

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/renxzen/golidator"
	"github.com/renxzen/golidator/internal/validate"
)

func LogErrorsJson(t *testing.T, errors []golidator.ValidationError) {
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
			errors, err := golidator.Validate(tt.input)
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
			var errors []golidator.ValidationError
			var err error
			if tt.badInput.Field1 != 0 {
				errors, err = golidator.Validate(tt.badInput)
			} else {
				errors, err = golidator.Validate(tt.input)
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
				message = validate.ErrMsgNotStringType
			} else {
				message = validate.ErrMsgInvalidEmail
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
			var errors []golidator.ValidationError
			var err error
			if tt.badInput.Field1 != 0 {
				errors, err = golidator.Validate(tt.badInput)
			} else {
				errors, err = golidator.Validate(tt.input)
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
				message = validate.ErrMsgNotStringType
			} else {
				message = validate.ErrMsgInvalidURL
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
			var errors []golidator.ValidationError
			var err error
			isInvalidTest := false

			if tt.intInput.Field1 != 0 {
				errors, err = golidator.Validate(tt.intInput)
			} else if tt.floatInput.Field1 != 0 {
				errors, err = golidator.Validate(tt.floatInput)
			} else if tt.stringInput.Field1 != "" {
				errors, err = golidator.Validate(tt.stringInput)
			} else {
				errors, err = golidator.Validate(tt.invalidInput)
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

			invalidTypeMessage := validate.ErrMsgNotStrIntType
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
			var errors []golidator.ValidationError
			var err error
			isInvalidTest := false

			if tt.intInput.Field1 != 0 {
				errors, err = golidator.Validate(tt.intInput)
			} else if tt.floatInput.Field1 != 0 {
				errors, err = golidator.Validate(tt.floatInput)
			} else if tt.stringInput.Field1 != "" {
				errors, err = golidator.Validate(tt.stringInput)
			} else {
				errors, err = golidator.Validate(tt.invalidInput)
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

			invalidTypeMessage := validate.ErrMsgNotStrIntType
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

func TestNotempty(t *testing.T) {
	type Request struct {
		Field1 []int  `validate:"notempty"`
		Field2 []*int `validate:"notempty"`
		Field3 *[]int `validate:"notempty"`
	}

	type BadRequest struct {
		Field1 int      `validate:"notempty"`
		Field2 *float64 `validate:"notempty"`
	}

	ok := 1
	ok2 := []int{1}
	notOk := []int{}
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
				Field1: []int{1},
				Field2: []*int{&ok},
				Field3: &ok2,
			},
			output: 0,
		},
		{
			name: "NotOk",
			input: Request{
				Field1: []int{},
				Field2: []*int{},
				Field3: &notOk,
			},
			output: 3,
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
			var errors []golidator.ValidationError
			var err error

			if tt.input.Field1 != nil {
				errors, err = golidator.Validate(tt.input)
			} else {
				errors, err = golidator.Validate(tt.badInput)
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
			if tt.input.Field1 != nil {
				message = validate.ErrMsgEmptyArray
			} else {
				message = validate.ErrMsgNotArrayType
			}

			for _, validationError := range errors {
				for _, error := range validationError.Errors {
					if error != message {
						t.Errorf("Expected: %v. Result: %v", message, error)
					}
				}
			}
		})
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

	errors, err := golidator.Validate(input)
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

	errors, err := golidator.Validate(input)
	if err != nil {
		t.Fatal(err)
	}

	LogErrorsJson(t, errors)

	if len(errors) != 5 {
		t.Errorf("\nExpected: %v.\nResult: %v.", 0, len(errors))
	}
}

func TestFieldName(t *testing.T) {
	type Request struct {
		Field1 string  `json:"field_1" validate:"notblank"`
		Field2 *string `json:"field_2" validate:"notblank"`
		Field3 string  `               validate:"notblank"`
		Field4 *string `               validate:"notblank"`
	}

	blankField := ""
	testTable := []struct {
		name   string
		input  Request
		output []string
	}{
		{
			name: "Ok",
			input: Request{
				Field1: "",
				Field2: &blankField,
				Field3: "",
				Field4: &blankField,
			},
			output: []string{"field_1", "field_2", "Field3", "Field4"},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			errors, err := golidator.Validate(tt.input)
			if err != nil {
				t.Fatal(err)
			}
			LogErrorsJson(t, errors)

			count := 0
			for _, fieldName := range tt.output {
				for _, err := range errors {
					if err.Field == fieldName {
						count++
					}
				}
			}
			if len(tt.output) != count {
				t.Errorf("\nExpected: %v.\nResult: %v.", len(tt.output), count)
			}
		})
	}
}

func TestNumeric(t *testing.T) {
	type Request struct {
		Field1 string  `validate:"numeric"`
		Field2 *string `validate:"numeric"`
	}

	type BadRequest struct {
		Field1 int      `validate:"numeric"`
		Field2 *float64 `validate:"numeric"`
	}

	numericOK := "12345"
	numericNotOK := "12345ABCDE"
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
				Field1: numericOK,
				Field2: &numericOK,
			},
			output: 0,
		},
		{
			name: "NotOk",
			input: Request{
				Field1: numericNotOK,
				Field2: &numericNotOK,
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
			var errors []golidator.ValidationError
			var err error
			if tt.badInput.Field1 != 0 {
				errors, err = golidator.Validate(tt.badInput)
			} else {
				errors, err = golidator.Validate(tt.input)
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
				message = validate.ErrMsgNotStringType
			} else {
				message = validate.ErrMsgNotNumeric
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

func TestLen(t *testing.T) {
	type Request[T string | []string] struct {
		Field1 T  `validate:"len=5"`
		Field2 *T `validate:"len=5"`
	}

	type BadRequest struct {
		Field1 int      `validate:"len=5"`
		Field2 *float64 `validate:"len=5"`
	}

	stringOK := "abcde"
	stringNotOK := "abcd"
	sliceOK := []string{"a", "b", "c", "d", "e"}
	sliceNotOK := []string{"a", "b", "c", "d"}
	testTable := []struct {
		name        string
		stringInput Request[string]
		sliceInput  Request[[]string]
		badInput    BadRequest
		output      int
	}{
		{
			name: "StringOK",
			stringInput: Request[string]{
				Field1: stringOK,
				Field2: &stringOK,
			},
			output: 0,
		},
		{
			name: "StringNotOk",
			stringInput: Request[string]{
				Field1: stringNotOK,
				Field2: &stringNotOK,
			},
			output: 2,
		},
		{
			name: "SliceOK",
			sliceInput: Request[[]string]{
				Field1: sliceOK,
				Field2: &sliceOK,
			},
			output: 0,
		},
		{
			name: "SliceNotOk",
			sliceInput: Request[[]string]{
				Field1: sliceNotOK,
				Field2: &sliceNotOK,
			},
			output: 2,
		},
		{
			name: "InvalidType",
			badInput: BadRequest{
				Field1: 69,
				Field2: new(float64),
			},
			output: 2,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			var errors []golidator.ValidationError
			var err error

			if tt.stringInput.Field1 != "" {
				errors, err = golidator.Validate(tt.stringInput)
			} else if tt.sliceInput.Field1 != nil {
				errors, err = golidator.Validate(tt.sliceInput)
			} else if tt.badInput.Field1 != 0 {
				errors, err = golidator.Validate(tt.badInput)
			}
			if err != nil {
				t.Fatal(err)
			}

			LogErrorsJson(t, errors)

			if len(errors) != tt.output {
				t.Errorf("\nExpected: %v.\nResult: %v.", tt.output, len(errors))
			}

			if len(errors) == 0 {
				return
			}

			message := fmt.Sprintf(validate.ErrMsgInvalidLength, 5)
			if tt.sliceInput.Field1 != nil {
				message = fmt.Sprintf(validate.ErrMsgInvalidLengthSlice, 5)
			} else if tt.badInput.Field1 != 0 {
				message = validate.ErrMsgNotStrSliceType
			}

			for _, validationError := range errors {
				for _, error := range validationError.Errors {
					if error != message {
						t.Errorf("Expected: %v. Result: %v", message, error)
					}
				}
			}
		})
	}
}
