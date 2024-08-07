package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать простейший telnet-клиент.

Примеры вызовов:
go-telnet --timeout=10s host port
go-telnet mysite.ru 8080
go-telnet --timeout=3s 1.1.1.1 123


Требования:
1. 	Программа должна подключаться к указанному хосту (ip или доменное имя + порт) по протоколу TCP.
	После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT

2.	Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s)

3.	При нажатии Ctrl+D программа должна закрывать сокет и завершаться.
	Если сокет закрывается со стороны сервера, программа должна также завершаться.
	При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

/*
go-telnet [OPTION...] <host> <port>

Устанавливает tcp соединение с host:port

OPTIONS
		-timeout <duration> - устанавливает задержку в duration. (default 10s)
			<duration> = <number><suffix>
			<suffix> = ms|s|m
			ms - миллисекунда
			s - секунда
			m - минута
*/

type Config struct {
	Host    string
	Port    int
	Timeout time.Duration
}

func parseConfig() *Config {
	var cfg Config

	flag.DurationVar(&cfg.Timeout, "timeout", 10*time.Second, "устанавливает максимальное время подключения.")
	flag.Parse()

	// parse host port
	args := flag.Args()
	if len(args) != 2 {
		flag.Usage()
		log.Fatalln("len(args) != 2")
	}
	cfg.Host = args[0]
	var err error
	cfg.Port, err = strconv.Atoi(args[1])
	if err != nil {
		flag.Usage()
		log.Fatalln("port should be a number")
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		flag.Usage()
		log.Fatalln("port should be in [1, 65535]")
	}

	return &cfg
}

func echo(r io.Reader, w io.Writer) error {
	_, err := io.Copy(w, r)
	return err
}

func telnet(cfg *Config) error {
	// notify sigint
	fmt.Println(cfg)
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// Попытаемся установить соединение
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	log.Printf("Trying %s...\n", addr)
	conn, err := net.DialTimeout("tcp", addr, cfg.Timeout)
	if err != nil {
		return err
	}
	log.Printf("Connected to %s\n", addr)

	senderDone := make(chan error)
	receiverDone := make(chan error)

	go func() {
		senderDone <- echo(os.Stdin, conn)
	}()

	go func() {
		receiverDone <- echo(conn, os.Stdout)
	}()

	select {
	case err := <-receiverDone:
		if err != nil {
			fmt.Println("receiver err:", err)
		}
		close(receiverDone)
	case err := <-senderDone:
		if err != nil {
			fmt.Println("sender err:", err)
		}
		close(senderDone)
	case s := <-sig:
		log.Println("receive os signal:", s.String())
	}

	conn.Close()
	if _, ok := <-receiverDone; ok {
		close(receiverDone)
	}
	if _, ok := <-senderDone; ok {
		close(senderDone)
	}

	return nil
}

func main() {
	cfg := parseConfig()
	if err := telnet(cfg); err != nil {
		log.Fatalln(err)
	}
}
