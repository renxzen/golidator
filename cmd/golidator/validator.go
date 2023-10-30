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
	value      reflect.Value
	errors     map[string][]string
	// validators map[string]func(map[string][]string, reflect.Value, int, int)

	fieldIndex int
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

	// validators := make(map[string]func(map[string][]string, reflect.Value, int, int))
	// validators["required"] = required
	// validators["notblank"] = NotBlank
	// validators["email"] = email
	// validators["url"] = url
	// validators["min"] = min
	// validators["max"] = max
	// validators["notempty"] = notempty
	// validators["valarray"] = valarray

	return &validator{
		value:      value,
		errors:     make(map[string][]string),
		// validators: validators,
	}
}

func (v *validator) GetErrors() (map[string][]string, error) {
	for i := 0; i < v.value.NumField(); i++ {
		field := v.value.Type().Field(i)
		tag := field.Tag.Get(TagName)
		validators := strings.Split(tag, ",")

		for _, validator := range validators {
			args := strings.Split(validator, "=")

			var err error
			var limit int

			if len(args) > 1 {
				limit, err = strconv.Atoi(args[1])
				if err != nil {
					return v.errors, errors.New("Invalid limit number used in validation")
				}
			}

			v.fieldIndex = i
			v.fieldLength = limit

			value := reflect.ValueOf(v)
			method := value.MethodByName(util.Capitalize(args[0]))
			if !method.IsValid() {
				// return v.errors, ERROR TODEFINE
			}

			method.Call([]reflect.Value{})
			// fn := v.validators[args[0]]
			// if fn != nil {
			// 	fn(v.errors, v.value, i, limit)
			// }
		}
	}
	return v.errors, nil
}
