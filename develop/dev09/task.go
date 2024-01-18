package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// функция реализующая работу wget
func myWget(url string) error {
	// отправления get запроса по ссылке
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// если статус запроса не равна 200 вернуть ошибку
	if resp.StatusCode != http.StatusOK {
		return errors.New("failed to fetch page: " + resp.Status)
	}
	// возвращение последнего элемента url адресса для наименования результирующего файла
	fileName := path.Base(url)
	// создание результирующего файла
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	// копирование тела запроса в результирующий файл
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("downloaded: %s\n", fileName)
	return nil
}

func main() {
	// если нет аргументов для команды вернуть ошибку
	if len(os.Args) < 2 {
		fmt.Println("go run main.go <URL>")
		return
	}
	// вынимание адреса
	url := os.Args[1]
	// возвращение ошибки или её отсутствия от работы функции  myWget(url)
	err := myWget(url)
	if err != nil {
		fmt.Printf("error downloading: %s\n", err)
	}
}
