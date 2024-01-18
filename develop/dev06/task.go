package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

/*
=== Утилита cut ===

Принимает STDIN, разбивает по разделителю (TAB) на колонки, выводит запрошенные

Поддержать флаги:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	// вводные строки
	lines := make([]string, 0)
	// чтение ввода
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		text := scanner.Text()
		// чтение прекратится когда встретится пустая строчка
		if strings.TrimSpace(text) == "" {
			break
		}
		lines = append(lines, text)
	}
	// флаги командной строки
	fields := flag.Int("f", 0, "выбрать поля (колонки)")
	delimetr := flag.String("d", "\t", "использовать другой разделитель")
	separated := flag.Bool("s", false, "только строки с разделителем")
	flag.Parse()

	// анализ каждой введёной строки
	for _, line := range lines {
		// разбиение строки на слоайс
		splited := strings.Split(line, *delimetr)
		// если длина слайса равна одному и нету ключа о обязательном присутствии разделителя, напечатать строку
		if len(splited) == 1 && !*separated {
			fmt.Println(line)
		} else if len(splited) == 1 && *separated { //если длина слайса равна одному и есть ключ о обязательном присутствии разделителя, пропустить
			continue
		} else if *fields > len(splited) { // если номер поля больше длины слайса напечатать пустую строку
			fmt.Println()
		} else { // иначе вывести нужный элемент поля
			fmt.Println(splited[*fields-1])
		}
	}
}
