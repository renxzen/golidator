package engine

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/renxzen/golidator/internal/cache"
	"github.com/renxzen/golidator/internal/fieldinfo"
	"github.com/renxzen/golidator/internal/validators"
)

type ValidationError struct {
	Field  string   `json:"field"`
	Errors []string `json:"errors"`
}

var (
	UseCaching = true
	typeCache  = cache.NewTypeCache()
)

func Validate(model any) ([]ValidationError, error) {
	value := reflect.ValueOf(model)
	kind := value.Kind()

	if kind == reflect.Pointer {
		if value.IsNil() {
			return nil, nil
		}

		value = value.Elem()
		kind = value.Kind()
	}

	if kind != reflect.Struct {
		return nil, fmt.Errorf("model must be a struct, got %s", kind)
	}

	if UseCaching {
		return validateWithCache(value)
	}
	return validateWithoutCache(value)
}

func validateWithCache(value reflect.Value) ([]ValidationError, error) {
	typeInfo := typeCache.GetWithValues(value.Type(), value)
	results := make([]ValidationError, 0, len(typeInfo.Fields))

	for _, fieldInfo := range typeInfo.Fields {
		if fieldInfo.ValidateTag == "" {
			continue
		}

		validationResults, err := executeFieldValidation(fieldInfo)
		if err != nil {
			return nil, err
		}
		results = append(results, validationResults...)
	}
	return results, nil
}

func validateWithoutCache(value reflect.Value) ([]ValidationError, error) {
	numField := value.NumField()
	results := make([]ValidationError, 0, numField)

	for i := range numField {
		fieldInfo := fieldinfo.ExtractInfo(value, i)

		if fieldInfo.ValidateTag == "" {
			continue
		}

		validationResults, err := executeFieldValidation(fieldInfo)
		if err != nil {
			return nil, err
		}
		results = append(results, validationResults...)
	}
	return results, nil
}

func executeFieldValidation(fieldInfo fieldinfo.Info) ([]ValidationError, error) {
	if fieldInfo.ValidateTag == "" {
		return nil, nil
	}

	validators := strings.Split(fieldInfo.ValidateTag, ",")
	var valErrors []string
	var results []ValidationError

	for _, validator := range validators {
		validator = strings.TrimSpace(validator)
		if validator == "" {
			continue
		}

		validatorName := validator
		if idx := strings.IndexByte(validator, '='); idx != -1 {
			validatorName = validator[:idx]
		}

		errorMsg := executeValidator(validatorName, fieldInfo)
		if errorMsg != "" {
			valErrors = append(valErrors, errorMsg)
		}

		if validatorName == "isarray" && errorMsg == "" {
			nestedResults := handleArrayValidation(fieldInfo)
			results = append(results, nestedResults...)
		}
	}

	if len(valErrors) > 0 {
		results = append(results, ValidationError{
			Field:  fieldInfo.JSONName,
			Errors: valErrors,
		})
	}

	return results, nil
}

func executeValidator(validatorName string, fieldInfo fieldinfo.Info) string {
	valFunc, exists := validators.Registry[validatorName]
	if !exists {
		return fmt.Sprintf("unknown validator: %s", validatorName)
	}

	return valFunc(fieldInfo)
}

func handleArrayValidation(fieldInfo fieldinfo.Info) []ValidationError {
	var results []ValidationError
	validationValue := fieldInfo.GetValue()

	if fieldInfo.Kind == reflect.Slice {
		for j := 0; j < validationValue.Len(); j++ {
			result, err := Validate(validationValue.Index(j).Interface())
			if err != nil {
				results = append(results, ValidationError{
					Field:  fmt.Sprintf("%s[%d]", fieldInfo.JSONName, j),
					Errors: []string{err.Error()},
				})
				continue
			}

			for _, arr := range result {
				subFieldName := fmt.Sprintf("%s[%d]: %s", fieldInfo.JSONName, j, arr.Field)
				results = append(results, ValidationError{
					Field:  subFieldName,
					Errors: arr.Errors,
				})
			}
		}
	}

	return results
}
