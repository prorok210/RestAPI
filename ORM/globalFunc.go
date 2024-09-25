package main

import (
	"strings"
	"unicode"
)

func Contains(slice []string, input string) bool {
	for _, str := range slice {
		if strings.Contains(str, input) {
			return true
		}
	}
	return false
}

// The function takes a string and returns it with the 1st character in large case
func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}
