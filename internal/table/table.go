package table

import (
	"errors"
	"sort"
	"sync"
	"sync/atomic"
)

type Table interface {
	Query(columnID int) Iterator
	Histogram(columnID int, g []Range) []uint64
}

type Range struct {
	From int
	To   int
}

type table struct {
	entries []entry
	chunks  []chunk
}

func (t table) Query(columnID int) Iterator {
	return nil
}

func (t table) index(columnID int) (int, int) {
	for _, e := range t.entries {
		if e.Id == columnID {
			return e.Idx, e.Width
		}
	}
	return -1, 0
}

func (t table) Histogram(columnID int, bins ...int) ([]uint64, error) {
	if !sort.IntsAreSorted(bins) {
		return nil, errors.New("bins have to be sorted")
	}
	if len(bins) < 2 {
		return nil, errors.New("need at least 2 bins")
	}

	idx, width := t.index(columnID)
	if idx == -1 {
		return nil, errors.New("column not found")
	}
	if width > 1 {
		return nil, errors.New("invalid column for binning")
	}

	hist := make([]uint64, len(bins)-1)
	var wg sync.WaitGroup
	for _, c := range t.chunks {
		iter := c.query(idx, 1)
		wg.Add(1)
		go func(g *sync.WaitGroup) {
			chist := make([]uint64, len(bins))
			for iter.Next() {
				v := iter.Value()
				idx := sort.Search(len(bins), func(i int) bool { return v >= bins[i] })
				chist[idx]++
			}
			for i, v := range chist[:len(bins)-1] {
				_ = atomic.AddUint64(&hist[i], v)
			}
		}(&wg)
	}
	wg.Wait()

	return hist[:len(hist)-1], nil
}

type entry struct {
	Id    int
	Idx   int
	Width int
}

func New(rows ...[][]int) (Table, error) {
	return nil, nil
}