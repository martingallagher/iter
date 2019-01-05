package iter

import (
	"bytes"
	"unicode/utf8"
)

// Bytes is a bytes substring iterator.
type Bytes struct {
	b           []byte
	needle      []byte
	v           []byte
	l, n, start int
	emitAll     bool
}

// NewBytes returns a new bytes iterator.
func NewBytes(b, needle []byte) *Bytes {
	return &Bytes{b: b, needle: needle, l: len(b), n: len(needle)}
}

// EmitAll emits all (including consecutive) matches.
func (b *Bytes) EmitAll() {
	b.emitAll = true
}

// Next iterates to the next value; returning false if the iterator is exhausted.
func (b *Bytes) Next() bool {
	if b.start > b.l {
		return false
	}

	if b.n == 0 {
		// Empty needle; replicate bytes.Split and split after each UTF-8 sequence
		return b.nextRune()
	}

	if b.l == 0 {
		if !b.emitAll {
			return false
		}

		b.v = b.b
		b.start = b.l + 1

		return true
	}

	match := false

	for i := b.start; i < b.l; i++ {
		end := i + b.n

		if end > b.l {
			break
		}

		if !bytes.Equal(b.b[i:end], b.needle) {
			// Within a value range; continue reading
			if match {
				continue
			}

			b.start = i
			match = true

			continue
		}

		if !b.emitAll && !match {
			b.start = i

			continue
		}

		// Emit current value
		b.v = b.b[b.start:i]
		b.start = i + b.n

		return true
	}

	if !match {
		if !b.emitAll || !bytes.Equal(b.b[b.l-b.n:], b.needle) {
			return false
		}

		b.v = b.b[b.start:]
		b.start = b.l + 1

		return true
	}

	// Emit remaining value
	b.v = b.b[b.start:]
	b.start += b.l

	return true
}

func (b *Bytes) nextRune() bool {
	if b.l == 0 || b.start == b.l {
		return false
	}

	size := 1

	if b.b[b.start] >= utf8.RuneSelf {
		_, size = utf8.DecodeRune(b.b[b.start:])
	}

	end := b.start + size
	b.v = b.b[b.start:end]
	b.start = end

	return true
}

// ForEach executes a provided function once for each byte slice value.
func (b *Bytes) ForEach(f func(b []byte)) {
	b.Reset()

	for b.Next() {
		f(b.Bytes())
	}
}

// Chan returns a channel for receiving iterator values.
func (b *Bytes) Chan() <-chan []byte {
	values := make(chan []byte)

	go func() {
		b.ForEach(func(b []byte) {
			values <- b
		})

		close(values)
	}()

	return values
}

// Bytes returns the current byte slice value.
func (b *Bytes) Bytes() []byte {
	return b.v
}

// Reset resets the iterator start position.
func (b *Bytes) Reset() {
	b.start = 0
}

// BytesFunc is a rune based iterator, iterating through a byte slice
// rune-by-rune applying a given field function.
type BytesFunc struct {
	f        func(rune) bool
	b        []byte
	v        []byte
	l, start int
}

// NewBytesFunc returns a new BytesFunc iterator.
func NewBytesFunc(b []byte, f func(rune) bool) *BytesFunc {
	return &BytesFunc{f: f, b: b, l: len(b)}
}

// Next iterates to the next value; returning false if the iterator is exhausted.
func (b *BytesFunc) Next() bool {
	if b.start >= b.l {
		return false
	}

	i := b.start
	match := false

	for i < b.l {
		size := 1
		r := rune(b.b[i])

		if r >= utf8.RuneSelf {
			r, size = utf8.DecodeRune(b.b[i:])
		}

		if b.f(r) {
			if match {
				b.v = b.b[b.start:i]
				b.start = i + size

				return true
			}

			b.start = i + size
			match = false
		} else if !match {
			b.start = i
			match = true
		}

		i += size
	}

	if match && b.start < b.l {
		b.v = b.b[b.start:]
		b.start += b.l - b.start

		return true
	}

	return false
}

// ForEach executes a provided function once for each byte slice value.
func (b *BytesFunc) ForEach(f func(b []byte)) {
	b.Reset()

	for b.Next() {
		f(b.Bytes())
	}
}

// Chan returns a channel for receiving iterator values.
func (b *BytesFunc) Chan() <-chan []byte {
	values := make(chan []byte)

	go func() {
		b.ForEach(func(b []byte) {
			values <- b
		})

		close(values)
	}()

	return values
}

// Bytes returns the current byte slice value.
func (b *BytesFunc) Bytes() []byte {
	return b.v
}

// Reset resets the iterator start position.
func (b *BytesFunc) Reset() {
	b.start = 0
}
