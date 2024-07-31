//go:build linux

package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"syscall"
)

// Настраиваем команды
func SetupCmds(intr *Interpreter) {
	intr.AddCmd("exit", &exit{})
	intr.AddCmd("cd", &cd{})
	intr.AddCmd("pwd", &pwd{})
	intr.AddCmd("echo", &echo{})
	intr.AddCmd("kill", &kill{})
	intr.AddCmd("ps", &ps{})

}

// Команда exit
type exit struct{}

func (_ *exit) Exec(r io.Reader, w io.Writer, args []string) int {
	os.Exit(0)
	return 0
}

// Команда cd <args>
type cd struct{}

func (_ *cd) Exec(r io.Reader, w io.Writer, args []string) int {
	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <path>\n", args[0])
		return 2
	}

	if err := os.Chdir(args[1]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

// Команда pwd
type pwd struct{}

func (_ *pwd) Exec(r io.Reader, w io.Writer, args []string) int {
	if wd, err := os.Getwd(); err == nil {
		fmt.Fprintln(w, wd)
	} else {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

// Команда echo <args>
type echo struct{}

func (_ *echo) Exec(r io.Reader, w io.Writer, args []string) int {
	s := strings.Join(args[1:], " ") + "\n"
	if _, err := w.Write([]byte(s)); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

// Команда kill <sig> <pid>
type kill struct{}

func (k *kill) ArgsFatal(args []string) int {
	fmt.Fprintf(os.Stderr, "usage: %s TERM|KILL <pid>\n", args[0])
	return 2
}
func (k *kill) Exec(r io.Reader, w io.Writer, args []string) int {
	if len(args) != 3 {
		return k.ArgsFatal(args)
	}
	var sig syscall.Signal
	switch args[1] {
	case "TERM":
		sig = syscall.SIGTERM
	case "KILL":
		sig = syscall.SIGKILL
	default:
		fmt.Fprintln(os.Stderr, "wrong sig")
		return k.ArgsFatal(args)
	}

	pid, err := strconv.Atoi(args[2])
	if err != nil {
		return k.ArgsFatal(args)
	}

	if err := syscall.Kill(pid, sig); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

// Команда ps
// out: PID TTY CMD
type ps struct{}
type proc struct {
	pid   int    // pid процесса
	name  string // имя запускаемого файла
	state byte   // состояние процесса
}

func (p *ps) setStats(n *proc) error {
	procstats := fmt.Sprintf("/proc/%d/stat", n.pid)
	data, err := os.ReadFile(procstats)
	if err != nil {
		return err
	}
	if n, err := fmt.Sscanf(string(data), "%d %s %c", &n.pid, &n.name, &n.state); err != nil || n != 3 {
		return err
	}
	return nil
}

func (p *ps) Exec(r io.Reader, w io.Writer, args []string) int {
	entries, err := os.ReadDir("/proc/")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	procs := []proc{}
	for _, entry := range entries {
		if entry.IsDir() {
			name := entry.Name()
			if pid, err := strconv.Atoi(name); err == nil {
				// Мы нашли процесс
				nproc := proc{
					pid: pid,
				}
				// Считываем статистику процесса
				if err := p.setStats(&nproc); err == nil {
					procs = append(procs, nproc)
				} else {
					// skip
				}
			}
		}
	}

	// Теперь выведем их
	_, err = fmt.Fprintf(w, "%6s %7s %8s\n", "PID", "STATE", "NAME")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	for _, proc := range procs {
		_, err := fmt.Fprintf(w, "%-10d %c    %s\n", proc.pid, proc.state, proc.name)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	}

	return 0
}
