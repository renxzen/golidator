package golidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/renxzen/golidator/internal/validate"
)

const (
	TagName = "validate"
	JsonTag = "json"
)

type ValidationError struct {
	Field  string   `json:"field"`
	Errors []string `json:"errors"`
}

func Validate(model any) ([]ValidationError, error) {
	value := reflect.ValueOf(model)
	kind := value.Kind()

	if kind == reflect.Ptr {
		if value.IsNil() {
			return nil, nil
		}

		value = value.Elem()
	}

	if kind != reflect.Struct {
		return nil, fmt.Errorf("model must be a struct, got %s", kind)
	}

	numField := value.NumField()
	results := make([]ValidationError, 0, numField) // estimate capacity

	for i := range numField {
		// get field value
		fieldValue := value.Field(i)
		fieldValueType := fieldValue.Type()
		fieldValueKind := fieldValue.Kind()

		// get field type
		fieldType := value.Type().Field(i)
		fieldTypeName := fieldType.Type.Name()

		// check if field is a pointer
		if fieldValueKind == reflect.Ptr && !fieldValue.IsNil() {
			fieldValue = fieldValue.Elem()
			fieldValueType = fieldValue.Type()
			fieldValueKind = fieldValue.Kind()
			fieldTypeName = fieldValueType.Name()
		}

		// get json name
		jsonTag := strings.Split(fieldType.Tag.Get(JsonTag), ",")

		var fieldName string
		if len(jsonTag) > 0 && jsonTag[0] != "" {
			fieldName = jsonTag[0]
		} else {
			fieldName = fieldType.Name
		}

		validateTag := fieldType.Tag.Get(TagName)
		if validateTag == "" {
			// skip validation if no validation tag is present
			continue
		}

		validators := strings.SplitSeq(validateTag, ",")

		var valErrors []string
		for validator := range validators {
			if validator == "" {
				continue
			}

			args := strings.Split(validator, "=")
			if len(args) == 0 {
				continue
			}

			validatorName := args[0]

			var fieldLength int
			if len(args) > 1 {
				var err error
				fieldLength, err = strconv.Atoi(args[1])
				if fieldLength < 0 || err != nil {
					return nil, fmt.Errorf("invalid parameter %q used in %q validation", args[1], validatorName)
				}
			}

			var err error
			switch validatorName {
			case "notblank":
				err = validate.NotBlank(fieldValue, fieldValueKind, fieldTypeName)
			case "email":
				err = validate.Email(fieldValue, fieldValueKind, fieldTypeName)
			case "numeric":
				err = validate.Numeric(fieldValue, fieldValueKind, fieldTypeName)
			case "url":
				err = validate.URL(fieldValue, fieldValueKind, fieldTypeName)
			case "required":
				err = validate.Required(fieldValue, fieldValueKind)
			case "notempty":
				err = validate.NotEmpty(fieldValue, fieldValueType, fieldValueKind)
			case "min":
				err = validate.Min(fieldValue, fieldValueKind, fieldTypeName, fieldLength)
			case "max":
				err = validate.Max(fieldValue, fieldValueKind, fieldTypeName, fieldLength)
			case "len":
				err = validate.Len(fieldValue, fieldValueType, fieldValueKind, fieldTypeName, fieldLength)
			case "isarray":
				err = validate.IsArray(fieldValue, fieldValueType, fieldValueKind)

				if err == nil && fieldValueKind == reflect.Slice {
					for j := range fieldValue.Len() {
						result, subErr := Validate(fieldValue.Index(j).Interface())
						if subErr != nil {
							valErrors = append(valErrors, subErr.Error())
							continue
						}

						for _, arr := range result {
							subFieldName := fmt.Sprintf("%v[%v]: %v", fieldName, j, arr.Field)
							results = append(results, ValidationError{
								Field:  subFieldName,
								Errors: arr.Errors,
							})
						}
					}
				}
			default:
				return nil, fmt.Errorf("unknown validator: %s", args[0])
			}

			if err != nil {
				valErrors = append(valErrors, err.Error())
			}
		}

		if len(valErrors) > 0 {
			results = append(results, ValidationError{
				Field:  fieldName,
				Errors: valErrors,
			})
		}
	}
	return results, nil
}
