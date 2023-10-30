package util

import (
	"unicode"
)

func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	firstRune := []rune(s)[0]
	return string(unicode.ToUpper(firstRune)) + s[1:]
}
