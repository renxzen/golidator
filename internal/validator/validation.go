package validator

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"
)

var (
	emailRegex     = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	numericalRegex = regexp.MustCompile(`^[0-9]+$`)
)

func (v *Validator) Required() {
	if v.fieldValue.Kind() == reflect.Ptr && v.fieldValue.IsNil() {
		v.setError(ErrMsgMissing)
	}
}

func (v *Validator) NotBlank() {
	if v.fieldValue.Kind() == reflect.Ptr {
		return
	}

	if v.typeFieldTypeName != "string" {
		v.setError(ErrMsgNotStringType)
		return
	}

	if v.fieldValue.String() == "" {
		v.setError(ErrMsgNotBlank)
		return
	}
}

func (v *Validator) Email() {
	if v.fieldValue.Kind() == reflect.Ptr {
		return
	}

	if v.typeFieldTypeName != "string" {
		v.setError(ErrMsgNotStringType)
		return
	}

	if !emailRegex.MatchString(v.fieldValue.String()) {
		v.setError(ErrMsgInvalidEmail)
		return
	}
}

func (v *Validator) URL() {
	if v.fieldValue.Kind() == reflect.Ptr {
		return
	}

	if v.typeFieldTypeName != "string" {
		v.setError(ErrMsgNotStringType)
		return
	}

	_, err := url.ParseRequestURI(v.fieldValue.String())
	if err != nil {
		v.setError(ErrMsgInvalidURL)
		return
	}
}

func (v *Validator) Min() {
	if v.fieldValue.Kind() == reflect.Ptr {
		return
	}

	if v.typeFieldTypeName == "string" {
		if len(v.fieldValue.String()) < v.fieldLength {
			v.setError(fmt.Sprintf(ErrMsgStrInvalidMin, v.fieldLength))
		}
		return
	}

	if v.fieldValue.CanInt() {
		if v.fieldValue.Int() < int64(v.fieldLength) {
			v.setError(fmt.Sprintf(ErrMsgStrInvalidInt, v.fieldLength))
		}
		return
	}

	if v.fieldValue.CanFloat() {
		if v.fieldValue.Float() < float64(v.fieldLength) {
			v.setError(fmt.Sprintf(ErrMsgStrInvalidInt, v.fieldLength))
		}
		return
	}

	v.setError(ErrMsgNotStrIntType)
}

func (v *Validator) Max() {
	if v.fieldValue.Kind() == reflect.Ptr {
		return
	}

	if v.typeFieldTypeName == "string" {
		if len(v.fieldValue.String()) > v.fieldLength {
			v.setError(fmt.Sprintf(ErrMsgStrInvalidMax, v.fieldLength))
		}
		return
	}

	if v.fieldValue.CanInt() {
		if v.fieldValue.Int() > int64(v.fieldLength) {
			v.setError(fmt.Sprintf(ErrMsgIntInvalidMax, v.fieldLength))
		}
		return
	}

	if v.fieldValue.CanFloat() {
		if v.fieldValue.Float() > float64(v.fieldLength) {
			v.setError(fmt.Sprintf(ErrMsgIntInvalidMax, v.fieldLength))
		}
		return
	}

	v.setError(ErrMsgNotStrIntType)
}

func (v *Validator) NotEmpty() {
	if v.fieldValueType.Kind() != reflect.Slice {
		v.setError(ErrMsgNotArrayType)
		return
	}

	value := v.fieldValue.Len()
	if value == 0 {
		v.setError(ErrMsgEmptyArray)
		return
	}
}

func (v *Validator) IsArray() {
	if v.fieldValue.Kind() == reflect.Ptr && v.fieldValue.IsNil() {
		return
	}

	if v.fieldValueType.Kind() != reflect.Slice {
		v.setError(ErrMsgNotArrayType)
		return
	}

	for i := 0; i < v.fieldValue.Len(); i++ {
		mapErrors, err := NewValidate(v.fieldValue.Index(i).Interface()).GetErrors()
		if err != nil {
			// TODO: do something to notify the error
			return
		}

		for subField, arr := range mapErrors {
			subFieldName := fmt.Sprintf(
				"%v[%v]: %v",
				v.typeFieldName,
				i,
				subField,
			)
			v.errors[subFieldName] = arr
		}
	}
}

func (v *Validator) Len() {
	if v.fieldValue.Kind() == reflect.Ptr {
		return
	}

	if v.typeFieldTypeName != "string" && v.fieldValueType.Kind() != reflect.Slice {
		v.setError(ErrMsgNotStrSliceType)
		return
	}

	if v.typeFieldTypeName == "string" && len(v.fieldValue.String()) != v.fieldLength {
		v.setError(fmt.Sprintf(ErrMsgInvalidLength, v.fieldLength))
		return
	}

	if v.fieldValueType.Kind() == reflect.Slice && v.fieldValue.Len() != v.fieldLength {
		v.setError(fmt.Sprintf(ErrMsgInvalidLengthSlice, v.fieldLength))
		return
	}
}

func (v *Validator) Numeric() {
	if v.fieldValue.Kind() == reflect.Ptr {
		return
	}

	if v.typeFieldTypeName != "string" {
		v.setError(ErrMsgNotStringType)
		return
	}

	if !numericalRegex.MatchString(v.fieldValue.String()) {
		v.setError(ErrMsgNotNumeric)
		return
	}
}
