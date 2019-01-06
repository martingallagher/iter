package iterutil_test

import (
	"reflect"
	"strings"
	"testing"
	"unicode"

	"github.com/martingallagher/iter"
	"github.com/martingallagher/iter/internal/function"
	"github.com/martingallagher/iter/internal/runes"
	"github.com/martingallagher/iter/iterutil"
)

func TestIterUtil(t *testing.T) {
	tests := []string{
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
	}
	iterators := []struct {
		iter func(string) *iter.FuncIter
		f    func(rune) bool
	}{
		{iterutil.Fields, unicode.IsSpace},
		{iterutil.Lines, runes.IsNewline},
		{iterutil.Numbers, unicode.IsNumber},
		{iterutil.Words, runes.IsNotLN},
	}

	for _, v := range iterators {
		name := function.Name(v.f)

		t.Run(name, func(t *testing.T) {
			for _, s := range tests {
				expected := strings.FieldsFunc(s, v.f)
				iter := v.iter(s)
				l := len(expected)
				values := make([]string, 0, l)

				for iter.Next() {
					values = append(values, iter.String())
				}

				if !reflect.DeepEqual(expected, values) {
					t.Errorf("%s iterator failed; expected %v (len=%d), got %v (len=%d)",
						name, expected, l, values, len(values))
				}
			}
		})
	}
}
