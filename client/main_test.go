package main

import (
	"testing"
)

func TestPalindrome(t *testing.T) {
	tests := []struct {
		Name     string
		Word     string
		Expected bool
	}{
		{Name: "ASCII Palindrome", Word: "testset", Expected: true},
		{Name: "ASCII Non-Palindrome", Word: "testsetx", Expected: false},
		{Name: "Unicode Palindrome", Word: "壹壱漢字漢壱壹", Expected: true},
		{Name: "Unicode Non-Palindrome", Word: "x壹壱漢字漢壱壹y", Expected: false},
		{Name: "Empty String", Word: "", Expected: true},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			result := IsPalindrome(test.Word)
			if result != test.Expected {
				t.Errorf("IsPanindrome(%q) returned %v; expected %v", test.Word, result, test.Expected)
			}
		})
	}
}
