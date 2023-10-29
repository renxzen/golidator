package util

import (
	"testing"
)

func TestUsage(t *testing.T) {
	input := "StringToTest"
	output := "string_to_test"

	result := ToSnakeCase(input)

	if output != result {
		t.Fatalf(`ToSnakeCase("StringToTest") = %s, want "%s"`, result, output)
	}
}
