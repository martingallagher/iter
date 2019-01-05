package iter

import (
	"reflect"
	"unsafe"
)

// String is a substring iterator.
type String struct {
	b *Bytes
}

// NewString returns a new string iterator.
func NewString(s, needle string) *String {
	return &String{NewBytes(unsafeBytes(s), unsafeBytes(needle))}
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

// EmitAll emits all (including consecutive) matches.
func (s *String) EmitAll() {
	s.b.emitAll = true
}

// Next iterates to the next value; returning false if the iterator is exhausted.
func (s *String) Next() bool {
	return s.b.Next()
}

// ForEach executes a provided function once for each string value.
func (s *String) ForEach(f func(string)) {
	s.b.Reset()

	for s.b.Next() {
		f(unsafeString(s.b.v))
	}
}

// Chan returns a channel for receiving iterator values.
func (s *String) Chan() <-chan string {
	values := make(chan string)

	go func() {
		s.ForEach(func(s string) {
			values <- s
		})

		close(values)
	}()

	return values
}

// String returns the current string value.
func (s *String) String() string {
	return unsafeString(s.b.v)
}

// Reset resets the iterator start position.
func (s *String) Reset() {
	s.b.Reset()
}

// StringFunc is a rune based iterator, iterating through a string
// rune-by-rune applying a given field function.
type StringFunc struct {
	b *BytesFunc
}

// NewStringFunc returns a new StringFunc iterator.
func NewStringFunc(s string, f func(rune) bool) *StringFunc {
	return &StringFunc{NewBytesFunc(unsafeBytes(s), f)}
}

// Next iterates to the next value; returning false if the iterator is exhausted.
func (s *StringFunc) Next() bool {
	return s.b.Next()
}

// ForEach executes a provided function once for each string value.
func (s *StringFunc) ForEach(f func(string)) {
	s.b.Reset()

	for s.b.Next() {
		f(unsafeString(s.b.v))
	}
}

// Chan returns a channel for receiving iterator values.
func (s *StringFunc) Chan() <-chan string {
	values := make(chan string)

	go func() {
		s.ForEach(func(s string) {
			values <- s
		})

		close(values)
	}()

	return values
}

// String returns the current string value.
func (s *StringFunc) String() string {
	return unsafeString(s.b.v)
}

// Reset resets the iterator start position.
func (s *StringFunc) Reset() {
	s.b.Reset()
}
