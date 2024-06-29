package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	input  string = "input.txt"
	output string = "output.txt"
)

func main() {
	err := mathExec(input, output)
	if err != nil {
		log.Fatal(err)
	}
}

// mathExec - основная функция задачи.
func mathExec(input, output string) error {
	// Читаем файл с входящими данными и делим его построчно.
	fileIn, err := os.ReadFile(input)
	if err != nil {
		return fmt.Errorf("cannot read input file: %w", err)
	}
	lines := bytes.Split(fileIn, []byte("\n"))

	// Удаляем файл вывода, если он уже существует.
	if _, err := os.Stat(output); !os.IsNotExist(err) {
		os.Remove(output)
	}

	// Создаем новый файл вывода.
	fileOut, err := os.OpenFile(output, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("cannot create output file: %w", err)
	}
	defer fileOut.Close()

	// Создаем буфер для результатов вычислений
	buf := bufio.NewWriter(fileOut)
	defer buf.Flush()

	// Задаем паттерн регулярного выражения.
	re := regexp.MustCompile(`([0-9\.]+)([\s\+\-\*\/]+)([0-9\.]+)([=\s]+)(\?)`)

	// Построчно прогоняем входящие данные.
	for _, l := range lines {
		line := re.FindAllStringSubmatch(string(l), -1)
		if len(line) == 0 {
			continue
		}

		// Цикл на случай, если в одной	строке несколько выражений.
		for i := 0; i < len(line); i++ {
			// Убираем возможные пробелы
			line[i][2] = strings.TrimSpace(line[i][2])
			line[i][4] = strings.TrimSpace(line[i][4])
			// Считаем выражение
			res, err := calculator(line[i][1], line[i][2], line[i][3])
			if err != nil {
				log.Printf("incorrect expression: %s", err)
				continue
			}
			// Записываем результат вычислений вместо знака вопроса.
			line[i][5] = fmt.Sprintf("%g", res)
			line[i] = append(line[i], "\n")
			// Собираем строку и записываем ее в буфер вывода.
			str := strings.Join(line[i][1:], "")
			_, err = buf.WriteString(str)
			if err != nil {
				log.Printf("cannot write result string: %s", err)
				continue
			}
		}
	}
	return nil
}

// calculator - выполняет вычисление арифметического выражения.
func calculator(opOne, operator, opTwo string) (float64, error) {
	numOne, err := strconv.ParseFloat(opOne, 64)
	if err != nil {
		return 0, fmt.Errorf("incorrect first operand")
	}
	numTwo, err := strconv.ParseFloat(opTwo, 64)
	if err != nil {
		return 0, fmt.Errorf("incorrect second operand")
	}
	if numTwo == 0 {
		return 0, fmt.Errorf("division by zero")
	}

	var res float64
	switch operator {
	case "+":
		res = numOne + numTwo
	case "-":
		res = numOne - numTwo
	case "*":
		res = numOne * numTwo
	case "/":
		res = numOne / numTwo
	default:
		return 0, fmt.Errorf("unknown operator")
	}
	return res, nil
}
