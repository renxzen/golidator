package validators

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/renxzen/golidator/internal/fieldinfo"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

func Required(field fieldinfo.Info) string {
	if field.IsNil() {
		return MessageMissing
	}

	return ""
}

func NotBlank(field fieldinfo.Info) string {
	if !field.IsRequired && field.IsNil() {
		return ""
	}

	if !field.IsString() {
		return MessageNotStringType
	}

	if field.String() == "" {
		return MessageNotBlank
	}

	return ""
}

func Email(field fieldinfo.Info) string {
	if !field.IsRequired && field.IsNil() {
		return ""
	}

	if !field.IsString() {
		return MessageNotStringType
	}

	if !emailRegex.MatchString(field.String()) {
		return MessageInvalidEmail
	}

	return ""
}

func URL(field fieldinfo.Info) string {
	if !field.IsRequired && field.IsNil() {
		return ""
	}

	if !field.IsString() {
		return MessageNotStringType
	}

	if _, err := url.ParseRequestURI(field.String()); err != nil {
		return MessageInvalidURL
	}

	return ""
}

func Min(field fieldinfo.Info) string {
	if !field.IsRequired && field.IsNil() {
		return ""
	}

	if field.IsNil() {
		return ""
	}

	minValue, exists := field.GetArgumentInt("min")
	if !exists {
		return ""
	}

	if field.IsString() {
		if field.Len() < minValue {
			return fmt.Sprintf(MessageStrInvalidMin, minValue)
		}
		return ""
	}

	if field.IsInt() {
		if field.Int() < int64(minValue) {
			return fmt.Sprintf(MessageStrInvalidInt, minValue)
		}
		return ""
	}

	if field.IsFloat() {
		if field.Float() < float64(minValue) {
			return fmt.Sprintf(MessageStrInvalidInt, minValue)
		}
		return ""
	}

	return MessageNotStrIntType
}

func Max(field fieldinfo.Info) string {
	if !field.IsRequired && field.IsNil() {
		return ""
	}

	maxValue, exists := field.GetArgumentInt("max")
	if !exists {
		return ""
	}

	if field.IsString() {
		if field.Len() > maxValue {
			return fmt.Sprintf(MessageStrInvalidMax, maxValue)
		}
		return ""
	}

	if field.IsInt() {
		if field.Int() > int64(maxValue) {
			return fmt.Sprintf(MessageIntInvalidMax, maxValue)
		}
		return ""
	}

	if field.IsFloat() {
		if field.Float() > float64(maxValue) {
			return fmt.Sprintf(MessageIntInvalidMax, maxValue)
		}
		return ""
	}

	return MessageNotStrIntType
}

func NotEmpty(field fieldinfo.Info) string {
	if !field.IsSlice() {
		return MessageNotArrayType
	}

	if field.Len() == 0 {
		return MessageEmptyArray
	}

	return ""
}

func IsArray(field fieldinfo.Info) string {
	if !field.IsSlice() {
		return MessageNotArrayType
	}

	return ""
}

func Len(field fieldinfo.Info) string {
	if field.IsNil() {
		return ""
	}

	fieldLength, exists := field.GetArgumentInt("len")
	if !exists {
		return ""
	}

	if !field.IsString() && !field.IsSlice() {
		return MessageNotStrSliceType
	}

	if field.IsString() && field.Len() != fieldLength {
		return fmt.Sprintf(MessageInvalidLength, fieldLength)
	}

	if field.IsSlice() && field.Len() != fieldLength {
		return fmt.Sprintf(MessageInvalidLengthSlice, fieldLength)
	}

	return ""
}

func Numeric(field fieldinfo.Info) string {
	if !field.IsRequired && field.IsNil() {
		return ""
	}

	if !field.IsString() {
		return MessageNotStringType
	}

	if field.Len() == 0 {
		return MessageNotNumeric
	}

	for _, r := range field.String() {
		if r < '0' || r > '9' {
			return MessageNotNumeric
		}
	}

	return ""
}
