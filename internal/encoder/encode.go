package encoder

import (
	"bytes"
	"sort"
)

type StringDeEncoder interface {
	Encode(string) ([]byte, error)
	Decode([]byte) string
	Bytes() int
}

type IntDeEncoder interface {
	// Encode transform the value
	// when the value has been used before
	Encode(int) ([]byte, bool)
	Decode([]byte) int
	Bytes() int

	// Transform can be used to transform
	// values that haven't been seen before.
	// Assuming we stored
	// 2 as 0000 0010 and
	// 4 as 0000 0011
	// When we are looking for 3 (a value that isn't stored)
	// we return both values and the callee can decide which
	// one to use.
	// The first value is the biggest value smaller than v.
	// The second value is the smallest value bigger than v.
	// When v is smaller than all values the first value will
	// be nil. Analog to this the last value will be nil when
	// v is bigger than all values.
	Transform(v int) ([]byte, []byte)
}

type intDeEncoder struct {
	values   []int
	encoding [][]byte
	width    int
}

func (de intDeEncoder) Bytes() int { return de.width }

func (de intDeEncoder) Encode(v int) ([]byte, bool) {
	idx := sort.SearchInts(de.values, v)
	if idx != len(de.values) && de.values[idx] != v {
		return nil, false
	}
	return de.encoding[idx], true
}

func (de intDeEncoder) Transform(v int) ([]byte, []byte) {
	idx := sort.SearchInts(de.values, v)
	if idx == len(de.values) {
		return de.encoding[idx-1], nil
	}
	if idx == 0 {
		return nil, de.encoding[idx]
	}
	return de.encoding[idx-1], de.encoding[idx]
}

func (de intDeEncoder) Decode(buf []byte) int {
	// Decode succeeds
	idx := sort.Search(len(de.encoding), func(i int) bool { return bytes.Compare(buf, de.encoding[i]) < 0 })
	return de.values[idx]
}

func intToBytes(v, bytes int) []byte {
	buf := make([]byte, bytes)
	for i := bytes - 1; i >= 0; i-- {
		buf[i] = byte(v & 0xff)
		v >>= 8
	}
	return buf
}

func NewIntDeEncoder(ints []int) IntDeEncoder {
	intset := map[int]struct{}{}

	for _, v := range ints {
		intset[v] = struct{}{}
	}

	//calc that!
	bytes := len(intset) << 8

	uniques := make([]int, len(intset))
	for k := range intset {
		uniques = append(uniques, k)
	}

	sort.Ints(uniques)

	res := intDeEncoder{
		width: bytes,
	}

	for i, v := range uniques {
		res.values = append(res.values, v)
		res.encoding = append(res.encoding, intToBytes(i, bytes))
	}
	return res
}
