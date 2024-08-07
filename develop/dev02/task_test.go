package main

import (
	"testing"
)

func TestUnpack(t *testing.T) {
	testCases := []struct {
		s    string
		want string
	}{
		{"", ""},                      // empty string
		{"abcd", "abcd"},              // without nums
		{"a4bc2d5e", "aaaabccddddde"}, // one digit nums
		{"a12", "aaaaaaaaaaaa"},       // big num
		// escape последовательности
		{`\02\7`, `007`},       // only digits
		{`qwe\4\5`, `qwe45`},   // without nums
		{`qwe\45`, `qwe44444`}, // one digit nums
		{`\\`, `\`},            // escape escape-rune
		{`qwe\\5`, `qwe\\\\\`}, // escape escape-rune
	}

	for _, tc := range testCases {
		got, err := Unpack(tc.s)
		if err != nil {
			t.Fatalf("Err should be nil, but: %s", err.Error())
		}
		if got != tc.want {
			t.Errorf("want \"%s\" got \"%s\"", tc.want, got)
		}
	}
}

func TestErrorUnpack(t *testing.T) {
	testCases := []struct {
		s      string
		errstr string
	}{
		{"45", "a string cannot start with a number"}, // just num
		{`\`, "incomplete escape sequence"},           // incomplete escape
		{`abc\`, "incomplete escape sequence"},        // incomplete escape

	}

	for _, tc := range testCases {
		got, err := Unpack(tc.s)
		if err == nil {
			t.Fatalf("Err == nil s: %s", tc.s)
		}
		if got != "" {
			t.Fatalf("Got should be \"\" but eq %s; s: %s", got, tc.s)
		}
		if err.Error() != tc.errstr {
			t.Errorf("want \"%s\" got \"%s\"", tc.errstr, err)
		}
	}
}
