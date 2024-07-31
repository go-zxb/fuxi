package pkg

import (
	"strings"
)

func InitialLetter(word string) string {
	if len(word) > 1 {
		return strings.ToUpper(word[0:1]) + word[1:]
	}
	return strings.ToUpper(word)
}
