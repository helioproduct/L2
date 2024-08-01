package main

import (
	"bytes"
	"reflect"
	"testing"
)

func TestSplitByteSep(t *testing.T) {
	sep := byte('\t')
	has := []byte("col1\tcol2\tcol3")
	want := [][]byte{
		[]byte("col1"),
		[]byte("col2"),
		[]byte("col3"),
	}

	splitter := SplitByteSep{sep: sep}
	got := splitter.Split(has)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:  %v\nwant: %v", got, want)
	}
}

func TestFilterSelectedCols(t *testing.T) {
	testCases := []struct {
		selectCols string
		has        [][]byte
		want       [][]byte
	}{
		{
			selectCols: "", // Не выбирать колонок
			has: [][]byte{
				[]byte("col1"),
				[]byte("col2"),
				[]byte("col3"),
			},
			want: nil,
		},
		{
			selectCols: "2", // Выбрать только вторую колонку
			has: [][]byte{
				[]byte("col1"),
				[]byte("col2"),
				[]byte("col3"),
			},
			want: [][]byte{
				[]byte("col2"),
			},
		},
		{
			selectCols: "-2", // Выбрать с начала и до второй
			has: [][]byte{
				[]byte("col1"),
				[]byte("col2"),
				[]byte("col3"),
			},
			want: [][]byte{
				[]byte("col1"),
				[]byte("col2"),
			},
		},
		{
			selectCols: "2-", // Выбрать с второй и до конца
			has: [][]byte{
				[]byte("col1"),
				[]byte("col2"),
				[]byte("col3"),
			},
			want: [][]byte{
				[]byte("col2"),
				[]byte("col3"),
			},
		},
		{
			selectCols: "1-3", // С первой по третью
			has: [][]byte{
				[]byte("col1"),
				[]byte("col2"),
				[]byte("col3"),
			},
			want: [][]byte{
				[]byte("col1"),
				[]byte("col2"),
				[]byte("col3"),
			},
		},
		{
			selectCols: "1, 3, 2-", // Первую, третью, вторую и третью
			has: [][]byte{
				[]byte("col1"),
				[]byte("col2"),
				[]byte("col3"),
			},
			want: [][]byte{
				[]byte("col1"),
				[]byte("col3"),
				[]byte("col2"),
				[]byte("col3"),
			},
		},
		{ // Если колонок слишком мало -> nil
			selectCols: "3",
			has: [][]byte{
				[]byte("col1"),
				[]byte("col2"),
			},
			want: nil,
		},
		{ // Если разделителей нет, то пропускаем
			selectCols: "3",
			has: [][]byte{
				[]byte("line without seps"),
			},
			want: [][]byte{
				[]byte("line without seps"),
			},
		},
		{ // Пустая строка, должна вернутся пустая строка
			selectCols: "3",
			has: [][]byte{
				[]byte(""),
			},
			want: [][]byte{
				[]byte(""),
			},
		},
	}

	for _, tc := range testCases {
		filter := NewFilterSelectedCols(tc.selectCols)
		got := filter.Filter(tc.has)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("\ngot:  %v\nwant: %v", got, tc.want)
		}
	}
}

func TestFilterSelectedColsOpts(t *testing.T) {
	testCases := []struct {
		selectCols string
		skip       bool
		has        [][]byte
		want       [][]byte
	}{
		{
			selectCols: "3", // Если разделителей нет, то не пропускаем
			skip:       true,
			has: [][]byte{
				[]byte("line without seps"),
			},
			want: nil,
		},
	}

	for _, tc := range testCases {
		filter := NewFilterSelectedColsWithOpts(tc.selectCols, FilterSelectedColsOpts{
			skip: tc.skip,
		})
		got := filter.Filter(tc.has)
		if !reflect.DeepEqual(got, tc.want) {
			t.Errorf("\ngot:  %v\nwant: %v", got, tc.want)
		}
	}
}

func TestMergeByteSep(t *testing.T) {
	sep := byte('\t')
	has := [][]byte{
		[]byte("col1"),
		[]byte("col2"),
		[]byte("col3"),
	}
	want := []byte("col1\tcol2\tcol3")

	merger := MergeByteSep{sep: sep}
	got := merger.Merge(has)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:  %v\nwant: %v", got, want)
	}
}

func TestCut(t *testing.T) {
	testCases := []struct {
		desc string
		opts CutOpts
		has  []byte
		want []byte
	}{
		{ // just cat
			desc: "just cat",
			opts: CutOpts{
				Fields:    "2",
				Delimiter: '\t',
				Separated: false,
			},
			has:  []byte("line1\nline2\nline3\n"),
			want: []byte("line1\nline2\nline3\n"),
		},
		{ // just cat with -s
			desc: "just cat with -s",
			opts: CutOpts{
				Fields:    "2",
				Delimiter: '\t',
				Separated: true,
			},
			has:  []byte("line1\nline2\nline3\n"),
			want: nil,
		},
		{ // print second col
			desc: "print second col",
			opts: CutOpts{
				Fields:    "2",
				Delimiter: '\t',
				Separated: false,
			},
			has:  []byte("line1\tcol2\nline2\n"),
			want: []byte("col2\nline2\n"),
		},
		{ // print second col with -s
			desc: "print second col with -s",
			opts: CutOpts{
				Fields:    "2",
				Delimiter: '\t',
				Separated: true,
			},
			has:  []byte("line1\tcol2\nline2\n"),
			want: []byte("col2\n"),
		},
		{ // eof
			desc: "eof",
			opts: CutOpts{
				Fields:    "2",
				Delimiter: '\t',
				Separated: false,
			},
			has:  []byte("line1"),
			want: []byte("line1\n"),
		},
	}

	for _, tc := range testCases {
		cut := NewCut(&tc.opts)
		var got bytes.Buffer
		err := cut.Cut(bytes.NewBuffer(tc.has), &got)
		if err != nil {
			t.Fatalf("err should be nil, but: %s", err)
		}
		if !reflect.DeepEqual(got.Bytes(), tc.want) {
			t.Errorf("desc: %s\ngot:  %v\nwant: %v", tc.desc, got.Bytes(), tc.want)
		}
	}
}
