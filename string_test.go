package iter

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"unicode"
)

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

func TestStringForEach(t *testing.T) {
	for _, sep := range seperators {
		for _, v := range tests {
			expected := removeEmptyStrings(strings.Split(v, sep))
			l := len(expected)
			values := make([]string, 0, l)
			iter := NewString(v, sep)

			iter.ForEach(func(s string) {
				values = append(values, s)
			})

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
			t.Run(funcName(f), func(t *testing.T) {
				for _, v := range tests {
					expected := strings.FieldsFunc(v, f)
					values := make([]string, 0, len(expected))
					iter := NewStringFunc(v, f)

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

func TestStringHelpers(t *testing.T) {
	iterators := []struct {
		iter func(string) *StringFunc
		f    func(rune) bool
	}{
		{Fields, unicode.IsSpace},
		{Lines, isNewline},
		{Numbers, unicode.IsNumber},
		{Words, isNotLN},
	}

	for _, v := range iterators {
		name := funcName(v.f)

		t.Run(name, func(t *testing.T) {
			for _, s := range tests {
				expected := strings.FieldsFunc(s, v.f)
				iter := v.iter(s)
				l := len(expected)
				values := make([]string, 0, l)

				for iter.Next() {
					values = append(values, iter.String())
				}

				if !equalStrings(expected, values) {
					t.Errorf("%s iterator failed; expected %v (len=%d), got %v (len=%d)",
						name, expected, l, values, len(values))
				}
			}
		})
	}
}

type IterReseter interface {
	Next() bool
	String() string
	Reset()
}

func testIterReseter(iter IterReseter) error {
	var expected []string

	for iter.Next() {
		expected = append(expected, iter.String())
	}

	if iter.Next() {
		return errors.New("unexpected iteration")
	}

	iter.Reset()

	l := len(expected)
	values := make([]string, 0, l)

	for iter.Next() {
		values = append(values, iter.String())
	}

	if !equalStrings(expected, values) {
		return fmt.Errorf("expected %v (len=%d), got %v (len=%d)",
			expected, l, values, len(values))
	}

	return nil
}

func TestStringReset(t *testing.T) {
	tests := []struct {
		name string
		iter IterReseter
	}{
		{"String", NewString(benchString, space)},
		{"StringFunc", Fields(benchString)},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			err := testIterReseter(v.iter)

			if err != nil {
				t.Error(err)
			}
		})
	}
}

func TestWordsCount(t *testing.T) {
	const needle = "Go"

	expected := countStd(benchString, needle)
	count := count(benchString, needle)

	t.Logf("Found %d occurences of %q", count, needle)

	if count != expected {
		t.Errorf("expected %d, got %d", expected, count)
	}
}

func testStrings(expected, got []string) error {
	if !equalStrings(expected, got) {
		return fmt.Errorf("string slice failed; expected %v (len=%d), got %v (len=%d)",
			expected, len(expected), got, len(got))
	}

	return nil
}

func count(s, needle string) int {
	count := 0
	iter := Words(s)

	for iter.Next() {
		if iter.String() == needle {
			count++
		}
	}

	return count
}

func countForEach(s, needle string) int {
	count := 0

	Words(s).ForEach(func(v string) {
		if v == needle {
			count++
		}
	})

	return count
}

func countChan(s, needle string) int {
	count := 0

	for v := range Words(s).Chan() {
		if v == needle {
			count++
		}
	}

	return count
}

func countStd(s, needle string) int {
	words := strings.FieldsFunc(s, isNotLN)
	count := 0

	for _, v := range words {
		if v == needle {
			count++
		}
	}

	return count
}
