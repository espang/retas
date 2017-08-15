package encoder

import (
	"bytes"
	"testing"
)

func TestDecode(t *testing.T) {
	testCases := []struct {
		name  string
		input []int
		bytes int
	}{
		{
			"two ints",
			[]int{1, 10},
			1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			de := NewIntDeEncoder(tc.input)

			_, _ = de.Encode(1)
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
