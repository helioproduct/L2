package main

/*
=== Задача на распаковку ===

Создать Go функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы / руны, например:
	- "a4bc2d5e" => "aaaabccddddde"
	- "abcd" => "abcd"
	- "45" => "" (некорректная строка)
	- "" => ""
Дополнительное задание: поддержка escape - последовательностей
	- qwe\4\5 => qwe45 (*)
	- qwe\45 => qwe44444 (*)
	- qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка функция должна возвращать ошибку. Написать unit-тесты.

Функция должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"fmt"
	// "regexp"
	"strconv"
	"strings"
)

func UnpackString(s string) string {

	var count string
	var rn rune
	var unpacked string
	var escape bool

	runes := []rune(s + " ")
	for i := 0; i < len(runes); i++ {

		digit, err := strconv.Atoi(string(runes[i]))
		if err == nil && !escape {
			count += strconv.Itoa(digit)
		} else {
			repeatTimes, err := strconv.Atoi(count)
			if err != nil {
				repeatTimes = 1
			}
			if rn != 0 {
				unpacked += strings.Repeat(string(rn), repeatTimes)
			}
			rn = runes[i]
			count = ""
		}
	}
	// fmt.Println(unpacked)
	return unpacked
}

func main() {
	s := "\t"
	for _, rn := range []rune(s) {
		fmt.Printf("%c\n", rn)
	}
	// UnpackString(s)
	// fmt.Println([]byte(s))
}
