//go:build linux

package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"syscall"
)

// С помощью фнукции Parse интепретатор получает сущность,
// которая хранит в себе информацию о том, что нужно сделать
type Parser interface {
	Parse() Entity
}

// Структура, которая представляет собой результат Parse.
type Entity struct {
	Cmds [][]string // Коамнды, которые нужно выполнить
	Bg   bool       // Нужно ли запустить их в фоновом режиме
	EOF  bool       // Если True, то это команда завершения работы
}

// Представялет из себя встроенную команду
type Executor interface {
	Exec(r io.Reader, w io.Writer, args []string) int
}

type forkFunc func() int

// Делаем форк и вызываем функцию forkFunc
// in -> forkFunc -> out
// ...............\_ err
func forkout(in, out, err *os.File, forkFunc forkFunc) (int, error) {
	pid, _, errno := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if errno != 0 {
		return 0, fmt.Errorf("can't fork: %d", errno)
	}
	if pid == 0 {
		// child
		// Заемняем stdin на in
		err := syscall.Dup2(int(in.Fd()), int(os.Stdin.Fd()))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// Заемняем stdout на out
		err = syscall.Dup2(int(out.Fd()), int(os.Stdout.Fd()))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		exitStatus := forkFunc()
		os.Exit(exitStatus)
	}
	return int(pid), nil
}

// Делаем форк и вызываем функцию forkFunc
// in -> forkFunc -> out
// ...............\_ err
// создает pipe out
func fork(in, errf *os.File, forkFunc forkFunc) (out *os.File, pid int, err error) {
	// Создаем Pipe
	r, w, err := os.Pipe()
	if err != nil {
		return nil, 0, err
	}

	// Делаем fork
	_pid, _, errno := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	pid = int(_pid)
	if errno != 0 {
		return nil, 0, fmt.Errorf("can't fork: %d", errno)
	}
	if pid == 0 {
		// child
		// Закрываем не нужный для ребенка reader
		r.Close()
		// stdin заменяем на in
		err := syscall.Dup2(int(in.Fd()), int(os.Stdin.Fd()))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		// stdout заменяем на out
		err = syscall.Dup2(int(w.Fd()), int(os.Stdout.Fd()))
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		exitStatus := forkFunc()
		os.Exit(exitStatus)
	}

	// Закрываем не нужный для родителя writer
	w.Close()

	return r, int(pid), nil
}

// Обворачивает Executor для fork
func wrapExecForFork(e Executor, args []string) forkFunc {
	return func() int {
		return e.Exec(os.Stdin, os.Stdout, args)
	}
}

// Отвечает за работу с парсером и управлениями командами
type Interpreter struct {
	parser   Parser
	commands map[string]Executor
}

func NewInterpreter() *Interpreter {
	var itrpr Interpreter
	itrpr.parser = NewDefaultParser(os.Stdin)
	itrpr.commands = make(map[string]Executor)
	return &itrpr
}

func (itrpr *Interpreter) AddCmd(name string, e Executor) {
	itrpr.commands[name] = e
}

// Выполняем команду и дожидаемся ее исполнения
// in -> cmd -> stdout
// ..........\_ stderr
func (itrpr *Interpreter) ewait(in *os.File, args []string) error {
	if len(args) == 0 {
		return nil
	}
	if e, ok := itrpr.commands[args[0]]; ok {
		// Встроенная команда
		exitStatus := e.Exec(in, os.Stdout, args)
		if exitStatus != 0 {
			return fmt.Errorf("exit status %d", exitStatus)
		}
	} else {
		// Используем exec, который является кроссплатформенной оберткой над fork, dup2, exec, wait
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = in
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		return err
	}
	return nil
}

// Создаем новый процесс, который выполняет команду
// in -> cmd -> out
// ..........\_ stderr
func (itrpr *Interpreter) eforkout(in *os.File, out *os.File, args []string) (int, error) {
	if len(args) == 0 {
		return 0, nil
	}

	if e, ok := itrpr.commands[args[0]]; ok {
		// Встроенная команда
		return forkout(in, out, os.Stderr, wrapExecForFork(e, args))
	} else {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = in
		cmd.Stderr = os.Stderr
		cmd.Stdout = out
		err := cmd.Start()
		return cmd.Process.Pid, err
	}
}

// Создаем новый процесс, который выполняет команду
// in -> cmd -> out
// ..........\_ stderr
func (itrpr *Interpreter) efork(in *os.File, args []string) (out *os.File, pid int, err error) {
	if len(args) == 0 {
		return in, 0, nil
	}

	if e, ok := itrpr.commands[args[0]]; ok {
		// Встроенная команда
		return fork(in, os.Stderr, wrapExecForFork(e, args))
	} else {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = in
		cmd.Stderr = os.Stderr

		if rp, err := cmd.StdoutPipe(); err == nil {
			out = rp.(*os.File) // TODO да, это не очень хорошо, но буду рассчитывать на неизменность модуля exec
		} else {
			return out, 0, err
		}
		err := cmd.Start()
		return out, cmd.Process.Pid, err
	}
}

func (itrpr *Interpreter) do(e Entity) error {
	if len(e.Cmds) == 0 {
		return errors.New("empty command")
	}

	in := os.Stdin
	defer func() {
		if in != os.Stdin {
			in.Close()
		}
	}()
	for _, cmd := range e.Cmds[:len(e.Cmds)-1] {
		// Нужно сделать fork и выполнить ее там
		out, pid, err := itrpr.efork(in, cmd)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "forked: %-10s pid: %6d\n", cmd[0], pid)

		if in != os.Stdin {
			in.Close() // Закрываем pipe, через который мы связывали два процесса
		}
		in = out
	}

	// Выполняем последнюю команду особенно
	lastArgs := e.Cmds[len(e.Cmds)-1]
	if e.Bg {
		// Нужно сделать fork и выполнить ее там
		pid, err := itrpr.eforkout(in, os.Stdout, lastArgs)
		if err != nil {
			return err
		}
		fmt.Fprintf(os.Stderr, "forked: %-10s pid: %6d\n", lastArgs[0], pid)
		return nil
	} else {
		return itrpr.ewait(in, lastArgs)
	}
}

func (itrpr *Interpreter) inviteFunc() string {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		pwd = ""
	}
	msg := fmt.Sprintf("%s$ ", pwd)
	return msg
}

func (itrpr *Interpreter) Start() {
	var e Entity
	fmt.Print(itrpr.inviteFunc())
	e = itrpr.parser.Parse()
	for !e.EOF {
		if err := itrpr.do(e); err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
		}
		fmt.Print(itrpr.inviteFunc())
		e = itrpr.parser.Parse()
	}
}
