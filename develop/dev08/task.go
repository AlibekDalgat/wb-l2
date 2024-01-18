package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// функция для перемещения между директориями
func osCD(args []string) error {
	// если нет аргументов для команды вернуть ошибку
	if len(args) < 2 {
		return errors.New("cd <directory>")
	}
	// смена директории
	err := os.Chdir(args[1])
	return err
}

// функция для показа текущегокаталога
func osPWD() (string, error) {
	// получение корневогопути до текущего каталога
	dir, err := os.Getwd()
	return dir, err
}

// функция вывода аргумента в STDOUT
func osECHO(args []string) {
	fmt.Println(strings.Join(args[1:], " "))
}

// функция терменирования процесса
func osKILL(args []string) error {
	// если нет аргументов для команды вернуть ошибку
	if len(args) < 2 {
		return errors.New("kill <pid>")
	}
	pidStr := args[1]
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return errors.New("Invalid PID")
	}
	cmd := exec.Command("kill", "-TERM", strconv.Itoa(pid))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// возвращение ошибки (или её отсутствия) запуска убийства процесса
	return cmd.Run()
}

// функция выводаобщей информации о запущенных процессах
func osPS() error {
	cmd := exec.Command("ps", "aux")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return err
}

func main() {
	// бесконечный ввод
	for {
		// приглашение
		fmt.Print("> ")
		// считывание команд
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()
		// если введена командавыход завершить цикл
		if input == "\\quit" {
			break
		}

		args := strings.Fields(input)
		if len(args) == 0 {
			continue
		}
		// обработка введённой команды
		switch args[0] {
		case "cd":
			err := osCD(args)
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "pwd":
			dir, err := osPWD()
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println(dir)
			}
		case "echo":
			osECHO(args)
		case "kill":
			err := osKILL(args)
			if err != nil {
				fmt.Println("Error:", err)
			}
		case "ps":
			err := osPS()
			if err != nil {
				fmt.Println("Error:", err)
			}
		default:
			fmt.Printf("%s: command not found\n", args[0])
		}
	}
}
