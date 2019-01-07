[![Build Status](https://travis-ci.org/martingallagher/iter.svg)](https://travis-ci.org/martingallagher/iter) [![GoDoc](https://godoc.org/github.com/martingallagher/iter?status.svg)](https://godoc.org/github.com/martingallagher/iter) [![Go Report Card](https://goreportcard.com/badge/github.com/martingallagher/iter)](https://goreportcard.com/report/github.com/martingallagher/iter) [![license](https://img.shields.io/github/license/martingallagher/iter.svg)](https://github.com/martingallagher/iter/blob/master/LICENSE)

# iter: Byte and String Iterators

`iter` provides low overhead (zero allocation where possible) iterators for strings and byte slices, fulfilling both `(bytes|strings).Split` and `(bytes|strings).FieldsFunc` functionality and additional helper functions.

## Examples

Word count:

```go
s := "My long string..."

func isNotLN(r rune) {
  return !unicode.IsLetter(r) && !unicode.IsNumber(r)
}

// Standard library
func stdCountOccurrences(s, word string) int {
  count := 0
  words := strings.FieldsFunc(s, isNotLN)

  for v := range words {
    if strings.EqualFold(word, v) {
      count++
    }
  }

  return count
}

// iter package
func iterCountOccurrences(s, word string) int {
  count := 0
  iter := iter.NewFuncString(s, isNotLN)

  for iter.Next() {
    if strings.EqualFold(word, iter.String()) {
      count++
    }
  }

  return count
}
```

## Benchmarks

    goos: linux
    goarch: amd64
    pkg: github.com/martingallagher/iter
    pkg: github.com/martingallagher/iter
    BenchmarkBytes-8           	  100000	     21530 ns/op	       0 B/op	       0 allocs/op
    BenchmarkString-8          	  100000	     21568 ns/op	       0 B/op	       0 allocs/op
    BenchmarkBytesEmitAll-8    	  100000	     21730 ns/op	       0 B/op	       0 allocs/op
    BenchmarkStringEmitAll-8   	  100000	     21678 ns/op	       0 B/op	       0 allocs/op
    BenchmarkBytesFunc-8       	  100000	     19977 ns/op	       0 B/op	       0 allocs/op
    BenchmarkStringFunc-8      	  100000	     20018 ns/op	       0 B/op	       0 allocs/op
    BenchmarkStdStringsMap-8   	   50000	     22758 ns/op	   11264 B/op	       4 allocs/op
    PASS
    ok  	github.com/martingallagher/iter	15.395s
