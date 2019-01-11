// Package iter implements non-allocating rune, byte and string iterators.
package iter

import (
	"bytes"
	"reflect"
	"unicode"
	"unicode/utf8"
	"unsafe"

	"github.com/martingallagher/iter/internal/runes"
)

// Iterator defines the iterator interface.
type Iterator interface {
	Next() bool
	Value() Value
}

// Valuer defines the bytes/string value interface.
type Valuer interface {
	Bytes() []byte
	String() string
}

// Value represents an iterator value.
type Value []byte

// Bytes returns the bytes value.
func (v Value) Bytes() []byte {
	return v
}

// String returns the string value.
func (v Value) String() string {
	return unsafeString(v)
}

// Iter implements a substring iterator.
type Iter struct {
	b           []byte
	needle      []byte
	v           []byte
	l, n, start int
	emitAll     bool
}

func unsafeBytes(s string) []byte {
	h := *(*reflect.StringHeader)(unsafe.Pointer(&s))

	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: h.Data,
		Len:  h.Len,
		Cap:  h.Len,
	}))
}

func unsafeString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// New returns a new iterator for the given haystack and needle bytes.
func New(b, needle []byte) *Iter {
	return &Iter{b: b, needle: needle, l: len(b), n: len(needle)}
}

// NewString returns a new iterator for the given haystack and needle strings.
func NewString(s, needle string) *Iter {
	return New(unsafeBytes(s), unsafeBytes(needle))
}

// EmitAll emits all (including consecutive) matches.
func (i *Iter) EmitAll() {
	i.emitAll = true
}

// Next iterates to the next value; returning false
// if the iterator is exhausted.
func (i *Iter) Next() bool {
	if i.start > i.l {
		return false
	}

	if i.n == 0 {
		// Empty needle; replicate bytes.Split and split after each UTF-8 sequence
		return i.nextRune()
	}

	if i.l == 0 {
		if !i.emitAll {
			return false
		}

		i.v = i.b
		i.start = i.l + 1

		return true
	}

	match := false

	for j := i.start; j < i.l; j++ {
		end := j + i.n

		if end > i.l {
			break
		}

		if !bytes.Equal(i.b[j:end], i.needle) {
			// Within a value range; continue reading
			if match {
				continue
			}

			i.start = j
			match = true

			continue
		}

		if !i.emitAll && !match {
			i.start = j

			continue
		}

		// Emit current value
		i.v = i.b[i.start:j]
		i.start = j + i.n

		return true
	}

	if !match {
		if !i.emitAll || !bytes.Equal(i.b[i.l-i.n:], i.needle) {
			return false
		}

		i.v = i.b[i.start:]
		i.start = i.l + 1

		return true
	}

	// Emit remaining value
	i.v = i.b[i.start:]
	i.start += i.l

	return true
}

func (i *Iter) nextRune() bool {
	if i.l == 0 || i.start == i.l {
		return false
	}

	size := 1

	if i.b[i.start] >= utf8.RuneSelf {
		_, size = utf8.DecodeRune(i.b[i.start:])
	}

	end := i.start + size
	i.v = i.b[i.start:end]
	i.start = end

	return true
}

// Value returns the current value.
func (i *Iter) Value() Value {
	return i.v
}

// Reset resets the iterator start position.
func (i *Iter) Reset() {
	i.start = 0
}

// FuncIter is a rune based iterator, iterating
// rune-by-rune applying a given function.
type FuncIter struct {
	fn       func(rune) bool
	b        []byte
	v        []byte
	l, start int
}

// NewFunc returns a new rune function iterator
// for the given bytes input.
func NewFunc(b []byte, fn func(rune) bool) *FuncIter {
	return &FuncIter{fn: fn, b: b, l: len(b)}
}

// NewFuncString returns a new rune function iterator
// for the given string input.
func NewFuncString(s string, f func(rune) bool) *FuncIter {
	return NewFunc(unsafeBytes(s), f)
}

// Next iterates to the next value; returning false if the iterator is exhausted.
func (f *FuncIter) Next() bool {
	if f.start >= f.l {
		return false
	}

	i := f.start
	match := false

	for i < f.l {
		size := 1
		r := rune(f.b[i])

		if r >= utf8.RuneSelf {
			r, size = utf8.DecodeRune(f.b[i:])
		}

		if f.fn(r) {
			if match {
				f.v = f.b[f.start:i]
				f.start = i + size

				return true
			}

			f.start = i + size
			match = false
		} else if !match {
			f.start = i
			match = true
		}

		i += size
	}

	if match && f.start < f.l {
		f.v = f.b[f.start:]
		f.start += f.l - f.start

		return true
	}

	return false
}

// Value returns the current value.
func (f *FuncIter) Value() Value {
	return f.v
}

// Reset resets the iterator start position.
func (f *FuncIter) Reset() {
	f.start = 0
}

// Chan returns a channel for receiving iterator values.
func Chan(i Iterator) <-chan Value {
	c := make(chan Value)

	go func() {
		for i.Next() {
			c <- i.Value()
		}

		close(c)
	}()

	return c
}

// ForEach provides for-each semantics for iterators.
func ForEach(i Iterator, fn func(Value)) {
	for i.Next() {
		fn(i.Value())
	}
}

// Fields returns a new string iterator emitting values between each instance
// of one or more consecutive white space runes.
func Fields(s string) *FuncIter {
	return NewFuncString(s, unicode.IsSpace)
}

// Lines returns a new string iterator emitting values between newlines.
func Lines(s string) *FuncIter {
	return NewFuncString(s, runes.IsNewline)
}

// Numbers returns a new string iterator emitting numeric values.
func Numbers(s string) *FuncIter {
	return NewFuncString(s, unicode.IsNumber)
}

// Words returns a new string iterator naively emitting words.
func Words(s string) *FuncIter {
	return NewFuncString(s, runes.IsNotLN)
}
