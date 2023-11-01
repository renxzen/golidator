package validator

import (
	"fmt"
	"net/url"
	"reflect"
	"regexp"

	"github.com/renxzen/golidator/internal/util"
)

func (v *validator) SetValues(i int) {
	v.fieldIndex = i
	v.fieldName = v.value.Type().Field(v.fieldIndex).Name
	v.fieldType = v.value.Type().Field(v.fieldIndex).Type.Name()
	v.fieldValue = v.value.Field(v.fieldIndex)

	if v.value.Field(v.fieldIndex).Kind() == reflect.Ptr {
		if v.value.Field(v.fieldIndex).IsNil() {
			return
		}
		v.fieldType = v.value.Field(v.fieldIndex).Elem().Type().Name()
		v.fieldValue = v.value.Field(v.fieldIndex).Elem()
	}
}

func (v *validator) setError(message string) {
	v.errors[v.fieldName] = append(v.errors[v.fieldName], message)
}

func (v *validator) Required() {
	if v.value.Field(v.fieldIndex).IsNil() {
		v.setError("Must not be missing from the body")
	}
}

func (v *validator) Notblank() {
	if v.fieldType != "string" {
		v.setError("Invalid type. Must be string")
		return
	}

	if v.fieldValue.String() == "" {
		v.setError("Must not be blank")
		return
	}
}

func (v *validator) Email() {
	if v.fieldType != "string" {
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
	if v.fieldType != "string" {
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
	if v.fieldType == "string" {
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

	v.setError("Invalid type. Must be string or numeric")
	return
}

func (v *validator) Max() {
	if v.fieldType == "string" {
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

	v.setError("Invalid type. Must be string or numeric")
	return
}

func (v *validator) Notempty() {
	if v.value.Field(v.fieldIndex).Type().Kind() != reflect.Slice {
		return
	}

	value := v.value.Field(v.fieldIndex).Len()
	if value == 0 {
		v.setError("Array must not be empty")
		return
	}
}

func (v *validator) Isarray() {
	fieldName := v.value.Type().Field(v.fieldIndex).Name
	array := v.value.Field(v.fieldIndex)
	if array.Kind() == reflect.Ptr {
		if array.IsNil() {
			return
		}

		array = array.Elem()
	}

	if array.Type().Kind() != reflect.Slice {
		v.setError("Invalid type. Must be array")
		return
	}

	for i := 0; i < array.Len(); i++ {
		mapErrors, err := NewValidate(array.Index(i).Interface()).GetErrors()
		if err != nil {
			// TODO: do something with error
			return
		}

		for subField, arr := range mapErrors {
			subFieldName := fmt.Sprintf("%v[%v]: %v", fieldName, i, util.ToSnakeCase(subField))
			v.errors[subFieldName] = arr
		}
	}
}
