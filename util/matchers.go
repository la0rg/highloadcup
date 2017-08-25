package util

import (
	"strings"
	"unicode"
)

func OnlyLetters(s string) bool {
	return unicode.IsLetter(rune(s[0]))
}

func IsGender(s string) bool {
	return s == "m" || s == "f"
}

func ContainsNull(bytes []byte) bool {
	return strings.Contains(string(bytes), ": null")
}
