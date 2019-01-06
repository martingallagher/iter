package iter

import (
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unicode"
)

const space = " "

var (
	funcs      = []func(rune) bool{unicode.IsSpace, isNewline, isNotLN}
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

func TestBytesForEach(t *testing.T) {
	for _, sep := range seperators {
		sep := []byte(sep)

		for _, v := range tests {
			b := []byte(v)
			expected := removeEmptyBytes(bytes.Split(b, sep))
			l := len(expected)
			values := make([][]byte, 0, l)
			iter := New(b, sep)

			iter.ForEach(func(b []byte) {
				values = append(values, b)
			})

			err := testBytes(expected, values)

			if err != nil {
				t.Error(err)
			}
		}
	}
}

func TestBytesFuncChan(t *testing.T) {
	b := []byte(tests[4])
	expected := bytes.FieldsFunc(b, isNotLN)
	values := make([][]byte, 0, len(expected))

	for v := range NewFunc(b, isNotLN).Chan() {
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
			t.Run(funcName(f), func(t *testing.T) {
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

func testBytes(expected, got [][]byte) error {
	if !equalBytes(expected, got) {
		return fmt.Errorf("byte slice failed; expected %v (len=%d), got %v (len=%d)",
			expected, len(expected), got, len(got))
	}

	return nil
}

func funcName(f interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	i := strings.IndexByte(name, '/')

	if i != -1 {
		name = name[i+1:]
	}

	return name
}

func equalBytes(a, b [][]byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !bytes.Equal(a[i], b[i]) {
			return false
		}
	}

	return true
}

func equalStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func removeEmptyBytes(s [][]byte) [][]byte {
	f := s[:0]

	for _, v := range s {
		if len(v) > 0 {
			f = append(f, v)
		}
	}

	return f
}

func removeEmptyStrings(s []string) []string {
	f := s[:0]

	for _, v := range s {
		if v != "" {
			f = append(f, v)
		}
	}

	return f
}
