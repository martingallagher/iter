// Package iter implements non-allocating rune, byte and string iterators.
package iter

// Segmenter defines the segmentation interface.
type Segmenter interface {
	Segment(b []byte) ([]byte, []byte)
}

type SegmenterFunc func(b []byte) ([]byte, []byte)

func (fn SegmenterFunc) Segment(b []byte) ([]byte, []byte) {
	return fn(b)
}

func Names(b []byte) ([]byte, []byte) {
	return nil, nil
}

var _ SegmenterFunc = Names
