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

func Required(fieldValue reflect.Value) error {
	if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
		return errors.New(ErrMsgMissing)
	}

	return nil
}

func NotBlank(fieldValue reflect.Value, typeFieldTypeName string) error {
	if fieldValue.Kind() == reflect.Ptr {
		return nil
	}

	if typeFieldTypeName != "string" {
		return errors.New(ErrMsgNotStringType)
	}

	if fieldValue.String() == "" {
		return errors.New(ErrMsgNotBlank)
	}

	return nil
}

func Email(fieldValue reflect.Value, typeFieldTypeName string) error {
	if fieldValue.Kind() == reflect.Ptr {
		return nil
	}

	if typeFieldTypeName != "string" {
		return errors.New(ErrMsgNotStringType)
	}

	if !emailRegex.MatchString(fieldValue.String()) {
		return errors.New(ErrMsgInvalidEmail)
	}

	return nil
}

func URL(fieldValue reflect.Value, typeFieldTypeName string) error {
	if fieldValue.Kind() == reflect.Ptr {
		return nil
	}

	if typeFieldTypeName != "string" {
		return errors.New(ErrMsgNotStringType)
	}

	_, err := url.ParseRequestURI(fieldValue.String())
	if err != nil {
		return errors.New(ErrMsgInvalidURL)
	}

	return nil
}

func Min(fieldValue reflect.Value, typeFieldTypeName string, fieldLength int) error {
	if fieldValue.Kind() == reflect.Ptr {
		return nil
	}

	if typeFieldTypeName == "string" {
		if len(fieldValue.String()) < fieldLength {
			return fmt.Errorf(ErrMsgStrInvalidMin, fieldLength)
		}

		return nil
	}

	if fieldValue.CanInt() {
		if fieldValue.Int() < int64(fieldLength) {
			return fmt.Errorf(ErrMsgStrInvalidInt, fieldLength)
		}
		return nil
	}

	if fieldValue.CanFloat() {
		if fieldValue.Float() < float64(fieldLength) {
			return fmt.Errorf(ErrMsgStrInvalidInt, fieldLength)
		}
		return nil
	}

	return errors.New(ErrMsgNotStrIntType)
}

func Max(fieldValue reflect.Value, typeFieldTypeName string, fieldLength int) error {
	if fieldValue.Kind() == reflect.Ptr {
		return nil
	}

	if typeFieldTypeName == "string" {
		if len(fieldValue.String()) > fieldLength {
			return fmt.Errorf(ErrMsgStrInvalidMax, fieldLength)
		}
		return nil
	}

	if fieldValue.CanInt() {
		if fieldValue.Int() > int64(fieldLength) {
			return fmt.Errorf(ErrMsgIntInvalidMax, fieldLength)
		}
		return nil
	}

	if fieldValue.CanFloat() {
		if fieldValue.Float() > float64(fieldLength) {
			return fmt.Errorf(ErrMsgIntInvalidMax, fieldLength)
		}
		return nil
	}

	return errors.New(ErrMsgNotStrIntType)
}

func NotEmpty(fieldValue reflect.Value, fieldValueType reflect.Type) error {
	if fieldValueType.Kind() != reflect.Slice {
		return errors.New(ErrMsgNotArrayType)
	}

	if value := fieldValue.Len(); value == 0 {
		return errors.New(ErrMsgEmptyArray)
	}

	return nil
}

func IsArray(fieldValue reflect.Value, fieldValueType reflect.Type, typeFieldName string) error {
	if fieldValue.Kind() == reflect.Ptr && fieldValue.IsNil() {
		return nil
	}

	if fieldValueType.Kind() != reflect.Slice {
		return errors.New(ErrMsgNotArrayType)
	}

	return nil
}

func Len(fieldValue reflect.Value, fieldValueType reflect.Type, typeFieldTypeName string, fieldLength int) error {
	if fieldValue.Kind() == reflect.Ptr {
		return nil
	}

	if typeFieldTypeName != "string" && fieldValueType.Kind() != reflect.Slice {
		return errors.New(ErrMsgNotStrSliceType)
	}

	if typeFieldTypeName == "string" && len(fieldValue.String()) != fieldLength {
		return fmt.Errorf(ErrMsgInvalidLength, fieldLength)
	}

	if fieldValueType.Kind() == reflect.Slice && fieldValue.Len() != fieldLength {
		return fmt.Errorf(ErrMsgInvalidLengthSlice, fieldLength)
	}

	return nil
}

func Numeric(fieldValue reflect.Value, typeFieldTypeName string) error {
	if fieldValue.Kind() == reflect.Ptr {
		return nil
	}

	if typeFieldTypeName != "string" {
		return errors.New(ErrMsgNotStringType)
	}

	if !numericalRegex.MatchString(fieldValue.String()) {
		return errors.New(ErrMsgNotNumeric)
	}

	return nil
}
