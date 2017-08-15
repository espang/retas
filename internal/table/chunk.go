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
	iter := &iterator{
		len:   len(c.data),
		c:     &c,
		start: idx,
		end:   idx + width,
	}
	return iter
}

func byteIndexer(width int, bins []int) func(v byteView) int {
	//lte := byteViewComparer(width)
	toInt := byteViewInt(width)
	max := bins[len(bins)-1]
	min := bins[0]
	binW := (max - min) / (len(bins) - 1)

	// return func(v byteView) int {
	// 	return sort.Search(len(bins), func(i int) bool { return lte(v, bins[i]) })
	// }

	return func(v byteView) int {
		iv := toInt(v)
		if iv >= max {
			return len(bins)
		}
		idx := ((iv - min) / binW) - 1
		dir := 0
		for {
			if iv < bins[idx+1] {
				if iv >= bins[idx] {
					return idx + 1
				}
				if dir == 1 || idx == 0 {
					return 0
				}
				idx--
				dir = -1
				continue
			}
			if dir == -1 || idx == len(bins)-1 {
				return len(bins)
			}
			idx++
			dir = 1
		}
	}
}

func byteViewInt(length int) func(byteView) int {
	if length == 1 {
		return func(b byteView) int {
			return int(b.bs[0])
		}
	}
	return func(bv byteView) int {
		var res int
		for i := 0; i < length; i++ {
			res += int(bv.bs[i]) << uint((length-i-1)*8)
		}
		return res
	}
}

func byteViewComparer(length int) func(byteView, int) bool {
	// Todo: write benchmark with and without this code
	// if length == 1 {
	// 	return func(bv byteView, cmp int) bool {
	// 		return bv.bs[0] < byte(cmp&0xff)
	// 	}
	// }
	// // if length == 2 {
	// 	return func(bv byteView, cmp int) bool {
	// 		return bv.bs[0] < byte(cmp&0xff00) ||
	// 			bv.bs[0] == byte(cmp&0xff00) && bv.bs[1] < byte(cmp&0xff)
	// 	}
	// }
	return func(bv byteView, cmp int) bool {
		for i := 0; i < length; i++ {
			if bv.bs[i] == byte(cmp>>uint((length-i-1)*8)&0xff) {
				continue
			}
			return bv.bs[i] < byte(cmp>>uint((length-i-1)*8)&0xff)
		}
		return false
	}
}

type byteView struct {
	bs []byte
}

func newByteView(r []byte, start, end int) byteView {
	return byteView{bs: r[start:end]}
}
