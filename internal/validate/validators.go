package validate

import (
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
)

var (
	emailRegex     = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	numericalRegex = regexp.MustCompile(`^[0-9]+$`)
)

func Required(value reflect.Value, kind reflect.Kind) error {
	if kind == reflect.Ptr && value.IsNil() {
		return errors.New(ErrMsgMissing)
	}

	return nil
}

func NotBlank(value reflect.Value, kind reflect.Kind, fieldType string) error {
	if kind == reflect.Ptr {
		return nil
	}

	if fieldType != "string" {
		return errors.New(ErrMsgNotStringType)
	}

	if value.String() == "" {
		return errors.New(ErrMsgNotBlank)
	}

	return nil
}

func Email(value reflect.Value, kind reflect.Kind, fieldType string) error {
	if kind == reflect.Ptr {
		return nil
	}

	if fieldType != "string" {
		return errors.New(ErrMsgNotStringType)
	}

	if !emailRegex.MatchString(value.String()) {
		return errors.New(ErrMsgInvalidEmail)
	}

	return nil
}

func URL(value reflect.Value, kind reflect.Kind, fieldType string) error {
	if kind == reflect.Ptr {
		return nil
	}

	if fieldType != "string" {
		return errors.New(ErrMsgNotStringType)
	}

	if _, err := url.ParseRequestURI(value.String()); err != nil {
		return errors.New(ErrMsgInvalidURL)
	}

	return nil
}

func Min(value reflect.Value, kind reflect.Kind, fieldType string, fieldLength int) error {
	if kind == reflect.Ptr {
		return nil
	}

	if fieldType == "string" {
		if len(value.String()) < fieldLength {
			return fmt.Errorf(ErrMsgStrInvalidMin, fieldLength)
		}

		return nil
	}

	if value.CanInt() {
		if value.Int() < int64(fieldLength) {
			return fmt.Errorf(ErrMsgStrInvalidInt, fieldLength)
		}
		return nil
	}

	if value.CanFloat() {
		if value.Float() < float64(fieldLength) {
			return fmt.Errorf(ErrMsgStrInvalidInt, fieldLength)
		}
		return nil
	}

	return errors.New(ErrMsgNotStrIntType)
}

func Max(value reflect.Value, kind reflect.Kind, fieldType string, fieldLength int) error {
	if kind == reflect.Ptr {
		return nil
	}

	if fieldType == "string" {
		if len(value.String()) > fieldLength {
			return fmt.Errorf(ErrMsgStrInvalidMax, fieldLength)
		}
		return nil
	}

	if value.CanInt() {
		if value.Int() > int64(fieldLength) {
			return fmt.Errorf(ErrMsgIntInvalidMax, fieldLength)
		}
		return nil
	}

	if value.CanFloat() {
		if value.Float() > float64(fieldLength) {
			return fmt.Errorf(ErrMsgIntInvalidMax, fieldLength)
		}
		return nil
	}

	return errors.New(ErrMsgNotStrIntType)
}

func NotEmpty(value reflect.Value, typeValue reflect.Type, kind reflect.Kind) error {
	if kind != reflect.Slice {
		return errors.New(ErrMsgNotArrayType)
	}

	if length := value.Len(); length == 0 {
		return errors.New(ErrMsgEmptyArray)
	}

	return nil
}

func IsArray(value reflect.Value, typeValue reflect.Type, kind reflect.Kind) error {
	if kind == reflect.Ptr && value.IsNil() {
		return nil
	}

	if typeValue.Kind() != reflect.Slice {
		return errors.New(ErrMsgNotArrayType)
	}

	return nil
}

func Len(value reflect.Value, typeValue reflect.Type, kind reflect.Kind, fieldType string, fieldLength int) error {
	if kind == reflect.Ptr {
		return nil
	}

	if fieldType != "string" && typeValue.Kind() != reflect.Slice {
		return errors.New(ErrMsgNotStrSliceType)
	}

	if fieldType == "string" && len(value.String()) != fieldLength {
		return fmt.Errorf(ErrMsgInvalidLength, fieldLength)
	}

	if typeValue.Kind() == reflect.Slice && value.Len() != fieldLength {
		return fmt.Errorf(ErrMsgInvalidLengthSlice, fieldLength)
	}

	return nil
}

func Numeric(value reflect.Value, kind reflect.Kind, fieldType string) error {
	if kind == reflect.Ptr {
		return nil
	}

	if fieldType != "string" {
		return errors.New(ErrMsgNotStringType)
	}

	if !numericalRegex.MatchString(value.String()) {
		return errors.New(ErrMsgNotNumeric)
	}

	return nil
}
