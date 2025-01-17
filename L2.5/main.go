package main

import (
	"sort"
	"strings"
)

func sortString(s string) string {
	runes := []rune(s)
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	return string(runes)
}

func FindAnagrams(dictionary *[]string) *map[string][]string {
	tempMap := make(map[string][]string)
	
	firstWord := make(map[string]string)
	
	for _, word := range *dictionary {

		word = strings.ToLower(word)
		sorted := sortString(word)

		if _, exists := firstWord[sorted]; !exists {
			firstWord[sorted] = word
		}
		
		tempMap[sorted] = append(tempMap[sorted], word)
	}
	

	result := make(map[string][]string)

	for sorted, words := range tempMap {
		if len(words) < 2 {
			continue
		}

		uniqueWords := make([]string, 0)
		seen := make(map[string]bool)
		
		for _, word := range words {
			if !seen[word] {
				uniqueWords = append(uniqueWords, word)
				seen[word] = true
			}
		}

		sort.Strings(uniqueWords)

		result[firstWord[sorted]] = uniqueWords
	}
	
	return &result
}