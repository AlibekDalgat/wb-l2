package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

/*
=== Утилита telnet ===

Реализовать примитивный telnet клиент:
Примеры вызовов:
go-telnet --timeout=10s host port go-telnet mysite.ru 8080 go-telnet --timeout=3s 1.1.1.1 123

Программа должна подключаться к указанному хосту (ip или доменное имя) и порту по протоколу TCP.
После подключения STDIN программы должен записываться в сокет, а данные полученные и сокета должны выводиться в STDOUT
Опционально в программу можно передать таймаут на подключение к серверу (через аргумент --timeout, по умолчанию 10s).

При нажатии Ctrl+D программа должна закрывать сокет и завершаться. Если сокет закрывается со стороны сервера, программа должна также завершаться.
При подключении к несуществующему сервер, программа должна завершаться через timeout.
*/

func main() {
	// флаги командной строки
	host := flag.String("host", "", "Хост (IP или доменное имя)")
	port := flag.Int("port", 0, "Порт")
	timeout := flag.Duration("timeout", 10*time.Second, "Таймаут подключения")
	flag.Parse()

	// проверка наличия хоста и порта
	if *host == "" || *port == 0 {
		fmt.Println("нужен хост и порт")
		// вывод описания флагов
		flag.PrintDefaults()
		return
	}

	// форматирование адреса для подключения
	addr := fmt.Sprintf("%s:%d", *host, *port)

	// процесс подключение к серверу
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Printf("ошибка подключения: %v\n", err)
		return
	}
	// установка таймаута
	conn.SetReadDeadline(time.Now().Add(*timeout))
	defer conn.Close()

	// чтения данных из сокета и вывода в stdout в отдельной горутине
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := conn.Read(buffer)
			// завершить программу если будет ошибка при чтении или окончания таймаута
			if err != nil {
				fmt.Printf("Ошибка чтения из сокета: %v\n", err)
				os.Exit(0)
			}
			fmt.Print(string(buffer[:n]))
		}
	}()

	// канал куда подастся сигнал завершения. Ctrl+D не работает, видимо, неподходящие свойства терминала
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)

	// ожидание сигнала завершения
	select {
	case <-sigCh:
		fmt.Println("Программа завершена по сигналу")
	}
}
