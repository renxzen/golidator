package golidator

import (
	"github.com/renxzen/golidator/cmd/validator"
)

type ValidationError struct {
	Field  string   `json:"field"`
	Errors []string `json:"errors"`
}

func Validate(i any) ([]ValidationError, error) {
	v := validator.NewValidate(i)
	mapErrors, err := v.GetErrors()
	errors := make([]ValidationError, 0)
	for field, arr := range mapErrors {
		err := ValidationError{
			Field:  field,
			Errors: arr,
		}

		errors = append(errors, err)
	}

	return errors, err
}
