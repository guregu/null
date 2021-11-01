package null

import (
	"testing"
)

func BenchmarkIntUnmarshalJSON(b *testing.B) {
	input := []byte("123456")
	var nullable Int
	for n := 0; n < b.N; n++ {
		nullable.UnmarshalJSON(input)
	}
}

func BenchmarkIntStringUnmarshalJSON(b *testing.B) {
	input := []byte(`"123456"`)
	var nullable String
	for n := 0; n < b.N; n++ {
		nullable.UnmarshalJSON(input)
	}
}

func BenchmarkNullIntUnmarshalJSON(b *testing.B) {
	input := []byte("null")
	var nullable Int
	for n := 0; n < b.N; n++ {
		nullable.UnmarshalJSON(input)
	}
}

func BenchmarkStringUnmarshalJSON(b *testing.B) {
	input := []byte(`"hello"`)
	var nullable String
	for n := 0; n < b.N; n++ {
		nullable.UnmarshalJSON(input)
	}
}

func BenchmarkNullStringUnmarshalJSON(b *testing.B) {
	input := []byte("null")
	var nullable String
	for n := 0; n < b.N; n++ {
		nullable.UnmarshalJSON(input)
	}
}
