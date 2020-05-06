package main

import (
	"testing"
)



func TestNormalize(t *testing.T) {
	testCases := []struct {
		input string
		expected string
	}{
		{"1234567890","1234567890"},
		{"123 456 7891","1234567891"},
		{"(123) 456 7892","1234567892"},
		{"(123) 456-7893","1234567893"},
		{"123-456-7894","1234567894"},
		{"123-456-7890","1234567890"},
		{"1234567892","1234567892"},
		{"(123)456-7892","1234567892"},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			result := normalize(tc.input)
			if result != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, result)
			} else {
				t.Logf("Yay it worked! Got %s, wanted %s", result, tc.expected)
			}
		})
	}
}