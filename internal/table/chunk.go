package table

type Iterator interface {
	Value() byteView
	Next() bool
}

type iterator struct {
	len    int
	cursor int
	c      *chunk
	start  int
	end    int
}

func (i *iterator) Value() byteView {
	return newByteView(i.c.data[i.cursor-1], i.start, i.end)
}

func (i *iterator) Next() bool {
	i.cursor++
	return i.cursor <= i.len
}

type chunk struct {
	data [][]byte
}

func (c chunk) query(idx, width int) Iterator {
	return &iterator{
		len:   len(c.data),
		c:     &c,
		start: idx,
		end:   idx + width,
	}
}

func byteViewLen1GTE(bv byteView, cmp int) bool {
	return bv.bs[0] >= byte(cmp&0xff)
}

type byteView struct {
	bs []byte
}

func newByteView(r []byte, start, end int) byteView {
	return byteView{bs: r[start:end]}
}
