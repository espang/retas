package table

type Iterator interface {
	Value() int
	Next() bool
}

type iterator struct {
	end    int
	cursor int
	c      *chunk
	s      int
}

func (i *iterator) Value() int {
	return i.c.data[i.cursor-1][i.s]
}

func (i *iterator) Next() bool {
	i.cursor++
	return i.cursor <= i.end
}

type chunk struct {
	data [][]int
}

func (c chunk) query(idx, width int) Iterator {
	return &iterator{
		end: len(c.data),
		c:   &c,
		s:   idx,
	}
}
