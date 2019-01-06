// Package iterutil provides common iterator functionality.
package iterutil

import (
	"unicode"

	"github.com/martingallagher/iter"
)

// Fields returns a new string iterator emitting values between each instance
// of one or more consecutive white space runes.
func Fields(s string) *iter.FuncIter {
	return iter.NewFuncString(s, unicode.IsSpace)
}

func isNewline(r rune) bool {
	return r == '\n' || r == '\r'
}

// Lines returns a new string iterator emitting values between newlines.
func Lines(s string) *iter.FuncIter {
	return iter.NewFuncString(s, isNewline)
}

// Numbers returns a new string iterator emitting numeric values.
func Numbers(s string) *iter.FuncIter {
	return iter.NewFuncString(s, unicode.IsNumber)
}

func isNotLN(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsNumber(r)
}

// Words returns a new string iterator naively emitting words.
func Words(s string) *iter.FuncIter {
	return iter.NewFuncString(s, isNotLN)
}
