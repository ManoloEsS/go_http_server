package handlers

import "testing"

func TestProfaneFilter(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{"kerfuffle", "****"},
		{"fornax", "****"},
		{"sharbert", "****"},
		{"Kerfuffle", "****"},
		{"KERFUFFLE", "****"},
		{"FORNAX", "****"},
		{"SHARBERT", "****"},
		{"FoRnAx", "****"},
		{"word fornax word", "word **** word"},
		{"fornax kerfuffle sharbert", "**** **** ****"},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.input, func(t *testing.T) {
			t.Parallel()
			actualOutput := profaneFilter(testCase.input, badWords)
			if actualOutput != testCase.expected {
				t.Errorf("expected profaneFilter(%s) to be %s but got %s", testCase.input, testCase.expected, actualOutput)
			}
		})
	}
}
