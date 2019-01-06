package runes

import "unicode"

// IsNewline naively checks if the rune is a newline character.
func IsNewline(r rune) bool {
	return r == '\n' || r == '\r'
}

// IsLN checks if the rune is classified as a Unicode letter or number.
func IsLN(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsNumber(r)
}

// IsNotLN checks if the rune is not classified as a Unicode letter or number.
func IsNotLN(r rune) bool {
	return !IsLN(r)
}
