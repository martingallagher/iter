// Package iter implements non-allocating rune, byte and string iterators.
package iter

import (
	"bytes"
	"reflect"
	"unicode/utf8"
	"unsafe"
)

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

// ForEach executes a provided function once for each byte slice value.
func (i *Iter) ForEach(fn func(b []byte)) {
	i.Reset()

	for i.Next() {
		fn(i.v)
	}
}

// ForEachString executes a provided function once for each string value.
func (i *Iter) ForEachString(fn func(s string)) {
	i.Reset()

	for i.Next() {
		fn(i.String())
	}
}

// Chan returns a channel for receiving iterator byte values.
func (i *Iter) Chan() <-chan []byte {
	values := make(chan []byte)

	go func() {
		i.ForEach(func(b []byte) {
			values <- b
		})

		close(values)
	}()

	return values
}

// ChanString returns a channel for receiving iterator string values.
func (i *Iter) ChanString() <-chan string {
	values := make(chan string)

	go func() {
		i.ForEachString(func(s string) {
			values <- s
		})

		close(values)
	}()

	return values
}

// Bytes returns the current byte slice value.
func (i *Iter) Bytes() []byte {
	return i.v
}

// String returns the current string value.
func (i *Iter) String() string {
	return unsafeString(i.v)
}

// Reset resets the iterator start position.
func (i *Iter) Reset() {
	i.start = 0
}

// FuncIter is a rune based iterator, iterating
// rune-by-rune applying a given field function.
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

// ForEach executes a provided function once for each byte slice value.
func (f *FuncIter) ForEach(fn func(b []byte)) {
	f.Reset()

	for f.Next() {
		fn(f.v)
	}
}

// ForEachString executes a provided function once for each string value.
func (f *FuncIter) ForEachString(fn func(s string)) {
	f.Reset()

	for f.Next() {
		fn(f.String())
	}
}

// Chan returns a channel for receiving iterator values.
func (f *FuncIter) Chan() <-chan []byte {
	values := make(chan []byte)

	go func() {
		f.ForEach(func(b []byte) {
			values <- b
		})

		close(values)
	}()

	return values
}

// ChanString returns a channel for receiving iterator string values.
func (f *FuncIter) ChanString() <-chan string {
	values := make(chan string)

	go func() {
		f.ForEachString(func(s string) {
			values <- s
		})

		close(values)
	}()

	return values
}

// Bytes returns the current byte slice value.
func (f *FuncIter) Bytes() []byte {
	return f.v
}

// String returns the current string value.
func (f *FuncIter) String() string {
	return unsafeString(f.v)
}

// Reset resets the iterator start position.
func (f *FuncIter) Reset() {
	f.start = 0
}
