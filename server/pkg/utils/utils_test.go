package utils

import (
	"regexp"
	"testing"
)

func TestGetRandomNumbers(t *testing.T) {
	value, err := GetRandomNumbers(32)
	if err != nil {
		t.Fatalf("generate random numbers: %v", err)
	}
	if matched := regexp.MustCompile(`^[0-9]{32}$`).MatchString(value); !matched {
		t.Fatalf("unexpected random number format: %q", value)
	}
}

func TestGetRandomNumbersRejectsInvalidLength(t *testing.T) {
	if _, err := GetRandomNumbers(0); err == nil {
		t.Fatal("expected an error for zero length")
	}
}
