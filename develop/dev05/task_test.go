package exercise_5

import (
	"reflect"
	"regexp"
	"testing"
)

func TestParseArgsNoPattern(t *testing.T) {
	var args []string
	re, file, err := parseArgs(args)
	if err == nil {
		t.Logf("parseArgs(args) = %v, %q, %v, expected nil, \"\", error", *re, file, err)
	}
}

func TestParseArgsNoFilename(t *testing.T) {
	args := []string{"a"}
	_, file, _ := parseArgs(args)
	if file != "" {
		t.Logf("filename = %s, expected \"\"", file)
		t.Fail()
	}
}

func TestParseArgsCorrectPattern(t *testing.T) {
	args := []string{"a"}
	re, _, _ := parseArgs(args)

	expected := "a"

	if re.String() != expected {
		t.Logf("regexp = %q, expected %q", re.String(), expected)
		t.Fail()
	}
}

func TestParseArgsCorrectRawPattern(t *testing.T) {
	fixed = true
	defer func() { fixed = false }()
	args := []string{"a"}
	re, _, _ := parseArgs(args)

	expected := `\Qa\E`

	if re.String() != expected {
		t.Logf("regexp = %q, expected %q", re.String(), expected)
		t.Fail()
	}
}

func TestParseArgsIncorrectPattern(t *testing.T) {
	args := []string{`\`}
	re, file, err := parseArgs(args)
	if err == nil {
		t.Logf("parseArgs(args) = %v, %q, %v, expected nil, \"\", error", *re, file, err)
		t.Fail()
	}
}

func TestFilter(t *testing.T) {
	lines := [][]byte{[]byte("a"), []byte("b")}
	re := regexp.MustCompile("a")

	s := filter(lines, re)

	expected := []int{0}

	if !reflect.DeepEqual(s, expected) {
		t.Logf("filter result: %v, expected: %v", s, expected)
		t.Fail()
	}
}

func TestFilterInverse(t *testing.T) {
	invert = true
	defer func() { invert = false }()

	lines := [][]byte{[]byte("a"), []byte("b")}
	re := regexp.MustCompile("a")

	s := filter(lines, re)

	expected := []int{1}

	if !reflect.DeepEqual(s, expected) {
		t.Logf("filter result: %v, expected: %v", s, expected)
		t.Fail()
	}
}
