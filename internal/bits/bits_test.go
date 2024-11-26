package bits_test

import (
	"testing"

	. "github.com/pchchv/mp3/internal/bits"
)

func TestBits(t *testing.T) {
	b1 := byte(85)  // 01010101
	b2 := byte(170) // 10101010
	b3 := byte(204) // 11001100
	b4 := byte(51)  // 00110011
	b := New([]byte{b1, b2, b3, b4})
	if b.Bits(1) != 0 {
		t.Fail()
	}

	if b.Bits(1) != 1 {
		t.Fail()
	}

	if b.Bits(1) != 0 {
		t.Fail()
	}

	if b.Bits(1) != 1 {
		t.Fail()
	}

	if b.Bits(8) != 90 /* 01011010 */ {
		t.Fail()
	}

	if b.Bits(12) != 2764 /* 101011001100 */ {
		t.Fail()
	}
}
