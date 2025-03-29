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

	fieldValue     reflect.Value
	fieldValueType reflect.Type
	fieldLength    int

	typeField         reflect.StructField
	typeFieldName     string
	typeFieldTypeName string
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

func (v *Validator) setError(message string) {
	v.errors[v.typeFieldName] = append(v.errors[v.typeFieldName], message)
}

func (v *Validator) setFieldData(i int) {
	// get field value
	v.fieldValue = v.value.Field(i)
	v.fieldValueType = v.fieldValue.Type()

	// get field type
	v.typeField = v.value.Type().Field(i)
	v.typeFieldTypeName = v.typeField.Type.Name()

	// check if field is a pointer
	if v.fieldValue.Kind() == reflect.Ptr && !v.fieldValue.IsNil() {
		v.fieldValue = v.fieldValue.Elem()
		v.fieldValueType = v.fieldValue.Type()
		v.typeFieldTypeName = v.fieldValueType.Name()
	}

	// get json name
	jsonTag := strings.Split(v.typeField.Tag.Get(JsonTag), ",")
	if len(jsonTag) > 0 && jsonTag[0] != "" {
		v.typeFieldName = jsonTag[0]
	} else {
		// get field name
		v.typeFieldName = v.typeField.Name
	}
}

func (v *Validator) GetErrors() (map[string][]string, error) {
	for i := 0; i < v.value.NumField(); i++ {
		v.setFieldData(i)
		validateTag := v.typeField.Tag.Get(TagName)
		validators := strings.Split(validateTag, ",")

		for _, validator := range validators {
			if validator == "" {
				continue
			}

			args := strings.Split(validator, "=")
			if len(args) == 0 {
				continue
			}

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
				v.fieldLength = limit
			}

			switch args[0] {
			case "notblank":
				v.NotBlank()
			case "email":
				v.Email()
			case "numeric":
				v.Numeric()
			case "url":
				v.URL()
			case "required":
				v.Required()
			case "notempty":
				v.NotEmpty()
			case "min":
				v.Min()
			case "max":
				v.Max()
			case "len":
				v.Len()
			case "isarray":
				v.IsArray()
			default:
				return nil, fmt.Errorf("unknown validator: %s", args[0])
			}
		}
	}
	return v.errors, nil
}
