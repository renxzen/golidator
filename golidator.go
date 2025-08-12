package golidator

import (
	"github.com/renxzen/golidator/internal/engine"
	"github.com/renxzen/golidator/internal/fieldinfo"
	"github.com/renxzen/golidator/internal/validators"
)

const (
	ValidateTag = "validate"
	JsonTag     = "json"
)

// ValidationError represents a validation error for a field
type ValidationError = engine.ValidationError

// ValidatorFunc represents a validator function
type ValidatorFunc = validators.ValidatorFunc

// FieldInfo represents field information for validation
type FieldInfo = fieldinfo.Info

// Validate validates a struct and returns validation errors
func Validate(model any) ([]ValidationError, error) {
	return engine.Validate(model)
}

// SetCaching enables or disables type caching for validation
func SetCaching(enabled bool) {
	engine.UseCaching = enabled
}

// AddValidator adds a new validator to the registry
func AddValidator(name string, validator ValidatorFunc) {
	if name == "" || validator == nil {
		panic("validator name cannot be empty")
	}
	validators.Registry[name] = validator
}
