package util

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"
)

/*
	The first implementation, which uses loops and buffers, is generally
	faster than the second one, which relies on regular expressions.
	There are a few reasons for this:
	- Overhead: Regular expressions have a parsing and compilation overhead.
	For simple string manipulation tasks like this, a straightforward
	loop can often outperform regex due to its directness.
	- Memory allocation: The first implementation uses a bytes.Buffer, which
	is efficient for string concatenation. Regular expressions might involve
	more internal memory allocations and copying, especially when using
	methods like ReplaceAllString.
	- Complexity: Regular expressions are powerful and can match complex
	patterns, but with that power comes a bit of overhead. The pattern we
	are usually matching and replacing here is relatively simple, so a
	loop can handle it very efficiently.
*/

func ToSnakeCase(s string) string {
	var result bytes.Buffer
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteByte('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	return result.String()
}

func ToSnakeCaseRegex(str string) string {
	allCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	newStr := allCap.ReplaceAllString(str, "${1}_${2}")
	return strings.ToLower(newStr)
}
