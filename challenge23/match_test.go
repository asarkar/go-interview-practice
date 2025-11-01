package match

import (
	"reflect"
	"sort"
	"testing"
)

func TestNaivePatternMatch(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		pattern  string
		expected []int
	}{
		{"Basic case", "ABABDABACDABABCABAB", "ABABCABAB", []int{10}},
		{"Multiple occurrences", "AABAACAADAABAABA", "AABA", []int{0, 9, 12}},
		{"Occurrence at the beginning", "GEEKSFORGEEKS", "GEEK", []int{0, 8}},
		{"Overlapping occurrences", "AAAAAA", "AA", []int{0, 1, 2, 3, 4}},
		{"No occurrences", "ABCDEFG", "XYZ", []int{}},
		{"Empty pattern", "ABCDEFG", "", []int{}},
		{"Empty text", "", "ABC", []int{}},
		{"Both empty", "", "", []int{}},
		{"Pattern longer than text", "ABC", "ABCDEF", []int{}},
		{"Pattern is the entire text", "ABCDEF", "ABCDEF", []int{0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NaivePatternMatch(tt.text, tt.pattern)
			sort.Ints(result) // Sort to ensure consistent order for comparison
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("NaivePatternMatch(%s, %s) = %v, expected %v",
					tt.text, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestKMPSearch(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		pattern  string
		expected []int
	}{
		{"Basic case", "ABABDABACDABABCABAB", "ABABCABAB", []int{10}},
		{"Multiple occurrences", "AABAACAADAABAABA", "AABA", []int{0, 9, 12}},
		{"Occurrence at the beginning", "GEEKSFORGEEKS", "GEEK", []int{0, 8}},
		{"Overlapping occurrences", "AAAAAA", "AA", []int{0, 1, 2, 3, 4}},
		{"No occurrences", "ABCDEFG", "XYZ", []int{}},
		{"Empty pattern", "ABCDEFG", "", []int{}},
		{"Empty text", "", "ABC", []int{}},
		{"Both empty", "", "", []int{}},
		{"Pattern longer than text", "ABC", "ABCDEF", []int{}},
		{"Pattern is the entire text", "ABCDEF", "ABCDEF", []int{0}},
		{"Complex pattern", "ACACACACGTACACACA", "ACACACA", []int{0, 10}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := KMPSearch(tt.text, tt.pattern)
			sort.Ints(result) // Sort to ensure consistent order for comparison
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("KMPSearch(%s, %s) = %v, expected %v",
					tt.text, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestRabinKarpSearch(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		pattern  string
		expected []int
	}{
		{"Basic case", "ABABDABACDABABCABAB", "ABABCABAB", []int{10}},
		{"Multiple occurrences", "AABAACAADAABAABA", "AABA", []int{0, 9, 12}},
		{"Occurrence at the beginning", "GEEKSFORGEEKS", "GEEK", []int{0, 8}},
		{"Overlapping occurrences", "AAAAAA", "AA", []int{0, 1, 2, 3, 4}},
		{"No occurrences", "ABCDEFG", "XYZ", []int{}},
		{"Empty pattern", "ABCDEFG", "", []int{}},
		{"Empty text", "", "ABC", []int{}},
		{"Both empty", "", "", []int{}},
		{"Pattern longer than text", "ABC", "ABCDEF", []int{}},
		{"Pattern is the entire text", "ABCDEF", "ABCDEF", []int{0}},
		{
			"Large text",
			"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla facilisi.",
			"sit",
			[]int{18},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RabinKarpSearch(tt.text, tt.pattern)
			sort.Ints(result) // Sort to ensure consistent order for comparison
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("RabinKarpSearch(%s, %s) = %v, expected %v",
					tt.text, tt.pattern, result, tt.expected)
			}
		})
	}
}

func TestExampleCases(t *testing.T) {
	// Example 1
	text1 := "ABABDABACDABABCABAB"
	pattern1 := "ABABCABAB"
	expected1 := []int{10}

	result1 := NaivePatternMatch(text1, pattern1)
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Example 1: NaivePatternMatch(%s, %s) = %v, expected %v",
			text1, pattern1, result1, expected1)
	}

	// Example 2
	text2 := "AABAACAADAABAABA"
	pattern2 := "AABA"
	expected2 := []int{0, 9, 12}

	result2 := KMPSearch(text2, pattern2)
	sort.Ints(result2) // Sort to ensure consistent order for comparison
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Example 2: KMPSearch(%s, %s) = %v, expected %v",
			text2, pattern2, result2, expected2)
	}

	// Example 3
	text3 := "GEEKSFORGEEKS"
	pattern3 := "GEEK"
	expected3 := []int{0, 8}

	result3 := RabinKarpSearch(text3, pattern3)
	sort.Ints(result3) // Sort to ensure consistent order for comparison
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Example 3: RabinKarpSearch(%s, %s) = %v, expected %v",
			text3, pattern3, result3, expected3)
	}

	// Example 4
	text4 := "AAAAAA"
	pattern4 := "AA"
	expected4 := []int{0, 1, 2, 3, 4}

	result4 := NaivePatternMatch(text4, pattern4)
	sort.Ints(result4) // Sort to ensure consistent order for comparison
	if !reflect.DeepEqual(result4, expected4) {
		t.Errorf("Example 4: NaivePatternMatch(%s, %s) = %v, expected %v",
			text4, pattern4, result4, expected4)
	}
}

// Benchmark to compare the performance of different algorithms
func BenchmarkPatternMatching(b *testing.B) {
	text := "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nulla facilisi. Sed euismod, nisl eget ultricies aliquam, nisl nisl ultricies nisl, eget ultricies nisl eget ultricies aliquam, nisl nisl ultricies nisl, eget ultricies nisl eget ultricies aliquam, nisl nisl ultricies nisl, eget ultricies nisl eget ultricies aliquam, nisl nisl ultricies nisl, eget ultricies nisl eget ultricies aliquam."
	pattern := "ultricies"

	b.Run("NaivePatternMatch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			NaivePatternMatch(text, pattern)
		}
	})

	b.Run("KMPSearch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			KMPSearch(text, pattern)
		}
	})

	b.Run("RabinKarpSearch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			RabinKarpSearch(text, pattern)
		}
	})
}
