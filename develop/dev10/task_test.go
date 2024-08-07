package main

import (
	"bytes"
	"testing"
)

func TestEcho(t *testing.T) {
	has := "echo\nabc"
	want := has

	got := bytes.NewBuffer([]byte{})
	err := echo(bytes.NewBufferString(has), got)

	if err != nil {
		t.Fatal("err should be nil", err)
	}
	if got.String() != want {
		t.Errorf("got: %s, want: %s", got.String(), want)
	}
}
