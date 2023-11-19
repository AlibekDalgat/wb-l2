package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

/*
=== Утилита grep ===

Реализовать утилиту фильтрации (man grep)

Поддержать флаги:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", печатать номер строки

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

// функция которая ищет подстроку в строке разными способами
func containStr(row, input string, fixed bool) bool {
	// если нужно искать точное совпадение
	if fixed {
		// вопспользоваться функцией из библиотеки strings ищущая точное вхождение
		if strings.Contains(row, input) {
			return true
		}
	} else { // если нужно искать с помощью по регулярному выражению
		// воспользоваться библиотекой regexp
		regex := regexp.MustCompile(input)
		matches := regex.FindString(row)
		if matches != "" {
			return true
		}
	}
	return false
}

func main() {
	// чтение ввода того что нужно искать
	input := os.Args[len(os.Args)-1]
	// вводные строки
	lines := make([]string, 0)
	// результирующий слайс
	res := make([]string, 0)
	// чтение ввода
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// флаги командной строки
	after := flag.Int("A", 0, "печатать +N строк после совпадения")
	before := flag.Int("B", 0, "печатать +N строк до совпадения")
	context := flag.Int("C", 0, "печатать ±N строк вокруг совпадения")
	count := flag.Bool("c", false, "количество строк")
	ignore := flag.Bool("i", false, "игнорировать регистр")
	invert := flag.Bool("v", false, "вместо совпадения исключать")
	fixed := flag.Bool("F", false, "точное совпадение со строкой, не паттерн")
	lineNum := flag.Bool("n", false, "печатать номер строки")
	flag.Parse()
	// счётчик строк которые выведутся
	k := 0
	// мапа указывающая на то, была ли добавлена строка в результирующий слайс
	met := make(map[int]struct{})
	for i, row := range lines {
		// строка без приведения к нижнему регистру
		originRow := row
		// если есть ключ игнорирование регистра привести строку к нижнему регистру по которому будет происходить поиск
		if *ignore {
			row = strings.ToLower(row)
			input = strings.ToLower(input)
		}
		// если есть ключ печати номера строки добавить номер строки в начале строки
		if *lineNum {
			originRow = strconv.Itoa(i+1) + ":" + originRow
		}
		// произвести поиск
		matched := containStr(row, input, *fixed)
		// если поиск положительный и ключа инверта нет или поиск отрицательный и инверт есть, начать процесс добавления строки в результат
		if matched && !*invert || !matched && *invert {
			// если есть ключи поиска до совпадения
			if *before > 0 || *context > 0 {
				var b int
				if *before > 0 {
					b = *before
				} else {
					b = *context
				}
				// добавить в результирующий слайс строки до совпадения
				for j := i - b; j < i; j++ {
					// если строки до совпадения уже встречались не добавлять
					if _, ok := met[j]; !ok && j >= 0 {
						res = append(res, lines[j])
						met[j] = struct{}{}
						k++
					}
				}
			}
			// если данная строка ещё добавлялась в результируюущю не добавлять
			if _, ok := met[i]; !ok {
				res = append(res, originRow)
				met[i] = struct{}{}
				k++
			}
			// если есть ключи поиска после совпадения
			if *after > 0 || *context > 0 {
				var a int
				if *after > 0 {
					a = *after
				} else {
					a = *context
				}
				// добавить в результирующий слайс строки после совпадения
				for j := i + 1; j <= i+a && j < len(lines); j++ {
					// если строки после совпадения уже встречались не добавлять
					if _, ok := met[j]; !ok && j >= 0 {
						res = append(res, lines[j])
						met[j] = struct{}{}
						k++
					}
				}
			}
		}
	}
	// если есть ключ вывода количества строк
	if *count {
		fmt.Println(k)
	} else {
		// иначе вывести результат
		for _, row := range res {
			fmt.Println(row)
		}
	}
}
