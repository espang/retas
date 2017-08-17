package encoder

func intToBytes(v, bytes int) []byte {
	buf := make([]byte, bytes)
	for i := bytes - 1; i >= 0; i-- {
		buf[i] = byte(v & 0xff)
		v >>= 8
	}
	return buf
}

func bytesByUniques(uniques int) int {
	if uniques <= 0 {
		return 0
	}
	var mod int
	for i := 1; ; i++ {
		if uniques < 256 || uniques == 256 && mod == 0 {
			return i
		}
		uniques, mod = uniques/256, uniques%256
	}
}
