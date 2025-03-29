package golidator

import (
	"github.com/renxzen/golidator/internal/inspect"
)

type ValidationError struct {
	Field  string   `json:"field"`
	Errors []string `json:"errors"`
}

func Validate(i any) ([]ValidationError, error) {
	inspector := inspect.NewInspector(i)
	mapErrors, err := inspector.GetErrors()
	if err != nil {
		return nil, err
	}

	valErrors := make([]ValidationError, 0, len(mapErrors))
	for field, arr := range mapErrors {
		valError := ValidationError{
			Field:  field,
			Errors: arr,
		}

		valErrors = append(valErrors, valError)
	}

	return valErrors, nil
}
