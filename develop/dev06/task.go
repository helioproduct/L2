package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// cut [OPTION...]

type Splitter interface {
	Split(b []byte) [][]byte
}

// Разделяет байты по байтовому сепаратору
type SplitByteSep struct {
	sep byte
}

func NewSplitByteSep(sep byte) *SplitByteSep {
	return &SplitByteSep{
		sep: sep,
	}
}

func (sb *SplitByteSep) Split(b []byte) [][]byte {
	var cols [][]byte
	s := 0
	for i, e := range b {
		if e == sb.sep {
			cols = append(cols, b[s:i])
			s = i + 1
		}
	}
	cols = append(cols, b[s:])

	return cols
}

type Filter interface {
	Filter(cols [][]byte) [][]byte
}

// Фильтрует заданные колонки
type FilterSelectedCols struct {
	selcols [][2]int // [[left right],...]
	skip    bool
}

// Создает filter
func NewFilterSelectedCols(selectedCols string) *FilterSelectedCols {
	colsParser := NewDefaultListParser()

	return &FilterSelectedCols{
		selcols: colsParser.Parse(selectedCols),
	}
}

type FilterSelectedColsOpts struct {
	skip bool // не пропускает строки без разделителей
}

// Создает filter со скипом строк без разделителей
func NewFilterSelectedColsWithOpts(selectedCols string, opts FilterSelectedColsOpts) *FilterSelectedCols {
	colsParser := NewDefaultListParser()

	return &FilterSelectedCols{
		selcols: colsParser.Parse(selectedCols),
		skip:    opts.skip,
	}
}

func (f *FilterSelectedCols) Filter(cols [][]byte) [][]byte {
	// Если мы не выбираем строчки, то сразу скипаем
	if len(f.selcols) == 0 {
		return nil
	}

	// Если разделителей нет и f.skip (скипнуть такую строчку) - скипаем
	if f.skip && len(cols) < 2 {
		return nil
	}
	// Если разделителей нет, то пропускаем строчку
	if len(cols) == 1 {
		return cols
	}

	var outCols [][]byte
	for _, r := range f.selcols {
		l, r := r[0]-1, r[1] // колонки в промежутке [ )
		if r == -1 {
			r = len(cols)
		}

		// Если некорректный промежуток, то пропускаем его
		validRange := l >= 0 && l < r && r <= len(cols)
		if !validRange {
			continue
		}

		outCols = append(outCols, cols[l:r]...) // Добавляем промежутки
	}

	return outCols
}

// Парсит LIST на последовательность интервалов [[left right],...]
type ListParser interface {
	Parse(list string) [][2]int
}

// Each LIST is made up of one range, or many ranges separated by commas.
// Selected  input  is written in the same order that it is read, and is written exactly once.  Each range is one
// of:
// N      N'th byte, character or field, counted from 1
// N-     from N'th byte, character or field, to end of line
// N-M    from N'th to M'th (included) byte, character or field
// -M     from first to M'th (included) byte, character or field
type DefaultListParser struct {
	rer      *regexp.Regexp
	redigits *regexp.Regexp
}

func NewDefaultListParser() *DefaultListParser {
	return &DefaultListParser{
		rer:      regexp.MustCompile(`(\d+-\d+|\d+-|-?\d+)`), // Находит промежутки
		redigits: regexp.MustCompile(`\d+`),
	}
}

func (p *DefaultListParser) Parse(list string) [][2]int {
	var selcols [][2]int

	ranges := p.rer.FindAllString(list, -1)

	for _, r := range ranges {
		digits := p.redigits.FindAllString(r, 2) // Находим числа

		if len(digits) == 2 { // N-M
			a, err := strconv.Atoi(digits[0])
			if err != nil {
				panic(digits[0] + " is not a number!")
			}
			b, err := strconv.Atoi(digits[1])
			if err != nil {
				panic(digits[1] + " is not a number!")
			}
			selcols = append(selcols, [2]int{a, b})
		} else if r[0] == '-' { // -M
			b, err := strconv.Atoi(digits[0])
			if err != nil {
				panic(digits[0] + " is not a number!")
			}
			selcols = append(selcols, [2]int{1, b})
		} else if r[len(r)-1] == '-' { // N-
			a, err := strconv.Atoi(digits[0])
			if err != nil {
				panic(digits[0] + " is not a number!")
			}
			selcols = append(selcols, [2]int{a, -1})
		} else { // N
			a, err := strconv.Atoi(digits[0])
			if err != nil {
				panic(digits[0] + " is not a number!")
			}
			selcols = append(selcols, [2]int{a, a})
		}
	}

	return selcols
}

type Merger interface {
	Merge(cols [][]byte) []byte
}

// Сливает колонки с указанным разделителем
type MergeByteSep struct {
	sep byte
}

func NewMergeByteSep(sep byte) *MergeByteSep {
	return &MergeByteSep{
		sep: sep,
	}
}

func (m *MergeByteSep) Merge(cols [][]byte) []byte {
	if len(cols) == 0 {
		return nil
	}

	var buf bytes.Buffer
	buf.Grow(len(cols[0]))
	buf.Write(cols[0])

	for i := 1; i < len(cols); i++ {
		buf.Grow(len(cols[i]) + 1)
		buf.WriteByte(m.sep)
		buf.Write(cols[i])
	}

	return buf.Bytes()
}

type Cut struct {
	splitter Splitter
	filter   Filter
	merger   Merger
}

func NewCut(opts *CutOpts) *Cut {
	sep := opts.Delimiter
	selcols := opts.Fields

	return &Cut{
		splitter: NewSplitByteSep(sep),
		filter: NewFilterSelectedColsWithOpts(selcols, FilterSelectedColsOpts{
			skip: opts.Separated,
		}),
		merger: NewMergeByteSep(sep),
	}
}

func (c *Cut) Cut(input io.Reader, output io.Writer) error {
	linesep := byte('\n')
	reader := bufio.NewReader(input)

	var err error
	linesepBytes := []byte{linesep}
	for err != io.EOF {
		var line []byte
		line, err = reader.ReadBytes(linesep)
		if err != nil && err != io.EOF {
			return err
		}
		if err != io.EOF {
			// Удаляем перенос строки
			line = line[:len(line)-1]
		}
		if err == io.EOF && len(line) == 0 {
			continue
		}

		sepline := c.splitter.Split(line)
		filtercols := c.filter.Filter(sepline)
		if filtercols == nil {
			continue
		}
		mergedline := c.merger.Merge(filtercols)

		if _, err := output.Write(mergedline); err != nil {
			return err
		}
		// Возвращаем перенос строки
		if _, err := output.Write(linesepBytes); err != nil {
			return err
		}
	}
	return nil
}

type CutOpts struct {
	Fields    string
	Delimiter byte
	Separated bool
}

func parseArgs() *CutOpts {
	var cutOpts CutOpts
	var del string

	flag.StringVar(&cutOpts.Fields, "f", "", "выбрать поля (колонки)")
	flag.StringVar(&del, "d", "", "использовать другой разделитель")
	flag.BoolVar(&cutOpts.Separated, "s", false, "только строки с разделителем")
	flag.Parse()

	if len(del) > 1 {
		fmt.Fprintln(os.Stderr, "delimiter must be a byte")
		flag.Usage()
		os.Exit(2)
	} else if len(del) == 1 {
		cutOpts.Delimiter = del[0]
	} else {
		cutOpts.Delimiter = '\t'
	}

	return &cutOpts
}

func main() {
	cfg := parseArgs()
	cut := NewCut(cfg)

	err := cut.Cut(os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}

}
