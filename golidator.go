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
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	var results []ValidationError

	for i := range value.NumField() {
		// get field value
		fieldValue := value.Field(i)
		fieldValueType := fieldValue.Type()

		// get field type
		typeField := value.Type().Field(i)
		typeFieldTypeName := typeField.Type.Name()

		// check if field is a pointer
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			fieldValue = fieldValue.Elem()
			fieldValueType = fieldValue.Type()
			typeFieldTypeName = fieldValueType.Name()
		}

		// get json name
		jsonTag := strings.Split(typeField.Tag.Get(JsonTag), ",")

		var typeFieldName string
		if len(jsonTag) > 0 && jsonTag[0] != "" {
			typeFieldName = jsonTag[0]
		} else {
			typeFieldName = typeField.Name
		}

		valError := ValidationError{Field: typeFieldName}

		validateTag := typeField.Tag.Get(TagName)
		validators := strings.SplitSeq(validateTag, ",")
		for validator := range validators {
			if validator == "" {
				continue
			}

			args := strings.Split(validator, "=")
			if len(args) == 0 {
				continue
			}

			var fieldLength int
			if len(args) > 1 {
				limit, err := strconv.Atoi(args[1])
				if err != nil {
					return nil, fmt.Errorf("invalid parameter %q used in %q validation", args[1], args[0])
				}
				fieldLength = limit
			}

			var err error
			switch args[0] {
			case "notblank":
				err = validate.NotBlank(fieldValue, typeFieldTypeName)
			case "email":
				err = validate.Email(fieldValue, typeFieldTypeName)
			case "numeric":
				err = validate.Numeric(fieldValue, typeFieldTypeName)
			case "url":
				err = validate.URL(fieldValue, typeFieldTypeName)
			case "required":
				err = validate.Required(fieldValue)
			case "notempty":
				err = validate.NotEmpty(fieldValue, fieldValueType)
			case "min":
				err = validate.Min(fieldValue, typeFieldTypeName, fieldLength)
			case "max":
				err = validate.Max(fieldValue, typeFieldTypeName, fieldLength)
			case "len":
				err = validate.Len(fieldValue, fieldValueType, typeFieldTypeName, fieldLength)
			case "isarray":
				err = validate.IsArray(fieldValue, fieldValueType)

				if err == nil && fieldValue.Kind() == reflect.Slice {
					for j := range fieldValue.Len() {
						result, subErr := Validate(fieldValue.Index(j).Interface())
						if subErr != nil {
							valError.Errors = append(valError.Errors, subErr.Error())
							continue
						}

						for _, arr := range result {
							subFieldName := fmt.Sprintf("%v[%v]: %v", typeFieldName, j, arr.Field)
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
				valError.Errors = append(valError.Errors, err.Error())
			}
		}

		if len(valError.Errors) > 0 {
			results = append(results, valError)
		}
	}
	return results, nil
}
