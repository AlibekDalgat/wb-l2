// main_test.go
package main

import (
	"os"
	"testing"
	"time"

	"github.com/beevik/ntp"
)

func TestMainFunction(t *testing.T) {
	// перед запуском тестируемой функции надо подготовить буферы, куда будут отправляться информация
	// сохранение оригинальных адресов на os.Stdout и os.Stderr, чтобы потом можно было их восстановить
	originalStdout := os.Stdout
	originalStderr := os.Stderr

	// временные буферы для stdout и stderr
	rStdout, wStdout, _ := os.Pipe()
	rStderr, wStderr, _ := os.Pipe()

	// перенаправляем stdout и stderr в временные буферы
	os.Stdout = wStdout
	os.Stderr = wStderr

	// запуск тестируемой функцию
	main()

	// закрытие временных буферов
	wStdout.Close()
	wStderr.Close()

	// восстанавление оригинальных адресов stdout и stderr
	os.Stdout = originalStdout
	os.Stderr = originalStderr

	// чтение данные из временных буферов
	stdoutBytes := make([]byte, 1024)
	stdoutLen, _ := rStdout.Read(stdoutBytes)
	stderrBytes := make([]byte, 1024)
	stderrLen, _ := rStderr.Read(stderrBytes)

	// проверка что не было ошибок при выполнении программы
	if stdoutLen == 0 || stderrLen > 0 {
		t.Errorf("ошибка при выполнения программы stdout: %s, stderr: %s", string(stdoutBytes[:stdoutLen]), string(stderrBytes[:stderrLen]))
		return
	}

	// проверка формата вывода времени
	_, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		t.Errorf("ошибка при получении времени с сервера NTP: %v", err)
		return
	}

	expectedOutput := []byte(time.Now().Format(time.RFC3339))
	if string(stdoutBytes[:stdoutLen]) != string(expectedOutput) {
		t.Errorf("ожидаемый вывод:\n%s\nполученный вывод:\n%s", expectedOutput, stdoutBytes[:stdoutLen])
	}
}
