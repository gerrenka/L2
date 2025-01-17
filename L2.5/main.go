package main

import (
	"sort"
	"strings"
)

// sortString возвращает отсортированную строку
func sortString(s string) string {
	// Преобразуем строку в слайс рун для корректной работы с UTF-8
	runes := []rune(s)
	// Сортируем руны
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	// Возвращаем отсортированную строку
	return string(runes)
}

// FindAnagrams находит все множества анаграмм в словаре
func FindAnagrams(dictionary *[]string) *map[string][]string {
	// Создаем мапу для хранения промежуточных результатов
	// Ключ - отсортированные буквы слова, значение - слайс анаграмм
	tempMap := make(map[string][]string)
	
	// Создаем мапу для хранения первого встретившегося слова для каждого множества
	firstWord := make(map[string]string)
	
	// Обрабатываем каждое слово из словаря
	for _, word := range *dictionary {
		// Приводим слово к нижнему регистру
		word = strings.ToLower(word)
		
		// Получаем отсортированные буквы слова
		sorted := sortString(word)
		
		// Если это первое слово с таким набором букв, сохраняем его
		if _, exists := firstWord[sorted]; !exists {
			firstWord[sorted] = word
		}
		
		// Добавляем слово в соответствующий слайс
		tempMap[sorted] = append(tempMap[sorted], word)
	}
	
	// Создаем результирующую мапу
	result := make(map[string][]string)
	
	// Заполняем результирующую мапу
	for sorted, words := range tempMap {
		// Пропускаем множества из одного элемента
		if len(words) < 2 {
			continue
		}
		
		// Удаляем дубликаты
		uniqueWords := make([]string, 0)
		seen := make(map[string]bool)
		
		for _, word := range words {
			if !seen[word] {
				uniqueWords = append(uniqueWords, word)
				seen[word] = true
			}
		}
		
		// Сортируем слова
		sort.Strings(uniqueWords)
		
		// Добавляем в результат, используя первое встретившееся слово как ключ
		result[firstWord[sorted]] = uniqueWords
	}
	
	return &result
}