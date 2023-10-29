package golidator

import (
	neturl "net/url"
	"errors"
	"fmt"
	"reflect"
	"regexp"
)

func getvalues(object reflect.Value, i int) (string, string, reflect.Value, error) {
	field := object.Type().Field(i).Name
	fieldType := object.Type().Field(i).Type.Name()
	value := object.Field(i)
	if object.Field(i).Kind() == reflect.Ptr {
		if object.Field(i).IsNil() {
			message := "Must not be missing from the body"
			return field, "", reflect.Value{}, errors.New(message)
		}
		fieldType = object.Field(i).Elem().Type().Name()
		value = object.Field(i).Elem()
	}

	return field, fieldType, value, nil
}

func required(errorsMap map[string][]string, rValue reflect.Value, i int, length int) {
	field, _, _, err := getvalues(rValue, i)
	if err != nil {
		message := err.Error()
		errorsMap[field] = append(errorsMap[field], message)
	}

	return
}

func NotBlank(errorsMap map[string][]string, rValue reflect.Value, i int, length int) {
	field, fieldType, value, err := getvalues(rValue, i)
	if err != nil {
		return
	}

	if fieldType != "string" {
		message := "Invalid type. Must be string"
		errorsMap[field] = append(errorsMap[field], message)
		return
	}

	if value.String() == "" {
		message := "Must not be blank"
		errorsMap[field] = append(errorsMap[field], message)
	}

	return
}

func ValidEmail(str string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(str)
}

func email(errorsMap map[string][]string, rValue reflect.Value, i int, length int) {
	field, fieldType, value, err := getvalues(rValue, i)
	if err != nil {
		return
	}

	if fieldType != "string" {
		message := "Invalid type. Must be string"
		errorsMap[field] = append(errorsMap[field], message)
		return
	}

	if !ValidEmail(value.String()) {
		message := "Must be a valid email"
		errorsMap[field] = append(errorsMap[field], message)
	}

	return
}

func ValidUrl(str string) bool {
	_, err := neturl.ParseRequestURI(str)
	return err == nil
}


func url(errorsMap map[string][]string, rValue reflect.Value, i int, length int) {
	field, fieldType, value, err := getvalues(rValue, i)
	if err != nil {
		return
	}

	if fieldType != "string" {
		message := "Invalid type. Must be string"
		errorsMap[field] = append(errorsMap[field], message)
		return
	}

	if !ValidUrl(value.String()) {
		message := "Must be a valid url"
		errorsMap[field] = append(errorsMap[field], message)
	}

	return
}

func min(errorsMap map[string][]string, rValue reflect.Value, i int, length int) {
	field, fieldType, value, err := getvalues(rValue, i)
	if err != nil {
		return
	}

	if fieldType != "string" {
		message := "Invalid type. Must be string"
		errorsMap[field] = append(errorsMap[field], message)
		return
	}

	if len(value.String()) < length {
		message := fmt.Sprintf("Must be at least %v characters long", length)
		errorsMap[field] = append(errorsMap[field], message)
	}

	return
}

func max(errorsMap map[string][]string, rValue reflect.Value, i int, length int) {
	field, fieldType, value, err := getvalues(rValue, i)
	if err != nil {
		return
	}

	if fieldType != "string" {
		message := "Invalid type. Must be string"
		errorsMap[field] = append(errorsMap[field], message)
		return
	}

	if len(value.String()) > length {
		message := fmt.Sprintf("Must be at least %v characters long", length)
		errorsMap[field] = append(errorsMap[field], message)
	}

	return
}

func notempty(errorsMap map[string][]string, rValue reflect.Value, i int, length int) {
	field := rValue.Type().Field(i).Name
	if rValue.Field(i).Type().Kind() != reflect.Slice {
		return
	}

	value := rValue.Field(i).Len()
	if value == 0 {
		message := "Array must not be empty"
		errorsMap[field] = append(errorsMap[field], message)
	}

	return
}

func valarray(errorsMap map[string][]string, rValue reflect.Value, i int, length int) {
	field := rValue.Type().Field(i).Name
	array := rValue.Field(i)
	if array.Kind() == reflect.Ptr {
		if array.IsNil() {
			return
		}

		array = array.Elem()
	}

	if array.Type().Kind() != reflect.Slice {
		message := "Invalid type. Must be array"
		errorsMap[field] = append(errorsMap[field], message)
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
			errorsMap[summaryField] = summary[j].Errors
		}
	}

	return
}
