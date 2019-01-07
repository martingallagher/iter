package iter

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unicode"

	"github.com/martingallagher/iter/internal/function"
	"github.com/martingallagher/iter/internal/runes"
)

const space = " "

var (
	funcs      = []func(rune) bool{unicode.IsSpace, runes.IsNewline, runes.IsNotLN}
	spaceBytes = []byte(space)
	seperators = []string{"", space, ", ", ",\n"}
	tests      = []string{
		"",                 // Empty string
		"    aa a      a ", // Spaces
		"Hi",
		" Hello World ",
		"   !Κάρολος       ...Δαρβίνος    123",
		`   An American monkey,
	after getting drunk on brandy,
	would never touch it again,
	and thus is much wiser than most men.   `,
		"The bigger the interface, the weaker the abstraction.",
		"A\tB\tC",
		`This royal throne of kings, this scepter'd isle,
This earth of majesty, this seat of Mars,
This other Eden, demi-paradise,
This fortress built by Nature for herself
Against infection and the hand of war,
This happy breed of men, this little world,
This precious stone set in the silver sea,
Which serves it in the office of a wall,
Or as a moat defensive to a house,
Against the envy of less happier lands,
This blessed plot, this earth, this realm, this England`,
		benchString,
	}
)

func TestBytes(t *testing.T) {
	for _, sep := range seperators {
		sep := []byte(sep)

		for _, v := range tests {
			b := []byte(v)
			expected := removeEmptyBytes(bytes.Split(b, sep))
			values := make([][]byte, 0, len(expected))
			iter := New(b, sep)

			for iter.Next() {
				values = append(values, iter.Bytes())
			}

			err := testBytes(expected, values)

			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestBytesFuncChan(t *testing.T) {
	b := []byte(tests[4])
	expected := bytes.FieldsFunc(b, runes.IsNotLN)
	values := make([][]byte, 0, len(expected))

	for v := range NewFunc(b, runes.IsNotLN).Chan() {
		values = append(values, v)
	}

	err := testBytes(expected, values)

	if err != nil {
		t.Error(err)
	}
}

func TestBytesEmitAll(t *testing.T) {
	for i, sep := range seperators {
		sep := []byte(sep)

		for j, v := range tests {
			b := []byte(v)
			expected := bytes.Split(b, sep)
			values := make([][]byte, 0, len(expected))
			iter := New(b, sep)
			iter.EmitAll()

			for iter.Next() {
				values = append(values, iter.Bytes())
			}

			err := testBytes(expected, values)

			if err != nil {
				t.Errorf("Test %d:%d: %s", i, j, err)
			}
		}
	}
}

func TestBytesFunc(t *testing.T) {
	t.Run("BytesFunc", func(t *testing.T) {
		for _, f := range funcs {
			t.Run(function.Name(f), func(t *testing.T) {
				for _, v := range tests {
					b := []byte(v)
					expected := bytes.FieldsFunc(b, f)
					values := make([][]byte, 0, len(expected))
					iter := NewFunc(b, f)

					for iter.Next() {
						values = append(values, iter.Bytes())
					}

					err := testBytes(expected, values)

					if err != nil {
						t.Error(err)
					}
				}
			})
		}
	})
}

func TestString(t *testing.T) {
	for _, sep := range seperators {
		for _, v := range tests {
			expected := removeEmptyStrings(strings.Split(v, sep))
			values := make([]string, 0, len(expected))
			iter := NewString(v, sep)

			for iter.Next() {
				values = append(values, iter.String())
			}

			err := testStrings(expected, values)

			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestStringEmitAll(t *testing.T) {
	for i, sep := range seperators {
		for j, v := range tests {
			expected := strings.Split(v, sep)
			values := make([]string, 0, len(expected))
			iter := NewString(v, sep)
			iter.EmitAll()

			for iter.Next() {
				values = append(values, iter.String())
			}

			err := testStrings(expected, values)

			if err != nil {
				t.Errorf("Test %d:%d: %s", i, j, err)
			}
		}
	}
}

func TestStringFunc(t *testing.T) {
	t.Run("StringFunc", func(t *testing.T) {
		for _, f := range funcs {
			t.Run(function.Name(f), func(t *testing.T) {
				for _, v := range tests {
					expected := strings.FieldsFunc(v, f)
					values := make([]string, 0, len(expected))
					iter := NewFuncString(v, f)

					for iter.Next() {
						values = append(values, iter.String())
					}

					err := testStrings(expected, values)

					if err != nil {
						t.Error(err)
					}
				}
			})
		}
	})
}

func testBytes(expected, got [][]byte) error {
	if !reflect.DeepEqual(expected, got) {
		return fmt.Errorf("byte slice failed; expected %v (len=%d), got %v (len=%d)",
			expected, len(expected), got, len(got))
	}

	return nil
}

func testStrings(expected, got []string) error {
	if !reflect.DeepEqual(expected, got) {
		return fmt.Errorf("string slice failed; expected %v (len=%d), got %v (len=%d)",
			expected, len(expected), got, len(got))
	}

	return nil
}

func removeEmptyBytes(s [][]byte) [][]byte {
	out := s[:0]

	for _, v := range s {
		if len(v) > 0 {
			out = append(out, v)
		}
	}

	return out
}

func removeEmptyStrings(s []string) []string {
	out := s[:0]

	for _, v := range s {
		if v != "" {
			out = append(out, v)
		}
	}

	return out
}
