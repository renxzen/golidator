package golidator

import (
	neturl "net/url"
	"fmt"
	"reflect"
	"regexp"
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

func (v *validator) setError(message string)  {
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

func (v *validator) ValidEmail(str string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(str)
}

func (v *validator) Email() {
	if v.fieldType != "string" {
		v.setError("Invalid type. Must be string")
		return
	}

	if !v.ValidEmail(v.fieldValue.String()) {
		v.setError("Must be a valid email")
		return
	}
}

func (v *validator) ValidUrl(str string) bool {
	_, err := neturl.ParseRequestURI(str)
	return err == nil
}


func (v *validator) Url() {
	if v.fieldType != "string" {
		v.setError("Invalid type. Must be string")
		return
	}

	if !v.ValidUrl(v.fieldValue.String()) {
		v.setError("Must be a valid url")
		return
	}
}

func (v *validator) Min() {
	if v.fieldType != "string" {
		v.setError("Invalid type. Must be string")
		return
	}

	if len(v.fieldValue.String()) < v.fieldLength {
		v.setError(fmt.Sprintf("Must be at least %v characters long", v.fieldLength))
		return
	}
}

func (v *validator) Max() {
	if v.fieldType != "string" {
		v.setError("Invalid type. Must be string")
		return
	}

	if len(v.fieldValue.String()) > v.fieldLength {
		v.setError(fmt.Sprintf("Must be at least %v characters long", v.fieldLength))
	}

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

func (v *validator) Valarray() {
	field := v.value.Type().Field(v.fieldIndex).Name
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

	leni := array.Len()
	for i := 0; i < leni; i++ {
		summary, err := Validate(array.Index(i).Interface())
		if err != nil {
			// TODO: do something with error
			return
		}

		for j := range summary {
			summaryField := fmt.Sprintf("%v[%v]: %v", field, i, summary[j].Field)
			v.errors[summaryField] = summary[j].Errors
		}
	}
}
