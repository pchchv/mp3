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
