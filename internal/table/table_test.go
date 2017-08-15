package table

import (
	"crypto/rand"
	"fmt"
	"testing"
)

// table to test 8 chunks when a chunk is upto 16bit
var (
	tab64k   Table
	tab524k  Table
	tab1048k Table
)

func init() {
	tab64k = CreateTable(256 * 256)
	tab524k = CreateTable(8 * 256 * 256)
	tab1048k = CreateTable(16 * 256 * 256)
}

func row(width, length int) [][]byte {
	arr := make([][]byte, length)
	for i := range arr {
		b := make([]byte, 8)
		_, _ = rand.Read(b)
		arr[i] = b
	}
	return arr
}

func CreateTable(l int) Table {
	t, err := New(
		Column{
			ID:    1,
			Width: 8,
			data:  row(8, l),
		},
		Column{
			ID:    2,
			Width: 9,
			data:  row(9, l),
		},
		Column{
			ID:    3,
			Width: 1,
			data:  row(1, l),
		},
		Column{
			ID:    4,
			Width: 2,
			data:  row(2, l),
		},
		Column{
			ID:    5,
			Width: 3,
			data:  row(3, l),
		},
		Column{
			ID:    6,
			Width: 4,
			data:  row(4, l),
		},
	)
	if err != nil {
		panic(err)
	}
	return t
}

func sum(vs ...uint64) uint64 {
	if len(vs) == 0 {
		return 0
	}
	if len(vs) == 1 {
		return vs[0]
	}
	return vs[0] + sum(vs[1:]...)
}

func TestTable1(t *testing.T) {
	bins := []int{0, 40, 80, 120, 160, 200, 240, 255}
	h, err := tab64k.Histogram(3, bins...)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("histogram: ", h, sum(h...))
}
func TestTable2(t *testing.T) {
	bins := []int{0 << 8, 40 << 8, 80 << 8, 120 << 8, 160 << 8, 200 << 8, 240 << 8, 255 << 8}
	h, err := tab64k.Histogram(4, bins...)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("histogram: ", h, sum(h...))
}

var result []uint64

func benchmarkTable(b *testing.B, t Table, l int, bins ...int) {
	var h []uint64
	var err error
	for i := range bins {
		bins[i] = bins[i] << uint((l-1)*8)
	}
	var idx int
	switch l {
	case 1:
		idx = 3
	case 2:
		idx = 4
	case 3:
		idx = 5
	case 4:
		idx = 6
	}
	b.ResetTimer()

	for n := 1; n < b.N; n++ {
		h, err = t.Histogram(idx, bins...)
		if err != nil {
			b.Error(err)
		}
	}
	result = h
}

func BenchmarkTableSmall(b *testing.B) {
	bm := struct {
		t    Table
		len  int
		bins []int
	}{tab64k, 1, []int{0, 40, 80, 120, 160, 200, 240, 255}}

	name := fmt.Sprintf("Tablesize %d chunks, ElementBytes %d", len(bm.t.(table).chunks), bm.len)
	b.Run(name, func(b *testing.B) {
		benchmarkTable(b, bm.t, bm.len, bm.bins...)
	})
}

func BenchmarkTableSmall2(b *testing.B) {
	bm := struct {
		t    Table
		len  int
		bins []int
	}{tab64k, 2, []int{0, 40, 80, 120, 160, 200, 240, 255}}

	name := fmt.Sprintf("Tablesize %d chunks, ElementBytes %d", len(bm.t.(table).chunks), bm.len)
	b.Run(name, func(b *testing.B) {
		benchmarkTable(b, bm.t, bm.len, bm.bins...)
	})
}

func BenchmarkTable(b *testing.B) {
	benchmarks := []struct {
		t    Table
		len  int
		bins []int
	}{
		{tab64k, 2, []int{0, 40, 80, 120, 160, 200, 240, 255}},
		{tab64k, 1, []int{0, 40, 80, 120, 160, 200, 240, 255}},
		{tab64k, 3, []int{0, 40, 80, 120, 160, 200, 240, 255}},
		{tab64k, 4, []int{0, 40, 80, 120, 160, 200, 240, 255}},
		{tab524k, 1, []int{0, 40, 80, 120, 160, 200, 240, 255}},
		{tab524k, 2, []int{0, 40, 80, 120, 160, 200, 240, 255}},
		{tab524k, 3, []int{0, 40, 80, 120, 160, 200, 240, 255}},
		{tab524k, 4, []int{0, 40, 80, 120, 160, 200, 240, 255}},
		{tab1048k, 1, []int{0, 40, 80, 120, 160, 200, 240, 255}},
		{tab1048k, 2, []int{0, 40, 80, 120, 160, 200, 240, 255}},
		{tab1048k, 3, []int{0, 40, 80, 120, 160, 200, 240, 255}},
		{tab1048k, 4, []int{0, 40, 80, 120, 160, 200, 240, 255}},
	}
	for _, bm := range benchmarks {
		name := fmt.Sprintf("Tablesize %d chunks, ElementBytes %d", len(bm.t.(table).chunks), bm.len)
		b.Run(name, func(b *testing.B) {
			benchmarkTable(b, bm.t, bm.len, bm.bins...)
		})
	}
}
