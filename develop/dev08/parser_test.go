package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestDefaultParser(t *testing.T) {
	testCases := []struct {
		has  []byte
		want Entity
	}{
		{
			has: []byte("cmd 1 2 3"),
			want: Entity{
				Cmds: [][]string{
					{"cmd", "1", "2", "3"},
				},
			},
		},
		{
			has: []byte("cmd1 a b c | cmd2 a b c | cmd3"),
			want: Entity{
				Cmds: [][]string{
					{"cmd1", "a", "b", "c"},
					{"cmd2", "a", "b", "c"},
					{"cmd3"},
				},
			},
		},
		{
			has: []byte("cmd1 a b c | cmd2 a b c | cmd3 &"),
			want: Entity{
				Cmds: [][]string{
					{"cmd1", "a", "b", "c"},
					{"cmd2", "a", "b", "c"},
					{"cmd3"},
				},
				Bg: true,
			},
		},
		{
			has: []byte("cmd1 a b c &"),
			want: Entity{
				Cmds: [][]string{
					{"cmd1", "a", "b", "c"},
				},
				Bg: true,
			},
		},
		{
			has: []byte(""),
			want: Entity{
				EOF: true,
			},
		},
	}

	for _, tc := range testCases {
		p := NewDefaultParser(bytes.NewBuffer(tc.has))
		got := p.Parse()

		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("\ngot:  %v\nwant: %v\n", got, tc.want)
		}
	}
}
