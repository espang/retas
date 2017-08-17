package encoder

import (
	"bytes"
	"fmt"
	"math"
	"sort"
	"strconv"
	"testing"
)

func arange(start, end, stepsize int) []int {
	var res []int
	for i := start; i < end; i += stepsize {
		res = append(res, i)
	}
	return res
}

func stringify(ints []int) (func(int) string, []string) {
	if len(ints) == 0 {
		panic("expect at least 1 value")
	}
	sort.Ints(ints)

	max := ints[len(ints)-1]
	length := math.Log10(float64(max))

	pattern := `%0` + strconv.Itoa(int(length)+1) + `d`

	t := func(v int) string {
		return fmt.Sprintf(pattern, v)
	}

	var res []string
	for _, i := range ints {
		res = append(res, t(i))
	}
	return t, res
}

var testCases = []struct {
	name  string
	input []int
	bytes int
	good  map[int][]byte
	bad   map[int][][]byte
}{
	{
		"two ints",
		[]int{1, 10},
		1,
		map[int][]byte{
			1:  []byte{0},
			10: []byte{1},
		},
		map[int][][]byte{
			0:  [][]byte{nil, []byte{0}},
			2:  [][]byte{[]byte{0}, []byte{1}},
			11: [][]byte{[]byte{1}, nil},
		},
	},
	{
		"ints 1, 2, ... 256",
		arange(1, 257, 1),
		1,
		map[int][]byte{
			1:   []byte{0},
			100: []byte{99},
			256: []byte{255},
		},
		map[int][][]byte{
			0:   [][]byte{nil, []byte{0}},
			257: [][]byte{[]byte{255}, nil},
		},
	},
	{
		"ints 1, 3, ... 513, 515",
		arange(1, 516, 2),
		2,
		map[int][]byte{
			1:   []byte{0, 0},
			513: []byte{1, 0},
			515: []byte{1, 1},
		},
		map[int][][]byte{
			0:   [][]byte{nil, []byte{0, 0}},
			2:   [][]byte{[]byte{0, 0}, []byte{0, 1}},
			516: [][]byte{[]byte{1, 1}, nil},
		},
	},
}

func TestIntDeEncode(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			de := NewIntDeEncoder(tc.input)

			if got := de.Bytes(); got != tc.bytes {
				t.Errorf("got %v; want %v", got, tc.bytes)
			}

			for k, v := range tc.good {
				got, ok := de.Encode(k)
				if !ok {
					t.Error("got !ok; want ok")
				}
				if !bytes.Equal(v, got) {
					t.Errorf("got %v; want %v", got, v)
				}

				val := de.Decode(got)
				if val != k {
					t.Errorf("got %v; want %v", val, k)
				}
			}

			for k, v := range tc.bad {
				_, ok := de.Encode(k)
				if ok {
					t.Error("got ok; want !ok")
				}

				lower, upper := de.Transform(k)
				if !bytes.Equal(lower, v[0]) {
					t.Errorf("got %v; want %v", lower, v[0])
				}

				if !bytes.Equal(upper, v[1]) {
					t.Errorf("got %v; want %v", upper, v[1])
				}
			}
		})
	}
}

func TestIntToByte(t *testing.T) {
	testCases := []struct {
		name   string
		input  int
		bytes  int
		output []byte
	}{
		{
			"single byte 1001 0110",
			150,
			1,
			[]byte{150},
		},
		{
			"single byte 1111 1111",
			255,
			1,
			[]byte{255},
		},
		{
			"single byte 0000 0000",
			0,
			1,
			[]byte{0},
		},
		{
			"two bytes 0001 0000 0000 1000",
			1<<3 + 1<<12,
			2,
			[]byte{16, 8},
		},
		{
			"three bytes 1000 0011 0001 0000 0000 1000",
			4104 + 1<<16 + 1<<17 + 1<<23,
			3,
			[]byte{131, 16, 8},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := intToBytes(tc.input, tc.bytes)
			if !bytes.Equal(result, tc.output) {
				t.Errorf("got %v; want %v", result, tc.output)
			}
		})
	}
}

func TestBytesByUniques(t *testing.T) {
	testCases := []struct {
		name   string
		input  int
		output int
	}{
		{"0", 0, 0},
		{"1", 1, 1},
		{"256", 256, 1},
		{"257", 257, 2},
		{"1^16", 1 << 16, 2},
		{"1^16+1", 1<<16 + 1, 3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := bytesByUniques(tc.input)
			if got != tc.output {
				t.Errorf("got %v; want %v", got, tc.output)
			}
		})
	}
}

func TestStringDeEncode(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			f, vals := stringify(tc.input)
			de := NewStringDeEncoder(vals)

			if got := de.Bytes(); got != tc.bytes {
				t.Errorf("got %v; want %v", got, tc.bytes)
			}

			for k, v := range tc.good {
				got, ok := de.Encode(f(k))
				if !ok {
					t.Error("got !ok; want ok")
				}
				if !bytes.Equal(v, got) {
					t.Errorf("got %v; want %v", got, v)
				}

				val := de.Decode(got)
				if val != f(k) {
					t.Errorf("got %v; want %v", val, k)
				}
			}

			for k, v := range tc.bad {
				_, ok := de.Encode(f(k))
				if ok {
					t.Error("got ok; want !ok")
				}

				lower, upper := de.Transform(f(k))
				if !bytes.Equal(lower, v[0]) {
					t.Errorf("got %v; want %v", lower, v[0])
				}

				if !bytes.Equal(upper, v[1]) {
					t.Errorf("got %v; want %v", upper, v[1])
				}
			}
		})
	}

}
