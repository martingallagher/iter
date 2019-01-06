// Package iterutil provides common iterator functionality.
package iterutil

import (
	"unicode"

	"github.com/martingallagher/iter"
	"github.com/martingallagher/iter/internal/runes"
)

// Fields returns a new string iterator emitting values between each instance
// of one or more consecutive white space runes.
func Fields(s string) *iter.FuncIter {
	return iter.NewFuncString(s, unicode.IsSpace)
}

// Lines returns a new string iterator emitting values between newlines.
func Lines(s string) *iter.FuncIter {
	return iter.NewFuncString(s, runes.IsNewline)
}

// Numbers returns a new string iterator emitting numeric values.
func Numbers(s string) *iter.FuncIter {
	return iter.NewFuncString(s, unicode.IsNumber)
}

// Words returns a new string iterator naively emitting words.
func Words(s string) *iter.FuncIter {
	return iter.NewFuncString(s, runes.IsNotLN)
}
