package fieldinfo

import (
	"reflect"
	"strconv"
	"strings"
)

const (
	JsonTag     = "json"
	ValidateTag = "validate"
)

func ExtractInfo(structValue reflect.Value, fieldIndex int) Info {
	structType := structValue.Type()
	field := structType.Field(fieldIndex)
	fieldValue := structValue.Field(fieldIndex)

	fieldType := field.Type
	fieldKind := fieldType.Kind()
	originalKind := fieldKind
	isPointer := false

	if fieldKind == reflect.Pointer {
		isPointer = true
		fieldType = fieldType.Elem()
		fieldKind = fieldType.Kind()
	}

	jsonName := field.Name
	if jsonTag := field.Tag.Get(JsonTag); jsonTag != "" {
		parts := strings.Split(jsonTag, ",")
		if len(parts) > 0 && parts[0] != "" {
			jsonName = parts[0]
		}
	}

	validateTag := field.Tag.Get(ValidateTag)
	validatorArgs, validatorInts, isRequired := parseValidatorArgs(validateTag)

	return Info{
		Index:         fieldIndex,
		Name:          field.Name,
		JSONName:      jsonName,
		Type:          fieldType,
		Kind:          fieldKind,
		TypeName:      fieldType.Name(),
		ValidateTag:   validateTag,
		IsPointer:     isPointer,
		OriginalKind:  originalKind,
		Value:         fieldValue,
		ValidatorStrs: validatorArgs,
		ValidatorInts: validatorInts,
		IsRequired:    isRequired,
	}
}

func parseValidatorArgs(validateTag string) (map[string]string, map[string]int, bool) {
	args := make(map[string]string)
	ints := make(map[string]int)
	isRequired := false
	if validateTag == "" {
		return args, ints, isRequired
	}

	validators := strings.SplitSeq(validateTag, ",")
	for validator := range validators {
		validator = strings.TrimSpace(validator)
		if validator == "" {
			continue
		}

		validatorName := validator
		if idx := strings.IndexByte(validator, '='); idx != -1 {
			validatorName = validator[:idx]
			value := validator[idx+1:]
			if value != "" {
				args[validatorName] = value

				if intValue, err := strconv.Atoi(value); err == nil && intValue >= 0 {
					ints[validatorName] = intValue
				}
			}
		}

		if validatorName == "required" {
			isRequired = true
		}
	}

	return args, ints, isRequired
}
