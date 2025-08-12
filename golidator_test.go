package golidator_test

import (
	"encoding/json"
	"regexp"
	"slices"
	"testing"

	"github.com/renxzen/golidator"
	"github.com/renxzen/golidator/internal/validators"
)

func TestStringValidators(t *testing.T) {
	type StringValidationStruct struct {
		NotBlankField string  `json:"notblank_field" validate:"notblank"`
		NotBlankPtr   *string `json:"notblank_ptr"   validate:"notblank"`
		EmailField    string  `json:"email_field"    validate:"email"`
		EmailPtr      *string `json:"email_ptr"      validate:"email"`
		URLField      string  `json:"url_field"      validate:"url"`
		URLPtr        *string `json:"url_ptr"        validate:"url"`
		NumericField  string  `json:"numeric_field"  validate:"numeric"`
		NumericPtr    *string `json:"numeric_ptr"    validate:"numeric"`
	}

	type InvalidTypeStruct struct {
		NotBlankField int `validate:"notblank"`
		EmailField    int `validate:"email"`
		URLField      int `validate:"url"`
		NumericField  int `validate:"numeric"`
	}

	tests := []validationTestCase{
		{
			name: "all_valid",
			input: StringValidationStruct{
				NotBlankField: "not blank",
				NotBlankPtr:   ptr("not blank ptr"),
				EmailField:    "test@example.com",
				EmailPtr:      ptr("ptr@example.com"),
				URLField:      "https://example.com",
				URLPtr:        ptr("http://example.org"),
				NumericField:  "12345",
				NumericPtr:    ptr("67890"),
			},
			expectedErrors: 0,
		},
		{
			name: "notblank_failures",
			input: StringValidationStruct{
				NotBlankField: "",
				NotBlankPtr:   ptr(""),
				EmailField:    "valid@example.com",
				EmailPtr:      ptr("valid@example.com"),
				URLField:      "https://example.com",
				URLPtr:        ptr("https://example.com"),
				NumericField:  "123",
				NumericPtr:    ptr("456"),
			},
			expectedErrors: 2,
			expectedFields: []string{"notblank_field", "notblank_ptr"},
		},
		{
			name: "email_failures",
			input: StringValidationStruct{
				NotBlankField: "valid",
				NotBlankPtr:   ptr("valid"),
				EmailField:    "invalid-email",
				EmailPtr:      ptr("also@invalid"),
				URLField:      "https://example.com",
				URLPtr:        ptr("https://example.com"),
				NumericField:  "123",
				NumericPtr:    ptr("456"),
			},
			expectedErrors: 2,
			expectedFields: []string{"email_field", "email_ptr"},
			errorMessages: []string{
				validators.MessageInvalidEmail,
				validators.MessageInvalidEmail,
			},
		},
		{
			name: "url_failures",
			input: StringValidationStruct{
				NotBlankField: "valid",
				NotBlankPtr:   ptr("valid"),
				EmailField:    "valid@example.com",
				EmailPtr:      ptr("valid@example.com"),
				URLField:      "not-a-url",
				URLPtr:        ptr("also.not.url"),
				NumericField:  "123",
				NumericPtr:    ptr("456"),
			},
			expectedErrors: 2,
			expectedFields: []string{"url_field", "url_ptr"},
			errorMessages: []string{
				validators.MessageInvalidURL,
				validators.MessageInvalidURL,
			},
		},
		{
			name: "numeric_failures",
			input: StringValidationStruct{
				NotBlankField: "valid",
				NotBlankPtr:   ptr("valid"),
				EmailField:    "valid@example.com",
				EmailPtr:      ptr("valid@example.com"),
				URLField:      "https://example.com",
				URLPtr:        ptr("https://example.com"),
				NumericField:  "abc123",
				NumericPtr:    ptr("123abc"),
			},
			expectedErrors: 2,
			expectedFields: []string{"numeric_field", "numeric_ptr"},
			errorMessages: []string{
				validators.MessageNotNumeric,
				validators.MessageNotNumeric,
			},
		},
		{
			name: "invalid_types",
			input: InvalidTypeStruct{
				NotBlankField: 123,
				EmailField:    456,
				URLField:      789,
				NumericField:  101112,
			},
			expectedErrors: 4,
			errorMessages: []string{
				validators.MessageNotStringType,
				validators.MessageNotStringType,
				validators.MessageNotStringType,
				validators.MessageNotStringType,
			},
		},
		{
			name: "edge_cases",
			input: StringValidationStruct{
				NotBlankField: "   ",
				NotBlankPtr:   ptr("\t\n"),
				EmailField:    "user+tag@example.co.uk",
				EmailPtr:      ptr("test.email@sub.domain.com"),
				URLField:      "https://example.com:8080/path?query=value#fragment",
				URLPtr:        ptr("ftp://files.example.com/file.txt"),
				NumericField:  "0",
				NumericPtr:    ptr("000123"),
			},
			expectedErrors: 0,
		},
	}

	runValidationTests(t, tests)
}

func TestNumericRangeValidators(t *testing.T) {
	type NumericValidationStruct struct {
		MinIntField    int      `json:"min_int"        validate:"min=5"`
		MinIntPtr      *int     `json:"min_int_ptr"    validate:"min=5"`
		MaxIntField    int      `json:"max_int"        validate:"max=10"`
		MaxIntPtr      *int     `json:"max_int_ptr"    validate:"max=10"`
		MinFloatField  float64  `json:"min_float"      validate:"min=5"`
		MinFloatPtr    *float64 `json:"min_float_ptr"  validate:"min=5"`
		MaxFloatField  float64  `json:"max_float"      validate:"max=10"`
		MaxFloatPtr    *float64 `json:"max_float_ptr"  validate:"max=10"`
		MinStringField string   `json:"min_string"     validate:"min=5"`
		MinStringPtr   *string  `json:"min_string_ptr" validate:"min=5"`
		MaxStringField string   `json:"max_string"     validate:"max=10"`
		MaxStringPtr   *string  `json:"max_string_ptr" validate:"max=10"`
	}

	type InvalidTypeStruct struct {
		MinField bool `validate:"min=5"`
		MaxField bool `validate:"max=10"`
	}

	tests := []validationTestCase{
		{
			name: "all_valid",
			input: NumericValidationStruct{
				MinIntField:    7,
				MinIntPtr:      ptr(6),
				MaxIntField:    8,
				MaxIntPtr:      ptr(9),
				MinFloatField:  6.5,
				MinFloatPtr:    ptr(7.2),
				MaxFloatField:  8.3,
				MaxFloatPtr:    ptr(9.1),
				MinStringField: "abcdef",
				MinStringPtr:   ptr("abcdefg"),
				MaxStringField: "abcdefghi",
				MaxStringPtr:   ptr("abcd"),
			},
			expectedErrors: 0,
		},
		{
			name: "min_failures",
			input: NumericValidationStruct{
				MinIntField:    3,
				MinIntPtr:      ptr(4),
				MaxIntField:    8,
				MaxIntPtr:      ptr(9),
				MinFloatField:  3.5,
				MinFloatPtr:    ptr(4.9),
				MaxFloatField:  8.3,
				MaxFloatPtr:    ptr(9.1),
				MinStringField: "abc",
				MinStringPtr:   ptr("abcd"),
				MaxStringField: "abcdefghi",
				MaxStringPtr:   ptr("abcd"),
			},
			expectedErrors: 6,
			expectedFields: []string{
				"min_int",
				"min_int_ptr",
				"min_float",
				"min_float_ptr",
				"min_string",
				"min_string_ptr",
			},
			errorPattern: `^must (have|be) more or equal than`,
		},
		{
			name: "max_failures",
			input: NumericValidationStruct{
				MinIntField:    7,
				MinIntPtr:      ptr(6),
				MaxIntField:    15,
				MaxIntPtr:      ptr(12),
				MinFloatField:  6.5,
				MinFloatPtr:    ptr(7.2),
				MaxFloatField:  11.3,
				MaxFloatPtr:    ptr(12.1),
				MinStringField: "abcdef",
				MinStringPtr:   ptr("abcdefg"),
				MaxStringField: "abcdefghijk",
				MaxStringPtr:   ptr("abcdefghijkl"),
			},
			expectedErrors: 6,
			expectedFields: []string{
				"max_int",
				"max_int_ptr",
				"max_float",
				"max_float_ptr",
				"max_string",
				"max_string_ptr",
			},
			errorPattern: `^must (have|be) less or equal than`,
		},
		{
			name: "invalid_types",
			input: InvalidTypeStruct{
				MinField: true,
				MaxField: false,
			},
			expectedErrors: 2,
			errorMessages: []string{
				validators.MessageNotStrIntType,
				validators.MessageNotStrIntType,
			},
		},
		{
			name: "edge_cases",
			input: NumericValidationStruct{
				MinIntField:    5,
				MinIntPtr:      ptr(5),
				MaxIntField:    10,
				MaxIntPtr:      ptr(10),
				MinFloatField:  5.0,
				MinFloatPtr:    ptr(5.0),
				MaxFloatField:  10.0,
				MaxFloatPtr:    ptr(10.0),
				MinStringField: "abcde",
				MinStringPtr:   ptr("abcde"),
				MaxStringField: "abcdefghij",
				MaxStringPtr:   ptr("abcdefghij"),
			},
			expectedErrors: 0,
		},
	}

	runValidationTests(t, tests)
}

func TestArraySliceValidators(t *testing.T) {
	type SubStruct struct {
		Name  string `json:"name"  validate:"notblank"`
		Email string `json:"email" validate:"email"`
	}

	type ArrayValidationStruct struct {
		NotEmptyIntSlice    []int       `json:"notempty_int"       validate:"notempty"`
		NotEmptyIntPtrSlice []*int      `json:"notempty_int_ptr"   validate:"notempty"`
		NotEmptySlicePtr    *[]int      `json:"notempty_slice_ptr" validate:"notempty"`
		IsArrayField        []SubStruct `json:"isarray_field"      validate:"isarray"`
		LenStringField      string      `json:"len_string"         validate:"len=5"`
		LenStringPtr        *string     `json:"len_string_ptr"     validate:"len=5"`
		LenSliceField       []string    `json:"len_slice"          validate:"len=3"`
		LenSlicePtr         *[]string   `json:"len_slice_ptr"      validate:"len=3"`
	}

	type InvalidTypeStruct struct {
		NotEmptyField int `validate:"notempty"`
		IsArrayField  int `validate:"isarray"`
		LenField      int `validate:"len=5"`
	}

	validInt := 1
	validSlice := []int{1, 2}
	validString := "abcde"
	validStringSlice := []string{"a", "b", "c"}

	tests := []validationTestCase{
		{
			name: "all_valid",
			input: ArrayValidationStruct{
				NotEmptyIntSlice:    []int{1, 2, 3},
				NotEmptyIntPtrSlice: []*int{&validInt},
				NotEmptySlicePtr:    &validSlice,
				IsArrayField: []SubStruct{
					{Name: "John", Email: "john@example.com"},
					{Name: "Jane", Email: "jane@example.com"},
				},
				LenStringField: "abcde",
				LenStringPtr:   &validString,
				LenSliceField:  []string{"a", "b", "c"},
				LenSlicePtr:    &validStringSlice,
			},
			expectedErrors: 0,
		},
		{
			name: "notempty_failures",
			input: ArrayValidationStruct{
				NotEmptyIntSlice:    []int{},
				NotEmptyIntPtrSlice: []*int{},
				NotEmptySlicePtr:    &[]int{},
				IsArrayField: []SubStruct{
					{Name: "John", Email: "john@example.com"},
				},
				LenStringField: "abcde",
				LenStringPtr:   &validString,
				LenSliceField:  []string{"a", "b", "c"},
				LenSlicePtr:    &validStringSlice,
			},
			expectedErrors: 3,
			expectedFields: []string{"notempty_int", "notempty_int_ptr", "notempty_slice_ptr"},
			errorMessages: []string{
				validators.MessageEmptyArray,
				validators.MessageEmptyArray,
				validators.MessageEmptyArray,
			},
		},
		{
			name: "len_failures",
			input: ArrayValidationStruct{
				NotEmptyIntSlice:    []int{1},
				NotEmptyIntPtrSlice: []*int{&validInt},
				NotEmptySlicePtr:    &validSlice,
				IsArrayField: []SubStruct{
					{Name: "John", Email: "john@example.com"},
				},
				LenStringField: "abc",
				LenStringPtr:   ptr("abcdefg"),
				LenSliceField:  []string{"a", "b"},
				LenSlicePtr:    &[]string{"a", "b", "c", "d"},
			},
			expectedErrors: 4,
			expectedFields: []string{"len_string", "len_string_ptr", "len_slice", "len_slice_ptr"},
		},
		{
			name: "isarray_with_nested_errors",
			input: ArrayValidationStruct{
				NotEmptyIntSlice:    []int{1},
				NotEmptyIntPtrSlice: []*int{&validInt},
				NotEmptySlicePtr:    &validSlice,
				IsArrayField: []SubStruct{
					{Name: "", Email: "invalid-email"},
					{Name: "Jane", Email: "jane@example.com"},
					{Name: "Bob", Email: "@invalid.com"},
				},
				LenStringField: "abcde",
				LenStringPtr:   &validString,
				LenSliceField:  []string{"a", "b", "c"},
				LenSlicePtr:    &validStringSlice,
			},
			expectedErrors: 3,
		},
		{
			name: "invalid_types",
			input: InvalidTypeStruct{
				NotEmptyField: 123,
				IsArrayField:  456,
				LenField:      789,
			},
			expectedErrors: 3,
			errorMessages: []string{
				validators.MessageNotArrayType,
				validators.MessageNotArrayType,
				validators.MessageNotStrSliceType,
			},
		},
	}

	runValidationTests(t, tests)
}

func TestComprehensiveEdgeCases(t *testing.T) {
	type NestedStruct struct {
		ID    int    `json:"id"    validate:"min=1"`
		Name  string `json:"name"  validate:"notblank,min=2,max=50"`
		Email string `json:"email" validate:"email"`
	}

	type ComplexStruct struct {
		RequiredString *string        `json:"required_string" validate:"required,notblank"`
		RequiredInt    *int           `json:"required_int"    validate:"required,min=0"`
		MultiValidated string         `json:"multi_validated" validate:"notblank,email,min=5,max=100"`
		Nested         NestedStruct   `json:"nested"`
		NestedPtr      *NestedStruct  `json:"nested_ptr"`
		NestedSlice    []NestedStruct `json:"nested_slice"    validate:"isarray,notempty"`
		StringSlice    []string       `json:"string_slice"    validate:"len=3"`
		IntSlice       []int          `json:"int_slice"       validate:"notempty"`
		OptionalURL    *string        `json:"optional_url"    validate:"url"`
		NumericStr     string         `json:"numeric_str"     validate:"numeric,len=5"`
	}

	tests := []validationTestCase{
		{
			name: "complex_all_valid",
			input: ComplexStruct{
				RequiredString: ptr("required_value"),
				RequiredInt:    ptr(10),
				MultiValidated: "test@example.com",
				Nested: NestedStruct{
					ID:    1,
					Name:  "John Doe",
					Email: "john@example.com",
				},
				NestedPtr: &NestedStruct{
					ID:    2,
					Name:  "Jane Smith",
					Email: "jane@example.com",
				},
				NestedSlice: []NestedStruct{
					{ID: 3, Name: "Bob", Email: "bob@example.com"},
					{ID: 4, Name: "Alice", Email: "alice@example.com"},
				},
				StringSlice: []string{"one", "two", "three"},
				IntSlice:    []int{1, 2, 3, 4, 5},
				OptionalURL: ptr("https://example.com"),
				NumericStr:  "12345",
			},
			expectedErrors: 0,
		},
		{
			name: "required_field_failures",
			input: ComplexStruct{
				RequiredString: nil,
				RequiredInt:    nil,
				MultiValidated: "test@example.com",
				Nested: NestedStruct{
					ID:    1,
					Name:  "John",
					Email: "john@example.com",
				},
				NestedSlice: []NestedStruct{
					{ID: 1, Name: "Test", Email: "test@example.com"},
				},
				StringSlice: []string{"a", "b", "c"},
				IntSlice:    []int{1},
				OptionalURL: ptr("invalid-url"),
				NumericStr:  "12345",
			},
			expectedErrors: 3,
			expectedFields: []string{"required_string", "required_int", "optional_url"},
		},
		{
			name: "multi_validation_failures",
			input: ComplexStruct{
				RequiredString: ptr("valid"),
				RequiredInt:    ptr(5),
				MultiValidated: "abc",
				Nested: NestedStruct{
					ID:    1,
					Name:  "John",
					Email: "john@example.com",
				},
				NestedSlice: []NestedStruct{
					{ID: 1, Name: "Test", Email: "test@example.com"},
				},
				StringSlice: []string{"a", "b", "c"},
				IntSlice:    []int{1},
				OptionalURL: ptr("invalid-url"),
				NumericStr:  "12345",
			},
			expectedErrors: 2,
			expectedFields: []string{"multi_validated", "optional_url"},
		},
		{
			name: "nested_validation_failures",
			input: ComplexStruct{
				RequiredString: ptr("valid"),
				RequiredInt:    ptr(5),
				MultiValidated: "test@example.com",
				Nested: NestedStruct{
					ID:    1,
					Name:  "Valid",
					Email: "valid@example.com",
				},
				NestedSlice: []NestedStruct{
					{ID: 1, Name: "Valid", Email: "valid@example.com"},
					{ID: 0, Name: "X", Email: "bad-email"},
				},
				StringSlice: []string{"a", "b", "c"},
				IntSlice:    []int{1},
				OptionalURL: ptr("invalid-url"),
				NumericStr:  "12345",
			},
			expectedErrors: 4,
		},
		{
			name: "array_slice_failures",
			input: ComplexStruct{
				RequiredString: ptr("valid"),
				RequiredInt:    ptr(5),
				MultiValidated: "test@example.com",
				Nested: NestedStruct{
					ID:    1,
					Name:  "Valid",
					Email: "valid@example.com",
				},
				NestedSlice: []NestedStruct{},
				StringSlice: []string{"a", "b"},
				IntSlice:    []int{},
				OptionalURL: ptr("invalid-url"),
				NumericStr:  "abc12",
			},
			expectedErrors: 5,
			expectedFields: []string{
				"nested_slice",
				"string_slice",
				"int_slice",
				"optional_url",
				"numeric_str",
			},
		},
		{
			name: "edge_cases_nil_empty",
			input: ComplexStruct{
				RequiredString: ptr(""),
				RequiredInt:    ptr(-1),
				MultiValidated: "",
				Nested: NestedStruct{
					ID:    1,
					Name:  "Valid",
					Email: "valid@example.com",
				},
				NestedPtr:   nil,
				NestedSlice: []NestedStruct{{ID: 1, Name: "Test", Email: "test@example.com"}},
				StringSlice: []string{"a", "b", "c"},
				IntSlice:    []int{1},
				OptionalURL: ptr("not-a-url"),
				NumericStr:  "1234",
			},
			expectedErrors: 5,
			expectedFields: []string{
				"required_string",
				"required_int",
				"multi_validated",
				"optional_url",
				"numeric_str",
			},
		},
	}

	runValidationTests(t, tests)
}

func TestCustomValidators(t *testing.T) {
	type ExpectedStringValueStruct struct {
		CorrectValue      string  `validate:"expected-value"`
		CorrectPtrValue   *string `validate:"expected-value"`
		IncorrectValue    string  `validate:"expected-value"`
		IncorrectPtrValue *string `validate:"expected-value"`
		InvalidType       int     `validate:"expected-value"`
		InvalidPtrType    *int    `validate:"expected-value"`
	}

	type ExpectedStringArgValueStruct struct {
		CorrectValue      string  `validate:"expected-str=val"`
		CorrectPtrValue   *string `validate:"expected-str=val"`
		IncorrectValue    string  `validate:"expected-str=val"`
		IncorrectPtrValue *string `validate:"expected-str=val"`
		InvalidType       int     `validate:"expected-str=val"`
		InvalidPtrType    *int    `validate:"expected-str=val"`
	}

	type ExpectedIntArgValueStruct struct {
		CorrectInt       int     `validate:"expected-int=10"`
		CorrectPtrInt    *int    `validate:"expected-int=10"`
		IncorrectInt     int     `validate:"expected-int=10"`
		IncorrectPtrInt  *int    `validate:"expected-int=10"`
		IncorrectType    string  `validate:"expected-int=10"`
		IncorrectPtrType *string `validate:"expected-int=10"`
	}

	golidator.AddValidator("expected-value", func(field golidator.FieldInfo) string {
		if !field.IsString() {
			return "not a string"
		}
		if field.String() != "expected-value" {
			return "incorrect value"
		}
		return ""
	})

	golidator.AddValidator("expected-str", func(field golidator.FieldInfo) string {
		if !field.IsString() {
			return "not a string"
		}
		value := field.GetArgumentStr("expected-str")
		if value == "" {
			return ""
		}
		if field.String() != value {
			return "incorrect value"
		}
		return ""
	})

	golidator.AddValidator("expected-int", func(field golidator.FieldInfo) string {
		if !field.IsInt() {
			return "not an int"
		}
		value, exists := field.GetArgumentInt("expected-int")
		if !exists {
			return ""
		}
		if field.Int() != int64(value) {
			return "incorrect value"
		}
		return ""
	})

	tests := []validationTestCase{
		{
			name: "custom_all_valid",
			input: ExpectedStringValueStruct{
				CorrectValue:      "expected-value",
				CorrectPtrValue:   ptr("expected-value"),
				IncorrectValue:    "invalid-value",
				IncorrectPtrValue: ptr("invalid-value"),
				InvalidType:       0,
				InvalidPtrType:    ptr(0),
			},
			expectedErrors: 4,
			errorMessages: []string{
				"incorrect value",
				"not a string",
			},
		},
		{
			name: "custom_arg_valid",
			input: ExpectedStringArgValueStruct{
				CorrectValue:      "val",
				CorrectPtrValue:   ptr("val"),
				IncorrectValue:    "invalid-value",
				IncorrectPtrValue: ptr("invalid-value"),
				InvalidType:       0,
				InvalidPtrType:    ptr(0),
			},
			expectedErrors: 4,
			errorMessages: []string{
				"incorrect value",
				"not a string",
			},
		},
		{
			name: "custom_arg_valid",
			input: ExpectedIntArgValueStruct{
				CorrectInt:       10,
				CorrectPtrInt:    ptr(10),
				IncorrectInt:     11,
				IncorrectPtrInt:  ptr(11),
				IncorrectType:    "invalid-value",
				IncorrectPtrType: ptr("invalid-value"),
			},
			expectedErrors: 4,
			errorMessages: []string{
				"incorrect value",
				"not an int",
			},
		},
	}

	runValidationTests(t, tests)
}

func TestFieldNameMapping(t *testing.T) {
	tests := []validationTestCase{
		{
			name: "json_tag_mapping",
			input: struct {
				Field1 string  `json:"field_1" validate:"notblank"`
				Field2 *string `json:"field_2" validate:"notblank"`
				Field3 string  `validate:"notblank"`
				Field4 *string `validate:"notblank"`
			}{
				Field1: "",
				Field2: ptr(""),
				Field3: "",
				Field4: ptr(""),
			},
			expectedErrors: 4,
			expectedFields: []string{"field_1", "field_2", "Field3", "Field4"},
		},
	}

	runValidationTests(t, tests)
}

func TestCachingBehavior(t *testing.T) {
	type TestStruct struct {
		Name  string `validate:"notblank"`
		Email string `validate:"email"`
	}

	validStruct := TestStruct{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	golidator.SetCaching(true)
	errors1, err1 := golidator.Validate(validStruct)
	if err1 != nil {
		t.Fatal(err1)
	}
	if len(errors1) != 0 {
		t.Errorf("Expected no errors with caching enabled, got %d", len(errors1))
	}

	golidator.SetCaching(false)
	errors2, err2 := golidator.Validate(validStruct)
	if err2 != nil {
		t.Fatal(err2)
	}
	if len(errors2) != 0 {
		t.Errorf("Expected no errors with caching disabled, got %d", len(errors2))
	}

	golidator.SetCaching(true)
}

func logErrorsJSON(t *testing.T, errors []golidator.ValidationError) {
	b, err := json.MarshalIndent(errors, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(b))
}

func ptr[T any](s T) *T {
	return &s
}

type validationTestCase struct {
	name           string
	input          any
	expectedErrors int
	expectedFields []string
	errorMessages  []string
	errorPattern   string
}

func runValidationTests(t *testing.T, tests []validationTestCase) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors, err := golidator.Validate(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			logErrorsJSON(t, errors)

			if len(errors) != tt.expectedErrors {
				t.Errorf("Expected %d errors, got %d", tt.expectedErrors, len(errors))
			}

			if tt.expectedErrors == 0 {
				return
			}

			if len(tt.expectedFields) > 0 {
				errorFields := make(map[string]bool)
				for _, e := range errors {
					errorFields[e.Field] = true
				}

				for _, field := range tt.expectedFields {
					if !errorFields[field] {
						t.Errorf("Expected error for field %q but none was reported", field)
					}
				}
			}

			if len(tt.errorMessages) > 0 {
				for _, validationError := range errors {
					for _, errorMsg := range validationError.Errors {
						if !slices.Contains(tt.errorMessages, errorMsg) {
							t.Errorf("Expected error message %q but none was reported", errorMsg)
						}
					}
				}
			} else if tt.errorPattern != "" {
				re := regexp.MustCompile(tt.errorPattern)
				for _, validationError := range errors {
					for _, errorMsg := range validationError.Errors {
						if !re.MatchString(errorMsg) {
							t.Errorf("Error message %q does not match pattern %q", errorMsg, tt.errorPattern)
						}
					}
				}
			}
		})
	}
}
