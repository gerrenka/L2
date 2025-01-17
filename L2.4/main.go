package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type SortOptions struct {
	column        int
	numeric       bool
	reverse       bool
	unique        bool
	byMonth       bool
	ignoreSpaces  bool
	check         bool
	humanNumeric  bool
}

type Line struct {
	content string
	parts   []string
}

func main() {
	// Определяем флаги командной строки
	column := flag.Int("k", 1, "column to sort by (1-based)")
	numeric := flag.Bool("n", false, "sort by numeric value")
	reverse := flag.Bool("r", false, "sort in reverse order")
	unique := flag.Bool("u", false, "output unique lines only")
	byMonth := flag.Bool("M", false, "sort by month name")
	ignoreSpaces := flag.Bool("b", false, "ignore trailing spaces")
	check := flag.Bool("c", false, "check if sorted")
	humanNumeric := flag.Bool("h", false, "sort by human readable numbers (2K, 1M, etc)")

	flag.Parse()

	// Получаем имя входного файла из аргументов
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: sort [options] input_file")
		os.Exit(1)
	}

	options := SortOptions{
		column:        *column - 1, // Преобразуем в 0-based индекс
		numeric:       *numeric,
		reverse:       *reverse,
		unique:        *unique,
		byMonth:       *byMonth,
		ignoreSpaces:  *ignoreSpaces,
		check:         *check,
		humanNumeric:  *humanNumeric,
	}

	// Читаем и сортируем файл
	lines, err := readLines(args[0], options)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	// Сортируем строки
	sorted := sortLines(lines, options)

	// Проверяем сортировку если установлен флаг -c
	if options.check {
		if !isSorted(sorted, options) {
			fmt.Println("Input is not sorted")
			os.Exit(1)
		}
		fmt.Println("Input is sorted")
		return
	}

	// Выводим результат
	printLines(sorted, options)
}

func readLines(filename string, options SortOptions) ([]Line, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []Line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		content := scanner.Text()
		if options.ignoreSpaces {
			content = strings.TrimRight(content, " \t")
		}
		parts := strings.Fields(content)
		lines = append(lines, Line{content: content, parts: parts})
	}

	return lines, scanner.Err()
}

func sortLines(lines []Line, options SortOptions) []Line {
	sorted := make([]Line, len(lines))
	copy(sorted, lines)

	sort.Slice(sorted, func(i, j int) bool {
		return compareLine(sorted[i], sorted[j], options)
	})

	if options.unique {
		sorted = uniqueLines(sorted)
	}

	return sorted
}

func compareLine(a, b Line, options SortOptions) bool {
	// Получаем значения для сравнения
	valA := getValue(a, options)
	valB := getValue(b, options)

	result := false

	switch {
	case options.numeric:
		numA, errA := strconv.ParseFloat(valA, 64)
		numB, errB := strconv.ParseFloat(valB, 64)
		if errA == nil && errB == nil {
			result = numA < numB
		} else {
			result = valA < valB
		}
	case options.byMonth:
		monthA := parseMonth(valA)
		monthB := parseMonth(valB)
		result = monthA < monthB
	case options.humanNumeric:
		numA := parseHumanNumber(valA)
		numB := parseHumanNumber(valB)
		result = numA < numB
	default:
		result = valA < valB
	}

	if options.reverse {
		return !result
	}
	return result
}

func getValue(line Line, options SortOptions) string {
	if len(line.parts) > options.column {
		return line.parts[options.column]
	}
	return line.content
}

func parseMonth(s string) time.Month {
	// Пытаемся распарсить название месяца
	months := map[string]time.Month{
		"jan": time.January,
		"feb": time.February,
		"mar": time.March,
		"apr": time.April,
		"may": time.May,
		"jun": time.June,
		"jul": time.July,
		"aug": time.August,
		"sep": time.September,
		"oct": time.October,
		"nov": time.November,
		"dec": time.December,
	}

	s = strings.ToLower(s)
	if month, ok := months[s[:3]]; ok {
		return month
	}
	return 0
}

func parseHumanNumber(s string) float64 {
	s = strings.TrimSpace(strings.ToUpper(s))
	multipliers := map[string]float64{
		"K": 1000,
		"M": 1000000,
		"G": 1000000000,
		"T": 1000000000000,
	}

	for suffix, multiplier := range multipliers {
		if strings.HasSuffix(s, suffix) {
			value, err := strconv.ParseFloat(s[:len(s)-1], 64)
			if err == nil {
				return value * multiplier
			}
		}
	}

	value, _ := strconv.ParseFloat(s, 64)
	return value
}

func uniqueLines(lines []Line) []Line {
	seen := make(map[string]bool)
	var result []Line

	for _, line := range lines {
		if !seen[line.content] {
			seen[line.content] = true
			result = append(result, line)
		}
	}

	return result
}

func isSorted(lines []Line, options SortOptions) bool {
	for i := 1; i < len(lines); i++ {
		if compareLine(lines[i], lines[i-1], options) {
			return false
		}
	}
	return true
}

func printLines(lines []Line, options SortOptions) {
	for _, line := range lines {
		fmt.Println(line.content)
	}
}