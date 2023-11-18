package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
=== Утилита sort ===

Отсортировать строки (man sort)
Основное

Поддержать ключи

-k — указание колонки для сортировки
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительное

Поддержать ключи

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учётом суффиксов

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

func main() {
	// офлаги командной строки
	column := flag.Int("k", 0, "указание колонки для сортировки.")
	numeric := flag.Bool("n", false, "сортировать по числовому значению.")
	reverse := flag.Bool("r", false, "сортировать в обратном порядке.")
	unique := flag.Bool("u", false, "не выводить повторяющиеся строки.")
	month := flag.Bool("M", false, "сортировать по названию месяца")
	ignoreTrailing := flag.Bool("b", false, "игнорирование хвостовых пробелом")
	isSorted := flag.Bool("c", false, "проверять отсортированы ли данные")
	flag.Parse()

	// чтение файла
	rows, err := readFile()
	if err != nil {
		fmt.Printf("Ошибка при чтении файла: %v\n", err)
		return
	}

	// функция сравнения для сортировки
	comparator := func(i, j int) bool {
		// сравнение двух конкретных строк
		valueI, valueJ := rows[i], rows[j]
		// дополнительные переменные для контроля того являются ли конкретная колонка числовым значением
		var (
			isNum            bool
			numValI, numValJ int
		)

		// если есть ключ игнорирования хвостовых пробелом - удалить все лишние про
		if *ignoreTrailing {
			valueI = strings.TrimRight(valueI, " ")
			valueJ = strings.TrimRight(valueJ, " ")
		}

		// если ключуказания колонки больше нуля
		if *column > 0 {
			// вынести определённую колонку
			fieldsI := strings.Fields(valueI)
			fieldsJ := strings.Fields(valueJ)

			// если номер колонки под=одит под количество колонок строки
			if *column <= len(fieldsI) && *column <= len(fieldsJ) {
				// приравнять всю строку этой колонке
				valueI = fieldsI[*column-1]
				valueJ = fieldsJ[*column-1]
			} else {
				// иначе сделать её пустой
				valueI = ""
				valueJ = ""
			}
		}
		// если есть ключ сортировки по ключевому значению
		if *numeric {
			numI, errI := strconv.Atoi(valueI)
			numJ, errJ := strconv.Atoi(valueJ)

			// если привести строку в число получается, определеить флаг isNum, указывающий на то что сравнивать нужно как числовые значения
			if errI == nil && errJ == nil {
				isNum = true
				numValI = numI
				numValJ = numJ
			}
		}

		// если есть ключсортировки по месяцам
		if *month {
			fieldsI := strings.Fields(valueI)
			fieldsJ := strings.Fields(valueJ)
			// сортировать по номеру месяца
			if *reverse {
				return monthes[fieldsI[5]] > monthes[fieldsJ[5]]
			}
			return monthes[fieldsI[5]] < monthes[fieldsJ[5]]
		}

		// возврашение результата сравнения
		// если есть ключ ревёрса сравнивать наоборот
		if *reverse {
			// если флаг isNum положительный, сравнить числовые значения
			if isNum {
				return numValI > numValJ
			}
			return valueI > valueJ
		}
		if isNum {
			return numValI < numValJ
		}
		return valueI < valueJ
	}

	// если есть ключпроверки сортирован ли данные
	if *isSorted {
		if !sort.SliceIsSorted(rows, comparator) {
			fmt.Println("Данные не отсортированы.")
		} else {
			fmt.Println("Данные отсортированы.")
		}
		return
	}

	// сортировка
	sort.SliceStable(rows, comparator)

	// если есть ключвывода уникальных значений
	if *unique {
		// запусить функцию удаления дублирующих значений
		rows = removeDuplicates(rows)
	}

	// Сохранение отсортированных строк в файл
	err = writeLines(rows)
	if err != nil {
		fmt.Printf("Ошибка при записи в файл: %v\n", err)
		return
	}
}

// функция чтения строк из файла и возвращает их в виде слайса
func readFile() ([]string, error) {
	file, err := os.Open("input.txt")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

// функция записывания строк в файл
func writeLines(lines []string) error {
	file, err := os.Create("output.txt")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}

	return writer.Flush()
}

// функция удаления дублирующих значений
func removeDuplicates(lines []string) []string {
	// мапа содержащая информациюю о уже встреченных строк
	uniqueLines := make(map[string]struct{})
	var result []string

	for _, line := range lines {
		// если строка ещё не встречалась (то если отсутвуюет в маппе) добавить в результирующий слайс
		if _, ok := uniqueLines[line]; !ok {
			uniqueLines[line] = struct{}{}
			result = append(result, line)
		}
	}

	return result
}

// мапа содержащая информацию о номерном значении названий месяцов
var monthes = map[string]int{
	"янв": 0,
	"фев": 1,
	"мар": 2,
	"апр": 3,
	"мая": 4,
	"июн": 5,
	"июл": 6,
	"авг": 7,
	"сен": 8,
	"окт": 9,
	"ноя": 10,
	"дек": 11,
}
