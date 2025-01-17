package main

import (
    "errors"
    "fmt"
    "strconv"
    "strings"
    "unicode"
)

func unpackString(s string) (string, error) {
    if s == "" {
        return "", nil
    }

    var result strings.Builder
    runes := []rune(s)
    isEscaped := false
    
    if len(runes) > 0 && unicode.IsDigit(runes[0]) {
        return "", errors.New("некорректная строка: начинается с цифры")
    }

    for i := 0; i < len(runes); i++ {
        currentRune := runes[i]

        if currentRune == '\\' && !isEscaped {
            isEscaped = true
            continue
        }

        if unicode.IsDigit(currentRune) && !isEscaped {
            if i == 0 {
                return "", errors.New("некорректная строка: начинается с цифры")
            }

            numStr := string(currentRune)
            for j := i + 1; j < len(runes) && unicode.IsDigit(runes[j]); j++ {
                numStr += string(runes[j])
                i = j
            }

            count, err := strconv.Atoi(numStr)
            if err != nil {
                return "", errors.New("ошибка при конвертации числа")
            }

            if count <= 0 {
                return "", errors.New("некорректное количество повторений: число должно быть положительным")
            }

            prevChar := runes[i-len(numStr)]
            if count > 1 {
                result.WriteString(strings.Repeat(string(prevChar), count-1))
            }
        } else {
            result.WriteRune(currentRune)
            isEscaped = false
        }
    }

    return result.String(), nil
}

func main() {
    examples := []string{
        "a4bc2d5e",
        "abcd",
        "45",
        "",
        "qwe\\4\\5",
        "qwe\\45",
        "qwe\\\\5",
        "a0",
    }

    for _, example := range examples {
        result, err := unpackString(example)
        if err != nil {
            fmt.Printf("Ошибка для строки %q: %v\n", example, err)
        } else {
            fmt.Printf("Результат для строки %q: %q\n", example, result)
        }
    }
}
	