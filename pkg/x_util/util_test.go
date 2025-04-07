package x_util

import (
	"testing"
	"unicode"
)

// TestRandomStringLength verifies that the generated string has the correct length.
func TestRandomStringLength(t *testing.T) {
	const length = 16

	str, err := RandomString(length)
	if err != nil {
		t.Fatalf("RandomString returned error: %v", err)
	}

	if len(str) != length {
		t.Errorf("Expected string of length %d, got %d", length, len(str))
	}
}

// TestRandomStringCharset ensures all characters are alphanumeric (A-Z, a-z, 0-9).
func TestRandomStringCharset(t *testing.T) {
	str, err := RandomString(100)
	if err != nil {
		t.Fatalf("RandomString returned error: %v", err)
	}

	for _, r := range str {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) {
			t.Errorf("Invalid character found in string: %q", r)
		}
	}
}

// TestRandomStringZeroLength checks behavior when length is zero.
func TestRandomStringZeroLength(t *testing.T) {
	str, err := RandomString(0)
	if err != nil {
		t.Fatalf("RandomString returned error for zero length: %v", err)
	}

	if str != "" {
		t.Errorf("Expected empty string, got %q", str)
	}
}
