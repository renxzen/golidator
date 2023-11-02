package validator

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"

	"github.com/renxzen/golidator/internal/util"
)

func (v *validator) Required() {
	if v.fieldValue.IsNil() {
		v.setError("Must not be missing from the body")
	}
}

func (v *validator) Notblank() {
	if v.typeFieldTypeName != "string" {
		v.setError("Invalid type. Must be string")
		return
	}

	if v.fieldValue.String() == "" {
		v.setError("Must not be blank")
		return
	}
}

func (v *validator) Email() {
	if v.typeFieldTypeName != "string" {
		v.setError("Invalid type. Must be string")
		return
	}

	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	if !emailRegex.MatchString(v.fieldValue.String()) {
		v.setError("Must be a valid email")
		return
	}
}

func (v *validator) Url() {
	if v.typeFieldTypeName != "string" {
		v.setError("Invalid type. Must be string")
		return
	}

	_, err := url.ParseRequestURI(v.fieldValue.String())
	if err != nil {
		v.setError("Must be a valid url")
		return
	}
}

func (v *validator) Min() {
	if v.typeFieldTypeName == "string" {
		if len(v.fieldValue.String()) < v.fieldLength {
			v.setError(fmt.Sprintf("Must have more than %v characters", v.fieldLength))
		}
		return
	}

	if v.fieldValue.CanInt() {
		if v.fieldValue.Int() < int64(v.fieldLength) {
			v.setError(fmt.Sprintf("Must be more than %v", v.fieldLength))
		}
		return
	}

	if v.fieldValue.CanFloat() {
		if v.fieldValue.Float() < float64(v.fieldLength) {
			v.setError(fmt.Sprintf("Must be more than %v", v.fieldLength))
		}
		return
	}

	v.setError("Invalid type. Must be string or numeric")
	return
}

func (v *validator) Max() {
	if v.typeFieldTypeName == "string" {
		if len(v.fieldValue.String()) > v.fieldLength {
			v.setError(fmt.Sprintf("Must have less than %v characters", v.fieldLength))
		}
		return
	}

	if v.fieldValue.CanInt() {
		if v.fieldValue.Int() > int64(v.fieldLength) {
			v.setError(fmt.Sprintf("Must be less than %v", v.fieldLength))
		}
		return
	}

	if v.fieldValue.CanFloat() {
		if v.fieldValue.Float() < float64(v.fieldLength) {
			v.setError(fmt.Sprintf("Must be less than %v", v.fieldLength))
		}
		return
	}

	v.setError("Invalid type. Must be string or numeric")
	return
}

func (v *validator) Notempty() {
	if v.fieldValueType.Kind() != reflect.Slice {
		return
	}

	value := v.fieldValue.Len()
	if value == 0 {
		v.setError("Array must not be empty")
		return
	}
}

func (v *validator) Isarray() {
	if v.fieldValue.Kind() == reflect.Ptr {
		if v.fieldValue.IsNil() {
			return
		}

		v.fieldValue = v.fieldValue.Elem()
	}

	if v.fieldValueType.Kind() != reflect.Slice {
		v.setError("Invalid type. Must be array")
		return
	}

	for i := 0; i < v.fieldValue.Len(); i++ {
		mapErrors, err := NewValidate(v.fieldValue.Index(i).Interface()).GetErrors()
		if err != nil {
			// TODO: do something with error
			return
		}

		for subField, arr := range mapErrors {
			subFieldName := fmt.Sprintf("%v[%v]: %v", v.typeFieldName, i, util.ToSnakeCase(subField))
			v.errors[subFieldName] = arr
		}
	}
}
