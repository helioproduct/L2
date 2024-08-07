package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

/*
Создать Go-функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы/руны, например:
"a4bc2d5e" => "aaaabccddddde"
"abcd" => "abcd"
"45" => "" (некорректная строка)
"" => ""

Дополнительно
Реализовать поддержку escape-последовательностей.
Например:
qwe\4\5 => qwe45 (*)
qwe\45 => qwe44444 (*)
qwe\\5 => qwe\\\\\ (*)


В случае если была передана некорректная строка, функция должна возвращать ошибку. Написать unit-тесты.
*/

// Распаковка
func Unpack(s string) (string, error) {
	var res strings.Builder
	var num int         // Сколько раз нужно повторить char
	var char rune       // Char, который нужно записать
	var isInitChar bool // Было ли присвоено char какое-то значение
	var openEscape bool // Нужно ли экранировать символ?

	for _, r := range s {

		if unicode.IsDigit(r) && !openEscape {
			if !isInitChar { // Если до числа не было символов
				return "", errors.New("a string cannot start with a number")
			}
			num = num*10 + int(r-'0')
		} else if r == '\\' && !openEscape { // Встречаем escape символ
			openEscape = true
		} else {
			if isInitChar {
				res.WriteRune(char)
				for ; num > 1; num-- {
					res.WriteRune(char)
				}
			}
			num = 0
			char = r
			openEscape = false
			isInitChar = true
		}

	}

	if openEscape { // При выходе из цикла - последовательность должна быть закрыта
		return "", errors.New("incomplete escape sequence")
	}

	if isInitChar {
		res.WriteRune(char)
		for ; num > 1; num-- {
			res.WriteRune(char)
		}
	}

	return res.String(), nil
}

func main() {
	s := `\02\7`
	if out, err := Unpack(s); err == nil {
		fmt.Printf("%s -> %s", s, out)
	} else {
		fmt.Println(err)
	}
}
