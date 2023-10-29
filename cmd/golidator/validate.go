package golidator

import "github.com/renxzen/golidator/internal/util"

func Validate(i interface{}) ([]ValidationError, error) {
	v := NewValidate(i)
	mapErrors, err := v.GetErrors()
	errors := make([]ValidationError, 0)
	for field, arr := range mapErrors {
		err := ValidationError{
			Field:  util.ToSnakeCase(field),
			Errors: arr,
		}

		errors = append(errors, err)
	}

	return errors, err
}

