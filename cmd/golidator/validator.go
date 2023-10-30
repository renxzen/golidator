package golidator

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/renxzen/golidator/internal/util"
)

const TagName = "validate"

type validator struct {
	value  reflect.Value
	errors map[string][]string

	fieldName   string
	fieldType   string
	fieldValue  reflect.Value
	fieldIndex  int
	fieldLength int
}

type Validator interface {
	GetErrors() (map[string][]string, error)
}

func NewValidate(model any) Validator {
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	return &validator{
		value:  value,
		errors: make(map[string][]string),
		// validators: validators,
	}
}

func (v *validator) GetErrors() (map[string][]string, error) {
	for i := 0; i < v.value.NumField(); i++ {
		field := v.value.Type().Field(i)
		tag := field.Tag.Get(TagName)
		validators := strings.Split(tag, ",")

		// get reflect values here to avoid calling "getvalues"
		// more than one time for each validator
		v.SetValues(i)

		for _, validator := range validators {
			args := strings.Split(validator, "=")

			if len(args) > 1 {
				limit, err := strconv.Atoi(args[1])
				if err != nil {
					return v.errors, errors.New("Invalid limit number used in validation")
				}
				v.fieldLength = limit
			}

			value := reflect.ValueOf(v)
			method := value.MethodByName(util.Capitalize(args[0]))
			if !method.IsValid() {
				// return v.errors, ERROR TODEFINE
			}

			method.Call([]reflect.Value{})
		}
	}
	return v.errors, nil
}
