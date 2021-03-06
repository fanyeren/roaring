package roaring

const (
	arrayDefaultMaxSize = 4096
	maxCapacity         = 1 << 16
	serial_cookie       = 12346
)

func getSizeInBytesFromCardinality(card int) int {
	if card > arrayDefaultMaxSize {
		return maxCapacity / 8
	} else {
		return 2 * int(card)
	}
}

// should be replaced with optimized assembly instructions
func numberOfTrailingZeros(i uint64) int {
	if i == 0 {
		return 64
	}
	x := i
	n := int64(63)
	y := x << 32
	if y != 0 {
		n -= 32
		x = y
	}
	y = x << 16
	if y != 0 {
		n -= 16
		x = y
	}
	y = x << 8
	if y != 0 {
		n -= 8
		x = y
	}
	y = x << 4
	if y != 0 {
		n -= 4
		x = y
	}
	y = x << 2
	if y != 0 {
		n -= 2
		x = y
	}
	return int(n - int64(uint64(x<<1)>>63))
}

func fill(arr []uint64, val uint64) {
	for i := range arr {
		arr[i] = val
	}
}
func fillRange(arr []uint64, start, end int, val uint64) {
	for i := start; i < end; i++ {
		arr[i] = val
	}
}

func fillArrayAND(container []uint16, bitmap1, bitmap2 []uint64) {
	if len(bitmap1) != len(bitmap2) {
		panic("array lengths don't match")
	}
	// TODO: rewrite in assembly
	pos := 0
	for k := range bitmap1 {
		bitset := bitmap1[k] & bitmap2[k]
		for bitset != 0 {
			t := bitset & -bitset
			container[pos] = uint16((k*64 + int(popcount(t-1))))
			pos = pos + 1
			bitset ^= t
		}
	}
}

func fillArrayANDNOT(container []uint16, bitmap1, bitmap2 []uint64) {
	if len(bitmap1) != len(bitmap2) {
		panic("array lengths don't match")
	}
	// TODO: rewrite in assembly
	pos := 0
	for k := range bitmap1 {
		bitset := bitmap1[k] &^ bitmap2[k]
		for bitset != 0 {
			t := bitset & -bitset
			container[pos] = uint16((k*64 + int(popcount(t-1))))
			pos = pos + 1
			bitset ^= t
		}
	}
}

func fillArrayXOR(container []uint16, bitmap1, bitmap2 []uint64) {
	if len(bitmap1) != len(bitmap2) {
		panic("array lengths don't match")
	}
	// TODO: rewrite in assembly
	pos := 0
	for k := 0; k < len(bitmap1); k++ {
		bitset := bitmap1[k] ^ bitmap2[k]
		for bitset != 0 {
			t := bitset & -bitset
			container[pos] = uint16((k*64 + int(popcount(t-1))))
			pos = pos + 1
			bitset ^= t
		}
	}
}

func highbits(x int) uint16 {
	u := uint(x)
	return uint16(u >> 16)
}
func lowbits(x int) uint16 {
	return uint16(x & 0xFFFF)
}

func maxLowBit() uint16 {
	return uint16(0xFFFF)
}

func toIntUnsigned(x uint16) int {
	return int(x & 0xFFFF)
}
