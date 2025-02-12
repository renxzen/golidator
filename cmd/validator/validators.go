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

func (v *validator) Required() {
	if v.fieldValue.Kind() == reflect.Ptr && v.fieldValue.IsNil() {
		v.setError(REQUIRED_ERROR)
	}
}

func (v *validator) Notblank() {
	if v.fieldValue.Kind() == reflect.Ptr {
		return
	}

	if v.typeFieldTypeName != "string" {
		v.setError(NOTSTRING_ERROR)
		return
	}

	if v.fieldValue.String() == "" {
		v.setError(NOTBLANK_ERROR)
		return
	}
}

func (v *validator) Email() {
	if v.fieldValue.Kind() == reflect.Ptr {
		return
	}

	if v.typeFieldTypeName != "string" {
		v.setError(NOTSTRING_ERROR)
		return
	}

	if !emailRegex.MatchString(v.fieldValue.String()) {
		v.setError(EMAIL_ERROR)
		return
	}
}

func (v *validator) Url() {
	if v.fieldValue.Kind() == reflect.Ptr {
		return
	}

	if v.typeFieldTypeName != "string" {
		v.setError(NOTSTRING_ERROR)
		return
	}

	_, err := url.ParseRequestURI(v.fieldValue.String())
	if err != nil {
		v.setError(URL_ERROR)
		return
	}
}

func (v *validator) Min() {
	if v.fieldValue.Kind() == reflect.Ptr {
		return
	}

	if v.typeFieldTypeName == "string" {
		if len(v.fieldValue.String()) < v.fieldLength {
			v.setError(fmt.Sprintf(MIN_STRING_ERROR, v.fieldLength))
		}
		return
	}

	if v.fieldValue.CanInt() {
		if v.fieldValue.Int() < int64(v.fieldLength) {
			v.setError(fmt.Sprintf(MIN_INTEGER_ERROR, v.fieldLength))
		}
		return
	}

	if v.fieldValue.CanFloat() {
		if v.fieldValue.Float() < float64(v.fieldLength) {
			v.setError(fmt.Sprintf(MIN_INTEGER_ERROR, v.fieldLength))
		}
		return
	}

	v.setError(NOTSTRINGORINTEGER_ERROR)
}

func (v *validator) Max() {
	if v.fieldValue.Kind() == reflect.Ptr {
		return
	}

	if v.typeFieldTypeName == "string" {
		if len(v.fieldValue.String()) > v.fieldLength {
			v.setError(fmt.Sprintf(MAX_STRING_ERROR, v.fieldLength))
		}
		return
	}

	if v.fieldValue.CanInt() {
		if v.fieldValue.Int() > int64(v.fieldLength) {
			v.setError(fmt.Sprintf(MAX_INTEGER_ERROR, v.fieldLength))
		}
		return
	}

	if v.fieldValue.CanFloat() {
		if v.fieldValue.Float() > float64(v.fieldLength) {
			v.setError(fmt.Sprintf(MAX_INTEGER_ERROR, v.fieldLength))
		}
		return
	}

	v.setError(NOTSTRINGORINTEGER_ERROR)
	return
}

func (v *validator) Notempty() {
	if v.fieldValueType.Kind() != reflect.Slice {
		v.setError(NOTARRAY_ERROR)
		return
	}

	value := v.fieldValue.Len()
	if value == 0 {
		v.setError(NOTEMPTY_ERROR)
		return
	}
}

func (v *validator) Isarray() {
	if v.fieldValue.Kind() == reflect.Ptr && v.fieldValue.IsNil() {
		return
	}

	if v.fieldValueType.Kind() != reflect.Slice {
		v.setError(NOTARRAY_ERROR)
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

func (v *validator) Len() {
	if v.fieldValue.Kind() == reflect.Ptr {
		return
	}

	if v.typeFieldTypeName != "string" && v.fieldValueType.Kind() != reflect.Slice {
		v.setError(NOTSTRINGORSLICE_ERROR)
		return
	}

	if v.typeFieldTypeName == "string" && len(v.fieldValue.String()) != v.fieldLength {
		v.setError(fmt.Sprintf(LEN_STRING_ERROR, v.fieldLength))
		return
	}

	if v.fieldValueType.Kind() == reflect.Slice && v.fieldValue.Len() != v.fieldLength {
		v.setError(fmt.Sprintf(LEN_SLICE_ERROR, v.fieldLength))
		return
	}
}

func (v *validator) Numeric() {
	if v.fieldValue.Kind() == reflect.Ptr {
		return
	}

	if v.typeFieldTypeName != "string" {
		v.setError(NOTSTRING_ERROR)
		return
	}

	if !numericalRegex.MatchString(v.fieldValue.String()) {
		v.setError(NUMERIC_ERROR)
		return
	}
}
