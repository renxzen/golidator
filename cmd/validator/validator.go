package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/renxzen/golidator/internal/util"
)

const TagName = "validate"

type validator struct {
	value  reflect.Value
	errors map[string][]string

	fieldValue     reflect.Value
	fieldValueType reflect.Type
	fieldLength    int

	typeField         reflect.StructField
	typeFieldName     string
	typeFieldTypeName string
}

type Validator interface {
	GetErrors() (map[string][]string, error)
}

// constructor

func NewValidate(model any) Validator {
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	return &validator{
		value:  value,
		errors: make(map[string][]string),
	}
}

// private methods

func (v *validator) setError(message string) {
	v.errors[v.typeFieldName] = append(v.errors[v.typeFieldName], message)
}

// public methods

func (v *validator) GetErrors() (map[string][]string, error) {
	reflection := reflect.ValueOf(v)

	for i := 0; i < v.value.NumField(); i++ {
		// set field data
		v.fieldValue = v.value.Field(i)
		v.fieldValueType = v.fieldValue.Type()

		v.typeField = v.value.Type().Field(i)
		v.typeFieldName = v.typeField.Name
		v.typeFieldTypeName = v.typeField.Type.Name()

		if v.fieldValue.Kind() == reflect.Ptr && !v.fieldValue.IsNil() {
			v.fieldValue = v.fieldValue.Elem()
			v.fieldValueType = v.fieldValue.Type()
			v.typeFieldTypeName = v.fieldValueType.Name()
		}

		validateTag := v.typeField.Tag.Get(TagName)
		validators := strings.Split(validateTag, ",")

		for _, validator := range validators {
			args := strings.Split(validator, "=")

			if len(args) > 1 {
				limit, err := strconv.Atoi(args[1])
				if err != nil {
					message := fmt.Sprintf(
						`Invalid parameter "%s" used in "%s" validation.`,
						args[1],
						args[0],
					)
					return v.errors, errors.New(message)
				}
				v.fieldLength = limit
			}

			// TODO: replace with a switch statement
			method := reflection.MethodByName(util.Capitalize(args[0]))
			if !method.IsValid() {
				message := fmt.Sprintf(`Validator "%s" not found.`, args[0])
				return v.errors, errors.New(message)
			}

			method.Call(nil)
		}
	}
	return v.errors, nil
}
