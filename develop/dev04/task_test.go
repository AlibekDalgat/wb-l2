package main

import (
	"reflect"
	"testing"
)

func TestMain_searchAnagramSets(t *testing.T) {
	testTable := []struct {
		name   string
		input  *[]string
		expRes *map[string][]string
	}{
		{
			name:   "OK1",
			input:  &[]string{"пятак", "слиток", "тяпка", "столик", "пятка", "листок"},
			expRes: &map[string][]string{"листок": {"листок", "слиток", "столик"}, "пятак": {"пятак", "пятка", "тяпка"}},
		},
		{
			name:   "OK2",
			input:  &[]string{"пятак", "слиток", "тяпка", "столик", "стол", "пятка", "листок", "лост", "лист"},
			expRes: &map[string][]string{"листок": {"листок", "слиток", "столик"}, "пятак": {"пятак", "пятка", "тяпка"}, "лост": {"лост", "стол"}},
		},
		{
			name:   "OK1",
			input:  &[]string{"пятак", "сЛитоК", "ТЯПКА", "столик", "пяткааа", "листок"},
			expRes: &map[string][]string{"листок": {"листок", "слиток", "столик"}, "пятак": {"пятак", "тяпка"}},
		},
	}

	for _, testCase := range testTable {
		res := searchAnagramSets(testCase.input)
		// исопльзование reflect.DeepEqua гарантирует сравнение глубоко равных значений на которое указывает указатель
		if !reflect.DeepEqual(res, testCase.expRes) {
			t.Errorf("\nОжидалось:\n%v\nПолучилось:\n%v\n", testCase.expRes, res)
		}
	}
}
