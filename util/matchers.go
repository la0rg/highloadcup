package util

import "unicode"

func OnlyLetters(s string) bool {
	return unicode.IsLetter(rune(s[0]))
}
