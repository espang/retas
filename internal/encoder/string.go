package encoder

import (
	"bytes"
	"sort"
)

type StringIterator interface {
	Next() bool
	Value() string
}

type StringDeEncoder interface {
	Encode(string) ([]byte, bool)
	Decode([]byte) string
	Bytes() int

	Transform(string) ([]byte, []byte)
}

type stringDeEncoder struct {
	values   []string
	encoding [][]byte
	width    int
}

func (de stringDeEncoder) Bytes() int { return de.width }

func (de stringDeEncoder) Encode(v string) ([]byte, bool) {
	idx := sort.SearchStrings(de.values, v)
	if idx == len(de.values) {
		return nil, false
	}
	if de.values[idx] != v {
		return nil, false
	}
	return de.encoding[idx], true
}

func (de stringDeEncoder) Transform(v string) ([]byte, []byte) {
	idx := sort.SearchStrings(de.values, v)
	if idx == len(de.values) {
		return de.encoding[idx-1], nil
	}
	if idx == 0 {
		return nil, de.encoding[idx]
	}
	return de.encoding[idx-1], de.encoding[idx]
}

func (de stringDeEncoder) Decode(buf []byte) string {
	// Decode succeeds
	idx := sort.Search(len(de.encoding), func(i int) bool { return bytes.Compare(buf, de.encoding[i]) < 0 })
	return de.values[idx-1]
}

func NewStringDeEncoder(iter StringIterator) StringDeEncoder {
	stringset := map[string]struct{}{}

	for iter.Next() {
		stringset[iter.Value()] = struct{}{}
	}

	bytes := bytesByUniques(len(stringset))
	uniques := make([]string, 0, len(stringset))
	for k := range stringset {
		uniques = append(uniques, k)
	}

	sort.Strings(uniques)
	res := stringDeEncoder{
		width: bytes,
	}

	for i, v := range uniques {
		res.values = append(res.values, v)
		res.encoding = append(res.encoding, intToBytes(i, bytes))
	}
	return res
}
