package main

import (
	"reflect"
	"testing"
)

func TestAnagramGroup(t *testing.T) {
	testCases := []struct {
		words []string
		want  map[string][]string
	}{
		{ // empty
			words: []string{},
			want:  map[string][]string{},
		},
		{ // one word groups
			words: []string{"а", "б", "в", "д", "дд"},
			want:  map[string][]string{},
		},
		{ // many words groups (first word of goup in dict has be key)
			words: []string{"тяпка", "слиток", "столик", "пятак", "пятка", "листок", "одинокий"},
			want: map[string][]string{
				"пятак":  {"пятак", "пятка", "тяпка"},
				"листок": {"листок", "слиток", "столик"},
			},
		},
		{ // lowercase
			words: []string{"ПяТаК", "слиТок", "сТОЛик", "тЯпКА"},
			want: map[string][]string{
				"пятак":  {"пятак", "тяпка"},
				"слиток": {"слиток", "столик"},
			},
		},
		{ // unique
			words: []string{"ПяТаК", "слиТок", "сТОЛик", "тЯпКА", "ПЯТАК", "СтоЛИк"},
			want: map[string][]string{
				"пятак":  {"пятак", "тяпка"},
				"слиток": {"слиток", "столик"},
			},
		},
		{ // unique & one word group
			words: []string{"П", "п"},
			want:  map[string][]string{},
		},
	}

	for _, tc := range testCases {
		got := FindAnagrams(tc.words)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("\ngot:  %v\nwant: %v\n", got, tc.want)
		}
	}
}
