package fieldinfo

import (
	"reflect"
)

// Info contains comprehensive information about a struct field for validation purposes.
// This struct is used by custom validators to access field metadata and values.
type Info struct {
	// Index is the zero-based position of this field within the struct definition.
	Index int

	// Name is the original field name as declared in the Go struct.
	Name string

	// JSONName is the field name used in validation error messages.
	// Derived from the "json" struct tag, or falls back to the original field name.
	JSONName string

	// Type is the reflect.Type of the field after dereferencing any pointer types.
	// For *string, this would be the string type.
	Type reflect.Type

	// Kind is the reflect.Kind of the field after dereferencing any pointer types.
	// For *string, this would be reflect.String.
	Kind reflect.Kind

	// TypeName is the string representation of the field's type name.
	// Examples: "string", "int", "bool", "MyCustomType"
	TypeName string

	// ValidateTag contains the complete validation tag string from the struct field.
	// Example: "required,min=5,max=100,email"
	ValidateTag string

	// IsPointer indicates whether the original field declaration is a pointer type.
	// True for *string, *int, etc. False for string, int, etc.
	IsPointer bool

	// OriginalKind is the reflect.Kind before dereferencing pointer types.
	// For *string, this would be reflect.Pointer.
	OriginalKind reflect.Kind

	// Value is the actual reflect.Value of the field from the struct instance.
	// This contains the runtime value being validated.
	Value reflect.Value

	// ValidatorStrs contains string arguments parsed from validation tags.
	// For validate:"custom=hello", this would contain {"custom": "hello"}
	ValidatorStrs map[string]string

	// ValidatorInts contains pre-parsed integer arguments from validation tags.
	// For validate:"min=10", this would contain {"min": 10}
	// This optimization avoids repeated string-to-int conversions during validation.
	ValidatorInts map[string]int

	// IsRequired indicates whether the field has the "required" validator.
	// This allows validators to skip validation on nil pointer fields that are not required.
	IsRequired bool
}

// GetValue returns the underlying reflect.Value of the field.
// If the field is a pointer and not nil, it returns the dereferenced value.
// Otherwise, it returns the original value.
func (f Info) GetValue() reflect.Value {
	if f.IsPointer && !f.Value.IsNil() {
		return f.Value.Elem()
	}
	return f.Value
}

// GetKind returns the reflect.Kind of the field.
// If the field is a pointer, it returns the Kind of the pointer type.
// Otherwise, it returns the original Kind of the field.
func (f Info) GetKind() reflect.Kind {
	if f.IsPointer {
		return f.Kind
	}
	return f.OriginalKind
}

// IsNil checks if the field value is nil.
// It returns false if the field is not a pointer.
// If the field is a pointer, it returns true if the value is nil.
func (f Info) IsNil() bool {
	if !f.IsPointer {
		return false
	}
	return f.Value.IsNil()
}

// IsInt checks if the field value can be converted to an int.
func (f Info) IsInt() bool {
	return f.GetValue().CanInt()
}

// Int returns the int value of the field.
func (f Info) Int() int64 {
	return f.GetValue().Int()
}

// IsFloat checks if the field value can be converted to a float.
func (f Info) IsFloat() bool {
	return f.GetValue().CanFloat()
}

// Float returns the float value of the field.
func (f Info) Float() float64 {
	return f.GetValue().Float()
}

// IsSlice checks if the field is a slice.
func (f Info) IsSlice() bool {
	return f.GetKind() == reflect.Slice
}

// Len returns the length of the field.
// It returns -1 if the field is not a string or slice.
// Otherwise, it returns the length of the string or slice.
func (f Info) Len() int {
	if f.IsString() {
		return len(f.GetValue().String())
	}
	if f.IsSlice() {
		return f.GetValue().Len()
	}
	return -1
}

// IsString checks if the field is a string.
func (f Info) IsString() bool {
	return f.TypeName == "string"
}

// String returns the string value of the field.
func (f Info) String() string {
	return f.GetValue().String()
}

// GetArgumentStr returns the string value associated with the given validator name
// from the Info's ValidatorStrs map. If the validator does not exist, it returns an empty string.
func (f Info) GetArgumentStr(validatorName string) string {
	if value, exists := f.ValidatorStrs[validatorName]; exists {
		return value
	}
	return ""
}

// GetArgumentInt returns the integer value associated with the given validator name
// from the Info's ValidatorInts map and a boolean indicating existence.
// If the validator does not exist, the returned boolean will be false.
func (f Info) GetArgumentInt(validatorName string) (int, bool) {
	value, exists := f.ValidatorInts[validatorName]
	return value, exists
}
