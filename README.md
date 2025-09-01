# GoLidator

GoLidator is a lightweight and flexible validation library for Go (Golang). It provides a simple and extensible way to validate data structures based on struct tags.

## Installation

To install GoLidator, use the `go get` command:

```bash
go get github.com/renxzen/golidator
```

## Usage

1. Import the `golidator` package in your code:

```go
import (
    "github.com/renxzen/golidator"
)
```

2. Create a struct with validation tags:

```go
type YourStruct struct {
    Field1 string `json:"field_1" validate:"notblank,email"`
    Field2 int    `json:"field_2" validate:"min=10,max=100"`
    // Add more fields and validation tags as needed
}
```

3. Use the `golidator.Validate` function to check for validation errors:

```go
data := YourStruct{
    Field1: "test@example",
    Field2: 5,
    // Set values for other fields
}

validationErrors, err := golidator.Validate(data)
if err != nil {
    // Handle error
}

if len(validationErrors) > 0 {
    // Handle validation errors
    for _, ve := range validationErrors {
        fmt.Printf("Field: %s, Errors: %v\n", ve.Field, ve.Errors)
    }
} else {
    // Data is valid, proceed with your logic
}
```

This would output:

```
Field: field_1, Errors: ["must be a valid email"]
Field: field_2, Errors: ["must be more or equal than 10"]
```

Or as a json:

```json
[
  {
    "field": "field_1",
    "errors": [
      "must be a valid email"
    ]
  },
  {
    "field": "field_2",
    "errors": [
      "must be more or equal than 10"
    ]
  }
]
```

## Performance & Caching

GoLidator includes an intelligent caching system that significantly improves performance for repeated validations of the same struct types.

### Cache Control

By default, caching is **enabled** for optimal performance. You can control caching behavior:

```go
import "github.com/renxzen/golidator"

// Disable caching (not recommended for production)
golidator.SetCaching(false)

// Enable caching (default behavior)
golidator.SetCaching(true)
```

### Performance Benefits

- **~2.4x faster** validation with caching enabled
- **Reduced memory allocations** through type information reuse
- **Automatic optimization** for repeated struct validations

## Custom Validators

GoLidator supports custom validators for specialized validation logic. Custom validators receive detailed field information and return error messages.

### Adding Custom Validators

```go
import "github.com/renxzen/golidator"

// Simple custom validator
golidator.AddValidator("custom-check", func(field golidator.FieldInfo) string {
    if !field.IsString() {
        return "field must be a string"
    }

    if field.String() != "expected-value" {
        return "field must equal 'expected-value'"
    }

    return "" // No error
})

// Custom validator with arguments
golidator.AddValidator("min-length", func(field golidator.FieldInfo) string {
    if !field.IsString() {
        return "field must be a string"
    }

    // Get the minimum length from validation tag: validate:"min-length=5"
    minLen, exists := field.GetArgumentInt("min-length")
    if !exists {
        return ""
    }

    if field.Len() < minLen {
        return fmt.Sprintf("field must be at least %d characters", minLen)
    }

    return ""
})
```

### Using Custom Validators

```go
type User struct {
    Username string `validate:"custom-check"`
    Password string `validate:"min-length=8"`
}

user := User{
    Username: "expected-value",
    Password: "short",
}

errors, err := golidator.Validate(user)
// Will return validation error for Password field
```

### FieldInfo API

Custom validators receive a `FieldInfo` struct with comprehensive field information:

```go
type FieldInfo struct {
    // Field metadata
    Name         string            // Original field name
    JSONName     string            // Name for error messages (from json tag)
    TypeName     string            // Type name ("string", "int", etc.)
    IsPointer    bool              // Whether field is a pointer type
    IsRequired   bool              // Whether field has "required" validator

    // Validation arguments
    ValidatorStrs map[string]string // String arguments from tags
    ValidatorInts map[string]int    // Pre-parsed integer arguments

    // ... other fields for advanced usage
}
```

### Helper Methods

```go
func customValidator(field golidator.FieldInfo) string {
    // Skip validation for non-required nil pointer fields
    if !field.IsRequired && field.IsNil() {
        return ""
    }

    // Type checking
    if field.IsString() { /* ... */ }
    if field.IsInt() { /* ... */ }
    if field.IsFloat() { /* ... */ }
    if field.IsSlice() { /* ... */ }

    // Nil checking
    if field.IsNil() { /* ... */ }

    // Value access
    strValue := field.String()
    intValue := field.Int()
    floatValue := field.Float()
    length := field.Len()

    // Argument access
    strArg := field.GetArgumentStr("arg-name")
    intArg, exists := field.GetArgumentInt("arg-name")

    return ""
}
```

## Supported Validators

- `notblank`: Ensures that a string is not empty.
- `email`: Validates that a string is a valid email address.
- `numeric`: Validates that a string contains only numbers.
- `url`: Validates that a string is a valid URL.
- `required`: Ensures that a field is not missing from the body.
- `notempty`: Ensures that an array is not empty.
- `min`: Validates that a string or numeric value is greater than or equal to a specified limit.
- `max`: Validates that a string or numeric value is less than or equal to a specified limit.
- `len`: Validates that a string or a slice value has the same amount of characters or elements.
- `isarray`: Ensures that a field is a non-nil slice and validates its elements recursively.

## TODO

- [x] optimize for speed
- [x] map for validators
- [x] cache to speed up performance
- [x] custom validators

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
