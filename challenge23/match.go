package match

import (
	perf "go-interview-practice/challenge16"
)

// NaivePatternMatch performs a brute force search for pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func NaivePatternMatch(text, pattern string) []int {
	n, m := len(text), len(pattern)
	res := make([]int, 0)
	if n < m || m == 0 {
		return res
	}

	for i := range n - m + 1 {
		s := text[i : i+m]
		if s == pattern {
			res = append(res, i)
		}
	}

	return res
}

// KMPSearch implements the Knuth-Morris-Pratt algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func KMPSearch(text, pattern string) []int {
	haystack := []rune(text)
	needle := []rune(pattern)
	n := len(needle)
	if len(haystack) < n || n == 0 {
		return make([]int, 0)
	}
	lps := make([]int, n)

	// Build LPS.
	perf.KMPLoop(needle, needle, 1, lps, true)

	// Search.
	return perf.KMPLoop(haystack, needle, 0, lps, false)
}

// RabinKarpSearch implements the Rabin-Karp algorithm to find pattern in text.
// Returns a slice of all starting indices where the pattern is found.
func RabinKarpSearch(text, pattern string) []int {
	// Hash Formula:
	// H = (c0 x b^(n-1) + c1 x b^(n-2) + ... + cn-1 x b^0) mod p
	// 	where, `ck` is ASCII code of the corresponding character, `b` is the size of the alphabet,
	// 	`n` is the length of the text, and `p` is a large prime number.
	//
	// Hash Update Formula:
	// (b x (H-old - c-old x b^(n-1)) + c-new) mod p
	const (
		base  int64 = 256
		prime int64 = 101
	)

	res := make([]int, 0)

	n, m := len(text), len(pattern)
	if n < m || m == 0 {
		return res
	}

	h := powMod(base, int64(m-1), prime)

	var patHash, txtHash int64
	for i := range m {
		patHash = (base*patHash + int64(pattern[i])) % prime
		txtHash = (base*txtHash + int64(text[i])) % prime
	}

	for i := range n - m + 1 {
		if patHash == txtHash && text[i:i+m] == pattern {
			res = append(res, i)
		}
		if i < n-m {
			txtHash = (base*(txtHash-int64(text[i])*h) + int64(text[i+m])) % prime
			if txtHash < 0 {
				txtHash += prime
			}
		}
	}
	return res
}

// powMod computes (base^exp) % mod efficiently using binary exponentiation.
// Time complexity: O(log exp)
// No overflow (as long as base, mod < 2^63)
func powMod(base, exp, mod int64) int64 {
	result := int64(1)
	base %= mod

	for exp > 0 {
		if exp&1 == 1 {
			result = (result * base) % mod
		}
		base = (base * base) % mod
		exp >>= 1 // divide by 2
	}
	return result
}
