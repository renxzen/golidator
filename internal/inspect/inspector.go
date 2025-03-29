package inspect

import (
	"errors"
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

type inspect struct {
	value  reflect.Value
	errors map[string][]string
}

func NewInspector(model any) *inspect {
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	return &inspect{
		value:  value,
		errors: make(map[string][]string),
	}
}

func (v *inspect) GetErrors() (map[string][]string, error) {
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
				err = validate.IsArray(fieldValue, fieldValueType, typeFieldName)

				for i := 0; i < fieldValue.Len(); i++ {
					result, subErr := NewInspector(fieldValue.Index(i).Interface()).GetErrors()
					if subErr != nil {
						v.errors[typeFieldName] = append(v.errors[typeFieldName], subErr.Error())
						continue
					}

					for subField, arr := range result {
						subFieldName := fmt.Sprintf("%v[%v]: %v", typeFieldName, i, subField)
						v.errors[subFieldName] = arr
					}
				}
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
