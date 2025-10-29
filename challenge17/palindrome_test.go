package palindrome

import (
	"strings"
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Empty string", "", true},
		{"Single character", "a", true},
		{"Simple palindrome", "racecar", true},
		{"Simple non-palindrome", "hello", false},
		{"Palindrome with spaces", "never odd or even", true},
		{"Palindrome with mixed case", "RaceCar", true},
		{"Palindrome with punctuation", "A man, a plan, a canal: Panama", true},
		{"Palindrome with numbers", "12321", true},
		{"Non-palindrome with numbers", "12345", false},
		{"Complex palindrome", "Madam, I'm Adam", true},
		{"Alphanumeric palindrome", "A1b2c3c2b1A", true},
		{"Alphanumeric non-palindrome", "A1b2c3d4e5", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPalindrome(tt.input)
			if result != tt.expected {
				t.Errorf("IsPalindrome(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

// TestIsPalindromeEdgeCases tests edge cases for the IsPalindrome function
func TestIsPalindromeEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			"Special characters only",
			"!@#$%^&*()",
			true,
		}, // Special chars should be ignored, leaving empty string
		{"Mixed with emoji", "😊1221😊", true}, // Should work with Unicode characters
		{"Very long palindrome", "a" + strings.Repeat("b", 10000) + "a", true},
		{"All spaces", "          ", true}, // All spaces should be ignored, leaving empty string
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPalindrome(tt.input)
			if result != tt.expected {
				t.Errorf("IsPalindrome(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}
