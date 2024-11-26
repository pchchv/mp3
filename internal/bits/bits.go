package bits

type Bits struct {
	vec     []byte
	bitPos  int
	bytePos int
}

func New(vec []byte) *Bits {
	return &Bits{
		vec: vec,
	}
}
func (b *Bits) Bit() int {
	if len(b.vec) <= b.bytePos {
		return 0
	}

	tmp := uint(b.vec[b.bytePos]) >> (7 - uint(b.bitPos))
	tmp &= 0x01
	b.bytePos += (b.bitPos + 1) >> 3
	b.bitPos = (b.bitPos + 1) & 0x07
	return int(tmp)
}

func (b *Bits) BitPos() int {
	return b.bytePos<<3 + b.bitPos
}

func (b *Bits) Bits(num int) int {
	if num == 0 {
		return 0
	} else if len(b.vec) <= b.bytePos {
		return 0
	}

	bb := make([]byte, 4)
	copy(bb, b.vec[b.bytePos:])
	tmp := (uint32(bb[0]) << 24) | (uint32(bb[1]) << 16) | (uint32(bb[2]) << 8) | (uint32(bb[3]))
	tmp <<= uint(b.bitPos)
	tmp >>= (32 - uint(num))
	b.bytePos += (b.bitPos + num) >> 3
	b.bitPos = (b.bitPos + num) & 0x07
	return int(tmp)
}

func (b *Bits) Tail(offset int) []byte {
	return b.vec[len(b.vec)-offset:]
}

func (b *Bits) SetPos(pos int) {
	b.bytePos = pos >> 3
	b.bitPos = pos & 0x7
}
