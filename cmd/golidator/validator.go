package golidator

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field  string   `json:"field"`
	Errors []string `json:"errors"`
}

type validator struct {
	tagName    string
	value      reflect.Value
	errors     map[string][]string
	validators map[string]func(map[string][]string, reflect.Value, int, int)
}

type Validator interface {
	GetErrors() (map[string][]string, error)
}

func NewValidate(model interface{}) Validator {
	value := reflect.ValueOf(model)
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	validators := make(map[string]func(map[string][]string, reflect.Value, int, int))
	validators["required"] = required
	validators["notblank"] = NotBlank
	validators["email"] = email
	validators["url"] = url
	validators["min"] = min
	validators["max"] = max
	validators["notempty"] = notempty
	validators["valarray"] = valarray

	return &validator{
		tagName:    "validate",
		value:      value,
		errors:     make(map[string][]string),
		validators: validators,
	}
}

func (v *validator) GetErrors() (map[string][]string, error) {
	for i := 0; i < v.value.NumField(); i++ {
				field := v.value.Type().Field(i)
		tag := field.Tag.Get(v.tagName)
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

			fn := v.validators[args[0]]
			if fn != nil {
				fn(v.errors, v.value, i, limit)
			}
		}
	}
	return v.errors, nil
}

