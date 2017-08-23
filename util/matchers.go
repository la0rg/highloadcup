package util

import "unicode"

func OnlyLetters(s string) bool {
	return unicode.IsLetter(rune(s[0]))
}

func IsGender(s string) bool {
	return s == "m" || s == "f"
}
