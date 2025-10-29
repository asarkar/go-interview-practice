package palindrome

import (
	"unicode"
)

// IsPalindrome checks if a string is a palindrome.
// A palindrome reads the same backward as forward, ignoring case, spaces, and punctuation.
func IsPalindrome(s string) bool {
	// 1. Clean the string (remove spaces, punctuation, and convert to lowercase)
	// 2. Check if the cleaned string is the same forwards and backwards
	runes := make([]rune, 0)
	for _, r := range s {
		if isAlphanum(r) {
			runes = append(runes, r)
		}
	}
	n := len(runes)
	for i := range (n + 1) / 2 {
		if unicode.ToLower(runes[i]) != unicode.ToLower(runes[n-i-1]) {
			return false
		}
	}
	return true
}

func isAlphanum(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
