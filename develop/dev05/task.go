package exercise_5

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"testing"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

var after uint64
var before uint64
var context uint64
var count bool
var ignore bool
var invert bool
var fixed bool
var numerate bool

func init() {
	testing.Init()
	flag.Uint64Var(&after, "A", 0, "")
	flag.Uint64Var(&before, "B", 0, "")
	flag.Uint64Var(&context, "C", 0, "")
	flag.BoolVar(&count, "c", false, "")
	flag.BoolVar(&ignore, "i", false, "")
	flag.BoolVar(&invert, "v", false, "")
	flag.BoolVar(&fixed, "F", false, "")
	flag.BoolVar(&numerate, "n", false, "")
	flag.Parse()
}

// parseArgs объединяет патерны в одно регулярное выражение с учетом -i -F,
// возвращает regex и имя файла.
func parseArgs(args []string) (re *regexp.Regexp, file string, err error) {
	switch {
	case len(args) == 0:
		return nil, "", errors.New("no pattern specified")
	case len(args) == 1:
		ps := args[0]
		if fixed {
			ps = `\Q` + ps + `\E`
		}
		if ignore {
			ps = `(?i)(` + ps + `)`
		}
		re, err := regexp.Compile(ps)
		if err != nil {
			return nil, "", err
		}
		return re, "", nil
	default:
		st, en := `(`, `)`
		if fixed {
			st, en = `(\Q`, `\E)`
		}
		ps := st + args[0] + en
		for i := 1; i < len(args)-1; i++ {
			ps += `|` + st + args[i] + en
		}
		if ignore {
			ps = `(?i)(` + ps + `)`
		}
		re, err := regexp.Compile(ps)
		if err != nil {
			return nil, "", err
		}
		return re, args[len(args)-1], nil
	}
}

func maxUint64(a, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type scanner struct {
	sc *bufio.Scanner
	f  *os.File
}

// newScanner возвращает сканер для работы с файлом,
// если имя файла - пустая строка, то сканер работает с os.Stdin
func newScanner(file string) (*scanner, error) {
	if file == "" {
		sc := bufio.NewScanner(os.Stdin)
		return &scanner{sc: sc}, nil
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	sc := bufio.NewScanner(f)
	return &scanner{sc, f}, nil
}

// readLines возвращет слайс строк из файла как []byte
func readLines(scan *scanner) (lines [][]byte) {
	for scan.sc.Scan() {
		lines = append(lines, scan.sc.Bytes())
	}
	return lines
}

// filter возвращает номера строк, соответствущих регулярному
// выражению с учетом инверсии
func filter(lines [][]byte, re *regexp.Regexp) (s []int) {
	for i, line := range lines {
		if re.Match(line) != invert {
			s = append(s, i)
		}
	}
	return s
}

// printLine печатает в Stdout строку, а также номер строки, если предоставлен флаг -n.
// Строки, соответствубщие паттернам, отделяются от номера символом ":", остальные - символом "-"
func printLine(num int, line []byte, match bool) {
	switch {
	case numerate && match:
		fmt.Printf("%d:%s\n", num, line)
	case numerate && !match:
		fmt.Printf("%d-%s\n", num, line)
	default:
		fmt.Println(string(line))
	}
}

// printLines печатает строки lines с номерами из n в соответствии со значениями,
// переданными флагами -A, -B, -C.
func printLines(lines [][]byte, n []int) {
	h := 0
	for i := 0; i < len(n)-1; i++ {
		b := maxInt(h, n[i]-int(before))
		a := minInt(n[i]+int(after), n[i+1]-1)
		for j := b; j <= a; j++ {
			printLine(j+1, lines[j], n[i] == j)
		}
		h = a + 1
	}
	b := maxInt(h, n[len(n)-1]-int(before))
	a := minInt(n[len(n)-1]+int(after), len(lines)-1)
	for j := b; j <= a; j++ {
		printLine(j+1, lines[j], n[len(n)-1] == j)
	}
}

func main() {
	// флаги -A, -B, -C могут использоваться совместно
	after = maxUint64(after, context)
	before = maxUint64(before, context)

	re, file, err := parseArgs(flag.Args())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	scan, err := newScanner(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	defer scan.f.Close()

	lines := readLines(scan)
	filtered := filter(lines, re)

	// Если предоставлен флаг -c, печатается только количество строк, соответствуюших патернам
	if count {
		fmt.Println(len(filtered))
	} else if len(filtered) > 0 {
		printLines(lines, filtered)
	}
}
