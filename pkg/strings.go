package pkg

import (
	"strings"
	"unicode"
)

func InitialLetter(word string) string {
	if len(word) > 1 {
		return strings.ToUpper(word[0:1]) + word[1:]
	}
	return strings.ToUpper(word)
}

func InitialLetterToLower(word string) string {
	if len(word) > 1 {
		return strings.ToLower(word[0:1]) + word[1:]
	}
	return strings.ToLower(word)
}

func CamelToSnake(camelStr string) string {
	var result strings.Builder
	for i, char := range camelStr {
		if unicode.IsUpper(char) {
			if i > 0 {
				result.WriteRune('_')
			}
			result.WriteRune(unicode.ToLower(char))
		} else {
			result.WriteRune(char)
		}
	}
	return result.String()
}
