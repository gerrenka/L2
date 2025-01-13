package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	// Убираем символ новой строки в конце
	input = strings.TrimSuffix(input, "\n")

	// Дополнительно уберём пробелы по краям (если нужно)
	trimmed := strings.TrimSpace(input)

	// Если строка пустая
	if trimmed == "" {
		fmt.Println("")
		return
	}

	// Если вся оставшаяся строка — это число (проверка через Atoi)
	if _, err := strconv.Atoi(trimmed); err == nil {
		fmt.Println("некорректная строка")
		return
	}

	runes := []rune(input)
	var result []rune

	for i := 0; i < len(runes); i++ {
		// Проверяем, что текущий символ — это backslash,
		// и что впереди есть ещё символ, который — цифра
		if runes[i] == '\\' && i+1 < len(runes) && unicode.IsDigit(runes[i+1]) {
			// Добавляем в результат только цифру (без backslash)
			result = append(result, runes[i+1])
			// Пропускаем следующий символ (т. е. сдвигаем индекс ещё на 1)
			i++
		} else {
			// Иначе добавляем текущий символ как есть
			result = append(result, runes[i])
		}
	}

	fmt.Println(string(result))
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
