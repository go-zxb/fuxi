package pkg

import (
	"unicode"
	"unicode/utf8"
)

func HasChinese(s string) bool {
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		if unicode.Is(unicode.Han, r) {
			return true
		}
		i += size
	}
	return false
}
