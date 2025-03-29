package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

const (
	TagName = "validate"
	JsonTag = "json"
)

type Validator struct {
	value  reflect.Value
	errors map[string][]string
}

func NewValidate(model any) *Validator {
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	return &Validator{
		value:  value,
		errors: make(map[string][]string),
	}
}

func (v *Validator) GetErrors() (map[string][]string, error) {
	for i := 0; i < v.value.NumField(); i++ {
		// get field value
		fieldValue := v.value.Field(i)
		fieldValueType := fieldValue.Type()

		// get field type
		typeField := v.value.Type().Field(i)
		typeFieldTypeName := typeField.Type.Name()

		// check if field is a pointer
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			fieldValue = fieldValue.Elem()
			fieldValueType = fieldValue.Type()
			typeFieldTypeName = fieldValueType.Name()
		}

		// get json name
		var typeFieldName string
		jsonTag := strings.Split(typeField.Tag.Get(JsonTag), ",")
		if len(jsonTag) > 0 && jsonTag[0] != "" {
			typeFieldName = jsonTag[0]
		} else {
			typeFieldName = typeField.Name
		}

		validateTag := typeField.Tag.Get(TagName)
		validators := strings.Split(validateTag, ",")

		for _, validator := range validators {
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
					message := fmt.Sprintf(
						`Invalid parameter "%s" used in "%s" validation.`,
						args[1],
						args[0],
					)
					return nil, errors.New(message)
				}
				fieldLength = limit
			}

			var err error
			switch args[0] {
			case "notblank":
				err = notBlank(fieldValue, typeFieldTypeName)
			case "email":
				err = email(fieldValue, typeFieldTypeName)
			case "numeric":
				err = numeric(fieldValue, typeFieldTypeName)
			case "url":
				err = checkURL(fieldValue, typeFieldTypeName)
			case "required":
				err = required(fieldValue)
			case "notempty":
				err = notEmpty(fieldValue, fieldValueType)
			case "min":
				err = checkMin(fieldValue, typeFieldTypeName, fieldLength)
			case "max":
				err = checkMax(fieldValue, typeFieldTypeName, fieldLength)
			case "len":
				err = checkLen(fieldValue, fieldValueType, typeFieldTypeName, fieldLength)
			case "isarray":
				err = isArray(fieldValue, fieldValueType, typeFieldName, v.errors)
			default:
				return nil, fmt.Errorf("unknown validator: %s", args[0])
			}

			if err != nil {
				v.errors[typeFieldName] = append(v.errors[typeFieldName], err.Error())
			}
		}
	}
	return v.errors, nil
}
