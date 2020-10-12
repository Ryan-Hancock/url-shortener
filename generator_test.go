package main

import (
	"math/rand"
	"strings"
	"testing"
)

func TestGenerateString(t *testing.T) {
	t.Run("should generate the correct length", func(t *testing.T) {
		length := rand.Intn(10)
		a := GenerateString(length, "abcde")

		if len(a) != length {
			t.Errorf("GenerateString() returned the wrong length %d instead of %d", len(a), length)
		}
	})

	t.Run("should contain only letters in the charset", func(t *testing.T) {
		a := GenerateString(3, "a")

		if !strings.Contains(a, "a") {
			t.Errorf("GenerateString() does not read from supplied charset")
		}
	})
}
