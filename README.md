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
import "github.com/renxzen/golidator"
```

2. Create a struct with validation tags:

```go
type YourStruct struct {
    Field1 string `validate:"notblank,email"`
    Field2 int    `validate:"min=10,max=100"`
    // Add more fields and validation tags as needed
}
```

3. Use the `golidator.Validate` function to check for validation errors:

```go
data := YourStruct{
    Field1: "test@example.com",
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
- `isarray`: Ensures that a field is a non-nil array and validates its elements recursively.

## Todo

- custom validators
- validate array without `isarray`
- optimize reflection in loops
- simplify conditional checks
- reduce redundancy with type errors

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
