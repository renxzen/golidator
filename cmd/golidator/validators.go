package golidator

import (
	neturl "net/url"
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

func (v *validator) getvalues(i int) (string, string, reflect.Value, error) {
	field := v.value.Type().Field(i).Name
	fieldType := v.value.Type().Field(i).Type.Name()
	value := v.value.Field(i)
	if v.value.Field(i).Kind() == reflect.Ptr {
		if v.value.Field(i).IsNil() {
			message := "Must not be missing from the body"
			return field, "", reflect.Value{}, errors.New(message)
		}
		fieldType = v.value.Field(i).Elem().Type().Name()
		value = v.value.Field(i).Elem()
	}

	return field, fieldType, value, nil
}

func (v *validator) Required() {
	field, _, _, err := v.getvalues(v.fieldIndex)
	if err != nil {
		message := err.Error()
		v.errors[field] = append(v.errors[field], message)
	}
}

func (v *validator) Notblank() {
	field, fieldType, value, err := v.getvalues(v.fieldIndex)
	if err != nil {
		return
	}

	if fieldType != "string" {
		message := "Invalid type. Must be string"
		v.errors[field] = append(v.errors[field], message)
		return
	}

	if value.String() == "" {
		message := "Must not be blank"
		v.errors[field] = append(v.errors[field], message)
		return
	}
}

func (v *validator) ValidEmail(str string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(str)
}

func (v *validator) Email() {
	field, fieldType, value, err := v.getvalues(v.fieldIndex)
	if err != nil {
		return
	}

	if fieldType != "string" {
		message := "Invalid type. Must be string"
		v.errors[field] = append(v.errors[field], message)
		return
	}

	if !v.ValidEmail(value.String()) {
		message := "Must be a valid email"
		v.errors[field] = append(v.errors[field], message)
		return
	}
}

func (v *validator) ValidUrl(str string) bool {
	_, err := neturl.ParseRequestURI(str)
	return err == nil
}


func (v *validator) Url() {
	field, fieldType, value, err := v.getvalues(v.fieldIndex)
	if err != nil {
		return
	}

	if fieldType != "string" {
		message := "Invalid type. Must be string"
		v.errors[field] = append(v.errors[field], message)
		return
	}

	if !v.ValidUrl(value.String()) {
		message := "Must be a valid url"
		v.errors[field] = append(v.errors[field], message)
		return
	}
}

func (v *validator) Min() {
	field, fieldType, value, err := v.getvalues(v.fieldIndex)
	if err != nil {
		return
	}

	if fieldType != "string" {
		message := "Invalid type. Must be string"
		v.errors[field] = append(v.errors[field], message)
		return
	}

	if len(value.String()) < v.fieldLength {
		message := fmt.Sprintf("Must be at least %v characters long", v.fieldLength)
		v.errors[field] = append(v.errors[field], message)
		return
	}
}

func (v *validator) Max() {
	field, fieldType, value, err := v.getvalues(v.fieldIndex)
	if err != nil {
		return
	}

	if fieldType != "string" {
		message := "Invalid type. Must be string"
		v.errors[field] = append(v.errors[field], message)
		return
	}

	if len(value.String()) > v.fieldLength {
		message := fmt.Sprintf("Must be at least %v characters long", v.fieldLength)
		v.errors[field] = append(v.errors[field], message)
	}

	return
}

func (v *validator) Notempty() {
	field := v.value.Type().Field(v.fieldIndex).Name
	if v.value.Field(v.fieldIndex).Type().Kind() != reflect.Slice {
		return
	}

	value := v.value.Field(v.fieldIndex).Len()
	if value == 0 {
		message := "Array must not be empty"
		v.errors[field] = append(v.errors[field], message)
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
		message := "Invalid type. Must be array"
		v.errors[field] = append(v.errors[field], message)
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
