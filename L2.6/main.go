package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

type GrepConfig struct {
	After      int
	Before     int
	Context    int
	Count      bool
	IgnoreCase bool
	Invert     bool
	Fixed      bool
	LineNum    bool
	Pattern    string
	InputFile  string
}

func main() {
	config := parseFlags()
	
	if err := grep(config); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func parseFlags() GrepConfig {
	after := flag.Int("A", 0, "print N lines after match")
	before := flag.Int("B", 0, "print N lines before match")
	context := flag.Int("C", 0, "print N lines around match")
	count := flag.Bool("c", false, "print count of matching lines")
	ignoreCase := flag.Bool("i", false, "ignore case")
	invert := flag.Bool("v", false, "invert match")
	fixed := flag.Bool("F", false, "fixed string match")
	lineNum := flag.Bool("n", false, "print line numbers")

	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "pattern is required")
		os.Exit(1)
	}

	inputFile := ""
	if len(args) > 1 {
		inputFile = args[1]
	}

	return GrepConfig{
		After:      *after,
		Before:     *before,
		Context:    *context,
		Count:      *count,
		IgnoreCase: *ignoreCase,
		Invert:     *invert,
		Fixed:      *fixed,
		LineNum:    *lineNum,
		Pattern:    args[0],
		InputFile:  inputFile,
	}
}

func grep(config GrepConfig) error {
	var reader io.Reader
	if config.InputFile != "" {
		file, err := os.Open(config.InputFile)
		if err != nil {
			return err
		}
		defer file.Close()
		reader = file
	} else {
		reader = os.Stdin
	}

	scanner := bufio.NewScanner(reader)
	
	// Буфер для хранения предыдущих строк (для опции -B)
	var beforeLines []string
	var matchCount int
	var lineNum int
	var afterCount int
	
	// Определяем максимальное количество строк до совпадения
	maxBefore := config.Before
	if config.Context > maxBefore {
		maxBefore = config.Context
	}

	// Подготавливаем паттерн
	pattern := config.Pattern
	if config.IgnoreCase {
		pattern = strings.ToLower(pattern)
	}

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		var matches bool

		if config.Fixed {
			if config.IgnoreCase {
				matches = strings.Contains(strings.ToLower(line), pattern)
			} else {
				matches = strings.Contains(line, pattern)
			}
		} else {
			if config.IgnoreCase {
				matches = strings.Contains(strings.ToLower(line), pattern)
			} else {
				matches = strings.Contains(line, pattern)
			}
		}

		if config.Invert {
			matches = !matches
		}

		if matches {
			matchCount++
			
			if config.Count {
				continue
			}

			// Печатаем предыдущие строки
			if maxBefore > 0 {
				for i, bLine := range beforeLines {
					if config.LineNum {
						fmt.Printf("%d:", lineNum-len(beforeLines)+i)
					}
					fmt.Println(bLine)
				}
			}

			// Печатаем текущую строку
			if config.LineNum {
				fmt.Printf("%d:", lineNum)
			}
			fmt.Println(line)
			
			afterCount = config.After
			if config.Context > afterCount {
				afterCount = config.Context
			}
		} else {
			if afterCount > 0 {
				if config.LineNum {
					fmt.Printf("%d:", lineNum)
				}
				fmt.Println(line)
				afterCount--
			}
		}

		// Обновляем буфер предыдущих строк
		if maxBefore > 0 {
			beforeLines = append(beforeLines, line)
			if len(beforeLines) > maxBefore {
				beforeLines = beforeLines[1:]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	if config.Count {
		fmt.Println(matchCount)
	}

	return nil
}