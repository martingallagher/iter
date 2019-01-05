package iter

import "unicode"

// Fields returns a new string iterator emitting values between each instance
// of one or more consecutive white space runes.
func Fields(s string) *StringFunc {
	return NewStringFunc(s, unicode.IsSpace)
}

func isNewline(r rune) bool {
	return r == '\n' || r == '\r'
}

// Lines returns a new string iterator emitting values between newlines.
func Lines(s string) *StringFunc {
	return NewStringFunc(s, isNewline)
}

// Numbers returns a new string iterator emitting numeric values.
func Numbers(s string) *StringFunc {
	return NewStringFunc(s, unicode.IsNumber)
}

func isNotLN(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsNumber(r)
}

// Words returns a new string iterator naively emitting words.
func Words(s string) *StringFunc {
	return NewStringFunc(s, isNotLN)
}
