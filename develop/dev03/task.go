package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"unicode"
)

/*
Отсортировать строки в файле по аналогии с консольной утилитой sort
(man sort — смотрим описание и основные параметры): на входе подается файл из несортированными строками, на выходе — файл с отсортированными.

Реализовать поддержку утилитой следующих ключей:

-k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел)
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительно

Реализовать поддержку утилитой следующих ключей:

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учетом суффиксов
*/

type Config struct {
	SortedColumn     int  // Номер колонки, по которой нужно сортировать -k
	NumSort          bool // Числовая сортировка -n
	Reversed         bool // Сортировка в обратном порядке -r
	Unique           bool // Не выводить повторяющиеся строки -u
	MonthSort        bool // Сортировать по месяцам -M
	Strip            bool // Игнорировать хвостовые пробелы -b
	CheckSorted      bool // Проверить, отсортированны ли данные
	HumanNumericSort bool // Сортировать с учетом суффиксов h

	InPath  string // Входной файл IN
	OutPath string // Выходной файл OUT
}

func parseArgs() *Config {
	var cfg Config
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `NAME:
	sort - сортирует строки в текстовых файлах	
SYNOPSIS:
	sort [OPTION...] [IN] [OUT]

DESCRIPTION:
	Сортирует строки из файла IN и сохраняет в OUT

	Если IN не указан, то считывает строки из stdin до его закрытия
	Если указан только один файл, то он считается IN
	Если OUT не указан, то выводит результат работы в stdout
	Если не задана сортировка, то сортирует в лексикографическом порядке

	Опции:

	-b 
		Игнорировать хвостовые пробелы (лидирующие и хвостовые)
	
	-c 
		Провкра на отсортированность (не сортирует)
	
	-k k=COL
		Сортировать по колонке COL
		
		Строка разбивается по пробельным символам, 
		Если использовать с ключом b, то не будут учитываться колонки 
		с конца и в начале

		Если в строке количетсво колонок меньше нужного, то ключ будет равен "", а утилита выдаст предупреждение в stderr

		Колонки нумеруются с 1, слева направо.
	
	-M
		Сравнение по месяцам (unknown) < 'JAN' < ... < 'DEC'

	-n
		Сравнение по строково-числовому значению
		В качестве числа берется максимальный целочисленный префикс

		Если в строке нет числа, то она сортируется так, 
		если бы у нее был префикс "0"

	-r
		Производит сортировку в обратном направлении. Опция применяется последней
	
	-u
		Выводит только уникальные строки

	Следующие опции являются взаимоисключающими:
		-n, -M
	Утилита сразу же завершится, если одновременно задано несколько типов сортировки.

	Утилита использует внутренюю сортировку, поэтому лучше не стоит сортировать большие файлы (Вес которых эквивалентен с объемом RAM).
`)
	}
	flag.IntVar(&cfg.SortedColumn, "k", 0, "указание колонки для сортировки")
	flag.BoolVar(&cfg.NumSort, "n", false, "сортировать по числовому значению")
	flag.BoolVar(&cfg.Reversed, "r", false, "сортировать в обратном порядке")
	flag.BoolVar(&cfg.Unique, "u", false, "не выводить повторяющиеся строки")
	flag.BoolVar(&cfg.MonthSort, "M", false, "сортировать по названию месяца")
	flag.BoolVar(&cfg.Strip, "b", false, "игнорировать хвостовые пробелы")
	flag.BoolVar(&cfg.CheckSorted, "c", false, "проверять отсортированы ли данные")
	flag.Parse()

	// Проверка на неотрицательность колонки
	if cfg.SortedColumn < 0 {
		flag.Usage()
		os.Exit(2)
	}

	// Считывание названий файлов
	args := flag.Args()
	if len(args) > 2 {
		flag.Usage()
		os.Exit(2)
	}
	if len(args) >= 1 {
		cfg.InPath = args[0]
	}
	if len(args) >= 2 {
		cfg.OutPath = args[1]
	}

	// Проверка на взаимоисключающие ключи
	check := func() func(bool) {
		var f bool
		return func(b bool) {
			if f && b {
				flag.Usage()
				os.Exit(2)
			}
			f = f || b
		}
	}()
	check(cfg.NumSort)
	check(cfg.MonthSort)
	check(cfg.HumanNumericSort)

	return &cfg
}

func getInput(cfg *Config) (*bufio.Scanner, func() error) {
	var in *bufio.Scanner
	if cfg.InPath != "" {
		if f, err := os.Open(cfg.InPath); err == nil {
			in = bufio.NewScanner(f)
			return in, f.Close
		} else {
			log.Fatalf("Error: %s", err)
		}
	}

	in = bufio.NewScanner(os.Stdin)
	return in, func() error { return nil }
}

func getOut(cfg *Config) (*bufio.Writer, func() error) {
	var out *bufio.Writer
	if cfg.OutPath != "" {
		if f, err := os.Create(cfg.OutPath); err == nil {
			out = bufio.NewWriter(f)
			return out, func() error {
				err1 := out.Flush()
				err2 := f.Close()
				return errors.Join(err1, err2)
			}
		} else {
			log.Fatalf("Error: %s", err)
		}
	}
	out = bufio.NewWriter(os.Stdout)
	return out, func() error {
		return out.Flush()
	}
}

var MONTHS = map[string]int{
	"JAN": 1,
	"FEB": 2,
	"MAR": 3,
	"APR": 4,
	"MAY": 5,
	"JUN": 6,
	"JUL": 7,
	"AUG": 8,
	"SEP": 9,
	"OCT": 10,
	"NOV": 11,
	"DEC": 12,
}

type lessFunc func(i, j int) bool
type keyFunc func(i int) string
type sortByLess struct {
	lines []string
	less  lessFunc
	key   keyFunc
}

func parsePrefixInt(s string) int {
	res := 0
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return res
		}
		res = res*10 + int(r-'0')
	}

	return res
}

func makeSBL(cfg *Config) *sortByLess {
	var sbl sortByLess

	// Модификаторы ключа
	sbl.key = func(i int) string { return sbl.lines[i] }
	if cfg.Strip { // -b
		sbl.key = func(i int) string {
			return strings.TrimSpace(sbl.lines[i])
		}
	}
	if cfg.SortedColumn > 0 { //-k
		prevKey := sbl.key
		col := cfg.SortedColumn - 1
		sbl.key = func(i int) string {
			k := prevKey(i)
			splited := strings.Split(k, " ")
			if len(splited) > col {
				return splited[col]
			} else {
				return ""
			}
		}
	}

	// Модификаторы сравнения
	switch {
	case cfg.MonthSort:
		sbl.less = func(i, j int) bool {
			a, b := sbl.key(i), sbl.key(j)
			ai, ok := MONTHS[a]
			if !ok {
				ai = 0
			}
			bi, ok := MONTHS[b]
			if !ok {
				bi = 0
			}
			return ai < bi
		}
	case cfg.NumSort: // строково-числовая сортировка
		sbl.less = func(i, j int) bool {
			a, b := sbl.key(i), sbl.key(j)
			ai := parsePrefixInt(a)
			bi := parsePrefixInt(b)

			if ai != bi {
				return ai < bi
			}
			return a < b
		}
	default: // Лексиграфическая сортировка
		sbl.less = func(i, j int) bool {
			return sbl.key(i) < sbl.key(j)
		}
	}
	return &sbl
}

func (sbl *sortByLess) Len() int {
	return len(sbl.lines)
}
func (sbl *sortByLess) Swap(i, j int) {
	sbl.lines[i], sbl.lines[j] = sbl.lines[j], sbl.lines[i]
}
func (sbl *sortByLess) Less(i, j int) bool {
	return sbl.less(i, j)
}

func main() {
	cfg := parseArgs()

	// Открываем файлы для чтения и записи
	scn, closeIn := getInput(cfg)
	defer closeIn()
	out, closeOut := getOut(cfg)
	defer closeOut()

	// Строим интерфейс sort.Interface
	sbl := makeSBL(cfg)

	// Считываем строчки
	for scn.Scan() {
		sbl.lines = append(sbl.lines, scn.Text())
	}

	// Если нужно только проверить на отсортированность:
	if cfg.CheckSorted {
		{ // Обворачиваем функцию less
			old := sbl.less
			sbl.less = func(i, j int) bool {
				f := old(i, j)
				if f == false {
					fmt.Fprintf(out, "lines %d and %d are unordered", i+1, j+1)
					os.Exit(1)
				}
				return f
			}
		}
		sort.IsSorted(sbl)
		return
	}

	if cfg.Reversed {
		sort.Sort(sort.Reverse(sbl))
	} else {
		sort.Sort(sbl)
	}

	// Если нужно вывести не повторяющиеся строчки
	if cfg.Unique {
		fmt.Fprintln(out, sbl.lines[0])
		for i := 1; i < len(sbl.lines); i++ {
			if sbl.lines[i-1] != sbl.lines[i] { // i указывает на начло нового блока
				fmt.Fprintln(out, sbl.lines[i]) // Выводим строку из нового блока
			}
		}
		return
	}

	// Если нужно просто вывести строчки
	for _, l := range sbl.lines {
		fmt.Fprintln(out, l)
	}
}
