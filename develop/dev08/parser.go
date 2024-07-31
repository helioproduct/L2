package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type DefaultParser struct {
	scn *bufio.Scanner
}

func NewDefaultParser(reader io.Reader) *DefaultParser {
	return &DefaultParser{
		scn: bufio.NewScanner(reader),
	}
}

func (p *DefaultParser) Parse() Entity {
	var e Entity

	for ok := true; ok; {
		ok = false
		e = Entity{}
		if !p.scn.Scan() {
			return Entity{
				EOF: true,
			}
		}
		line := p.scn.Text()
		tokens := strings.Fields(line)

		var currcmd []string
		for i := 0; i < len(tokens); i++ {
			if tokens[i] == "&" {
				if len(tokens)-1 != i {
					fmt.Fprintln(os.Stderr, "& может быть только в конце.")
					ok = true // Нужно еще раз проитерироваться
					break
				}
				e.Bg = true
			} else if tokens[i] == "|" { // Если встретили канал, то сохраняем команду
				e.Cmds = append(e.Cmds, currcmd)
				currcmd = nil // Обнуляем массив
			} else { // Добавляем токен в команду
				currcmd = append(currcmd, tokens[i])
			}
		}
		e.Cmds = append(e.Cmds, currcmd)
	}

	return e
}
