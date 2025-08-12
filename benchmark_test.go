package golidator

import (
	"testing"
)

type BenchmarkStruct struct {
	Field1  string   `json:"field1" validate:"notblank,min=3,max=50"`
	Field2  *string  `json:"field2" validate:"required,email"`
	Field3  int      `json:"field3" validate:"min=1,max=100"`
	Field4  []string `json:"field4" validate:"notempty,len=3"`
	Field5  *int     `json:"field5" validate:"required,min=0"`
	Field6  string   `json:"field6" validate:"url"`
	Field7  string   `json:"field7" validate:"numeric"`
	Field8  []int    `json:"field8" validate:"isarray,notempty"`
	Field9  string   `json:"field9" validate:"len=10"`
	Field10 *string  `json:"field10" validate:"notblank"`
}

// Benchmark complex struct validation - Cache vs No Cache
func BenchmarkComplexStruct(b *testing.B) {
	email := "test@example.com"
	five := 5
	url := "https://example.com"
	numeric := "12345"
	tenChars := "1234567890"
	notBlank := "value"

	testStruct := BenchmarkStruct{
		Field1:  "valid",
		Field2:  &email,
		Field3:  50,
		Field4:  []string{"a", "b", "c"},
		Field5:  &five,
		Field6:  url,
		Field7:  numeric,
		Field8:  []int{1, 2, 3},
		Field9:  tenChars,
		Field10: &notBlank,
	}

	b.Run("WithCache", func(b *testing.B) {
		SetCaching(true)
		b.ResetTimer()
		for b.Loop() {
			if _, err := Validate(testStruct); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Benchmark validation with errors - Cache vs No Cache
func BenchmarkValidationWithErrors(b *testing.B) {
	invalidEmail := "invalid-email"
	zero := 0
	invalidURL := "not-a-url"
	nonNumeric := "abc123"
	wrongLength := "short"
	empty := ""

	testStruct := BenchmarkStruct{
		Field1:  "",            // notblank error
		Field2:  &invalidEmail, // email error
		Field3:  200,           // max error
		Field4:  []string{},    // notempty error
		Field5:  &zero,         // min error (assuming min=1)
		Field6:  invalidURL,    // url error
		Field7:  nonNumeric,    // numeric error
		Field8:  []int{},       // notempty error
		Field9:  wrongLength,   // len error
		Field10: &empty,        // notblank error
	}

	b.Run("WithCache", func(b *testing.B) {
		SetCaching(true)
		b.ResetTimer()
		for b.Loop() {
			_, err := Validate(testStruct)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("WithoutCache", func(b *testing.B) {
		SetCaching(false)
		defer func() { SetCaching(true) }()
		b.ResetTimer()
		for b.Loop() {
			if _, err := Validate(testStruct); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Benchmark simple struct validation - Cache vs No Cache
func BenchmarkSimpleStruct(b *testing.B) {
	type SimpleStruct struct {
		Name  string `json:"name" validate:"notblank"`
		Email string `json:"email" validate:"email"`
	}

	simple := SimpleStruct{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	b.Run("WithCache", func(b *testing.B) {
		SetCaching(true)
		b.ResetTimer()
		for b.Loop() {
			if _, err := Validate(simple); err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("WithoutCache", func(b *testing.B) {
		SetCaching(false)
		defer func() { SetCaching(true) }()
		b.ResetTimer()
		for b.Loop() {
			if _, err := Validate(simple); err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Benchmark multiple struct types alternating - Cache vs No Cache
func BenchmarkMultipleTypes(b *testing.B) {
	type SimpleStruct struct {
		Name  string `json:"name" validate:"notblank"`
		Email string `json:"email" validate:"email"`
	}

	type ComplexStruct struct {
		ID       int      `json:"id" validate:"min=1"`
		Tags     []string `json:"tags" validate:"notempty"`
		Optional *string  `json:"optional" validate:"notblank"`
	}

	simple := SimpleStruct{
		Name:  "John",
		Email: "john@example.com",
	}

	optional := "value"
	complexStruct := ComplexStruct{
		ID:       123,
		Tags:     []string{"tag1", "tag2"},
		Optional: &optional,
	}

	b.Run("WithCache", func(b *testing.B) {
		SetCaching(true)
		b.ResetTimer()
		for i := 0; b.Loop(); i++ {
			if i%2 == 0 {
				if _, err := Validate(simple); err != nil {
					b.Fatal(err)
				}
			} else {
				if _, err := Validate(complexStruct); err != nil {
					b.Fatal(err)
				}
			}
		}
	})

	b.Run("WithoutCache", func(b *testing.B) {
		SetCaching(false)
		defer func() { SetCaching(true) }()
		b.ResetTimer()
		for i := 0; b.Loop(); i++ {
			if i%2 == 0 {
				if _, err := Validate(simple); err != nil {
					b.Fatal(err)
				}
			} else {
				if _, err := Validate(complexStruct); err != nil {
					b.Fatal(err)
				}
			}
		}
	})
}
