package main

import "testing"

func TestUnpackString(t *testing.T) {
    tests := []struct {
        input    string
        expected string
        hasError bool
    }{
        {"a4bc2d5e", "aaaabccddddde", false},
        {"abcd", "abcd", false},
        {"45", "", true},
        {"", "", false},
        {"qwe\\4\\5", "qwe45", false},
        {"qwe\\45", "qwe44444", false},
        {"qwe\\\\5", "qwe\\\\\\\\\\", false},
        {"a2b3c4", "aabbbcccc", false},
        {"a", "a", false},
        {"a0", "", true},
        {"привет2", "приветт", false},
    }

    for _, test := range tests {
        result, err := unpackString(test.input)
        
        // Проверка на наличие ошибки
        if test.hasError && err == nil {
            t.Errorf("Ожидалась ошибка для входной строки: %q", test.input)
            continue
        }
        
        if !test.hasError && err != nil {
            t.Errorf("Неожиданная ошибка для входной строки: %q: %v", test.input, err)
            continue
        }

        // Проверка результата
        if !test.hasError && result != test.expected {
            t.Errorf("Для входной строки %q ожидалось %q, получено %q", 
                     test.input, test.expected, result)
        }
    }
}