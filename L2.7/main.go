package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	fields    string
	delimiter string
	separated bool
}

// parseFields преобразует строку с номерами полей в слайс индексов
func parseFields(fieldsStr string) ([]int, error) {
	if fieldsStr == "" {
		return nil, fmt.Errorf("fields cannot be empty")
	}

	var fields []int
	parts := strings.Split(fieldsStr, ",")

	for _, part := range parts {
		// Проверяем на диапазон (например, 1-3)
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) != 2 {
				return nil, fmt.Errorf("invalid range format: %s", part)
			}

			start, err := strconv.Atoi(rangeParts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid start of range: %s", rangeParts[0])
			}

			end, err := strconv.Atoi(rangeParts[1])
			if err != nil {
				return nil, fmt.Errorf("invalid end of range: %s", rangeParts[1])
			}

			if start > end {
				return nil, fmt.Errorf("invalid range: start > end")
			}

			for i := start; i <= end; i++ {
				fields = append(fields, i)
			}
		} else {
			// Одиночное число
			field, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid field number: %s", part)
			}
			fields = append(fields, field)
		}
	}

	return fields, nil
}

// processLine обрабатывает одну строку ввода
func processLine(line string, config Config, fields []int) (string, error) {
	// Если флаг -s установлен и в строке нет разделителя, пропускаем строку
	if config.separated && !strings.Contains(line, config.delimiter) {
		return "", nil
	}

	parts := strings.Split(line, config.delimiter)
	var result []string

	for _, field := range fields {
		// Поля в cut нумеруются с 1
		idx := field - 1
		if idx >= 0 && idx < len(parts) {
			result = append(result, parts[idx])
		}
	}

	return strings.Join(result, config.delimiter), nil
}

func main() {
	config := Config{}

	// Определение флагов командной строки
	flag.StringVar(&config.fields, "f", "", "select only these fields")
	flag.StringVar(&config.delimiter, "d", "\t", "use delimiter instead of TAB")
	flag.BoolVar(&config.separated, "s", false, "only lines containing delimiter")

	flag.Parse()

	// Проверяем обязательный параметр -f
	if config.fields == "" {
		fmt.Fprintln(os.Stderr, "Error: -f (fields) parameter is required")
		flag.Usage()
		os.Exit(1)
	}

	// Парсим поля
	fields, err := parseFields(config.fields)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing fields: %v\n", err)
		os.Exit(1)
	}

	// Читаем ввод построчно
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		result, err := processLine(line, config, fields)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error processing line: %v\n", err)
			continue
		}
		if result != "" {
			fmt.Println(result)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}
}