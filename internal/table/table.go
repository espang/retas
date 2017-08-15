package table

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
)

type Table interface {
	//Query(columnID int) Iterator
	Histogram(columnID int, bins ...int) ([]uint64, error)
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
		if e.ID == columnID {
			return e.Idx, e.Width
		}
	}
	return -1, 0
}

func (t table) Histogram(columnID int, bins ...int) ([]uint64, error) {
	if !sort.IntsAreSorted(bins) {
		return nil, fmt.Errorf("bins have to be sorted %v", bins)
	}
	if len(bins) < 2 {
		return nil, errors.New("need at least 2 bins")
	}

	idx, width := t.index(columnID)
	if idx == -1 {
		return nil, errors.New("column not found")
	}

	indexer := byteIndexer(width, bins)

	hist := make([]uint64, len(bins)-1)
	work := make(chan Iterator, 4)
	var wg sync.WaitGroup
	for w := 0; w < 8; w++ {
		wg.Add(1)
		go func(g *sync.WaitGroup) {
			chist := make([]uint64, len(bins)+1)
			for iter := range work {
				for iter.Next() {
					chist[indexer(iter.Value())]++
				}
			}
			for i, v := range chist[1:len(bins)] {
				_ = atomic.AddUint64(&hist[i], v)
			}
			g.Done()
		}(&wg)
	}
	for _, c := range t.chunks {
		work <- c.query(idx, width)
	}
	close(work)
	wg.Wait()

	return hist, nil
}

type entry struct {
	ID    int
	Idx   int
	Width int
}

type Column struct {
	ID    int // address for the column
	Width int
	data  [][]byte
}

func New(cols ...Column) (Table, error) {
	s := map[int]struct{}{}
	s2 := map[int]struct{}{}
	var idx, length int
	var entries []entry
	for _, col := range cols {
		length = len(col.data)
		s2[len(col.data)] = struct{}{}
		s[col.ID] = struct{}{}

		entries = append(entries, entry{
			ID:    col.ID,
			Idx:   idx,
			Width: col.Width,
		})
		idx += col.Width
	}
	if len(s2) > 1 {
		return table{}, errors.New("n")
	}
	if len(s) != len(cols) {
		return table{}, errors.New("o")
	}

	chunksize := 256 * 256
	var chunks []chunk
	var current chunk
	for i := 0; i < length; i++ {
		if i != 0 && i%chunksize == 0 {
			chunks = append(chunks, current)
			current = chunk{}
		}
		var row []byte
		for _, col := range cols {
			row = append(row, col.data[i]...)
		}
		current.data = append(current.data, row)
	}
	if len(current.data) > 0 {
		chunks = append(chunks, current)
	}

	return table{
		chunks:  chunks,
		entries: entries,
	}, nil
}
