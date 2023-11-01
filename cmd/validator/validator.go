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
	}
}

func (v *validator) GetErrors() (map[string][]string, error) {
	for i := 0; i < v.value.NumField(); i++ {
		field := v.value.Type().Field(i)
		validateTag := field.Tag.Get(TagName)
		validators := strings.Split(validateTag, ",")

		v.SetValues(i)
		for _, validator := range validators {
			args := strings.Split(validator, "=")

			if len(args) > 1 {
				limit, err := strconv.Atoi(args[1])
				if err != nil {
					message := fmt.Sprintf(`Invalid parameter "%s" used in "%s" validation.`, args[1], args[0])
					return v.errors, errors.New(message)
				}
				v.fieldLength = limit
			}

			value := reflect.ValueOf(v)
			method := value.MethodByName(util.Capitalize(args[0]))
			if !method.IsValid() {
				message := fmt.Sprintf(`Validator "%s" not found.`, args[0])
				return v.errors, errors.New(message)
			}

			method.Call([]reflect.Value{})
		}
	}
	return v.errors, nil
}

// validators


