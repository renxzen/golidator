package validators

import (
	"github.com/renxzen/golidator/internal/fieldinfo"
)

type ValidatorFunc func(fieldInfo fieldinfo.Info) string

var Registry = map[string]ValidatorFunc{
	"notblank": NotBlank,
	"email":    Email,
	"numeric":  Numeric,
	"url":      URL,
	"required": Required,
	"notempty": NotEmpty,
	"min":      Min,
	"max":      Max,
	"len":      Len,
	"isarray":  IsArray,
}
