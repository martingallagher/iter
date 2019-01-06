package iter

import (
	"fmt"
	"strings"
	"testing"
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

			iter.ForEachString(func(s string) {
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

func testStrings(expected, got []string) error {
	if !equalStrings(expected, got) {
		return fmt.Errorf("string slice failed; expected %v (len=%d), got %v (len=%d)",
			expected, len(expected), got, len(got))
	}

	return nil
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
