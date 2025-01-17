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
    
    // Проверка первого символа
    if len(runes) > 0 && unicode.IsDigit(runes[0]) {
        return "", errors.New("некорректная строка: начинается с цифры")
    }

    for i := 0; i < len(runes); i++ {
        currentRune := runes[i]

        // Обработка escape-последовательности
        if currentRune == '\\' && !isEscaped {
            isEscaped = true
            continue
        }

        // Если текущий символ - цифра
        if unicode.IsDigit(currentRune) && !isEscaped {
            if i == 0 {
                return "", errors.New("некорректная строка: начинается с цифры")
            }

            // Собираем все последующие цифры
            numStr := string(currentRune)
            for j := i + 1; j < len(runes) && unicode.IsDigit(runes[j]); j++ {
                numStr += string(runes[j])
                i = j
            }

            // Конвертируем строку в число
            count, err := strconv.Atoi(numStr)
            if err != nil {
                return "", errors.New("ошибка при конвертации числа")
            }

            // Проверяем корректность числа повторений
            if count <= 0 {
                return "", errors.New("некорректное количество повторений: число должно быть положительным")
            }

            // Повторяем предыдущий символ
            prevChar := runes[i-len(numStr)]
            // Важно: теперь вычитаем 1 из count, так как сам символ уже записан
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
    // Примеры использования
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
	

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"os"
// 	"strconv"
// 	"unicode"
// )

// func main() {
// 	reader := bufio.NewReader(os.Stdin)
// 	stroka, _ := reader.ReadString('\n')
// 	stroka = stroka[:len(stroka)-1] // Убираем символ новой строки
// 	runes := []rune(stroka)

// 	if stroka == "" { // Проверяем, является ли строка пустой
// 		fmt.Println("")
// 		return
// 	} else if _, err := strconv.Atoi(stroka); err == nil { // Проверяем, является ли строка числом
// 		fmt.Println("некорректная строка")
// 		return
// 	} else {
// 		for i := 0; i < len(runes); i++ {
// 			if unicode.IsDigit(runes[i]) {
// 				num, err := strconv.Atoi(string(runes[i])) // Преобразуем символ в число
// 				if err != nil {
// 					return
// 				}
// 				if i > 0 { // Проверяем, чтобы i-1 не выходил за пределы
// 					if unicode.IsDigit(runes[i-1]) {
						
// 					}else{
// 					for j := 0; j < num-1; j++ {
// 						if runes[i-1] == '\\' {
// 							if unicode.IsDigit(runes[i]){
// 								fmt.Printf("%c", runes[i])
// 								break
// 							}
// 					}else if unicode.IsDigit(runes[i-1]) {
// 						fmt.Printf("%c", runes[i])
// 					}else {
// 						fmt.Printf("%c", runes[i-1])
					
// 					}}
// 				}
// 			} else if unicode.IsLetter(runes[i]) {
// 				fmt.Printf("%c", runes[i])
// 			} else if runes[i] == '\\' {
// 				continue
				
// 			}
// 		}
// 	}
// }

// }
