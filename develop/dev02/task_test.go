package main

import (
	"testing"
)

func TestMain_unpacking(t *testing.T) {
	testTable := []struct {
		name      string
		input     []rune
		expectRes string
	}{
		{
			name:      "OK1",
			input:     []rune("a4bc2d5e"),
			expectRes: "aaaabccddddde",
		},
		{
			name:      "OK2",
			input:     []rune("abcd"),
			expectRes: "abcd",
		},
		{
			name:      "Incorrect",
			input:     []rune("45"),
			expectRes: "",
		},
		{
			name:      "OK3",
			input:     []rune(""),
			expectRes: "",
		},
		{
			name:      "OK4",
			input:     []rune("qwe\\4\\5"),
			expectRes: "qwe45",
		},
		{
			name:      "OK5",
			input:     []rune("qwe\\45"),
			expectRes: "qwe44444",
		},
		{
			name:      "OK6",
			input:     []rune("qwe\\\\5"),
			expectRes: "qwe\\\\\\\\\\",
		},
		{
			name:      "OK7",
			input:     []rune("abc4e\\22d3\\\\\\2d3\\\\w"),
			expectRes: "abcccce22ddd\\2ddd\\w",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := unpacking(testCase.input)
			if err != nil && err.Error() != "некорректная строка" {
				t.Errorf("некорректная строка")
			}

			if res != testCase.expectRes {
				t.Errorf("Ожидалось: %s\nПолучилось: %s\n", testCase.expectRes, res)
			}
		})
	}
}
