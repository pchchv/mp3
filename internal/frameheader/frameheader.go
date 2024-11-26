package frameheader

import (
	"errors"

	"github.com/pchchv/mp3/internal/consts"
)

// mepg1FrameHeader is MPEG1 Layer 1-3 frame header.
type FrameHeader uint32

// ID returns this header's ID stored in position 20,19
func (f FrameHeader) ID() consts.Version {
	return consts.Version((f & 0x00180000) >> 19)
}

// Layer returns the mpeg layer of this frame stored in position 18,17
func (f FrameHeader) Layer() consts.Layer {
	return consts.Layer((f & 0x00060000) >> 17)
}

// BirateIndex returns the bitrate index stored in position 15,12
func (f FrameHeader) BitrateIndex() int {
	return int(f&0x0000f000) >> 12
}

// LowSamplingFrequency returns whether the frame is encoded in a
// low sampling frequency => 0 = MPEG-1, 1 = MPEG-2/2.5
func (f FrameHeader) LowSamplingFrequency() int {
	if f.ID() == consts.Version1 {
		return 0
	}
	return 1
}

// SamplingFrequency returns the SamplingFrequency in Hz stored in position 11,10
func (f FrameHeader) SamplingFrequency() consts.SamplingFrequency {
	return consts.SamplingFrequency(int(f&0x00000c00) >> 10)
}

func (f FrameHeader) SamplingFrequencyValue() (int, error) {
	switch f.SamplingFrequency() {
	case 0:
		return 44100 >> uint(f.LowSamplingFrequency()), nil
	case 1:
		return 48000 >> uint(f.LowSamplingFrequency()), nil
	case 2:
		return 32000 >> uint(f.LowSamplingFrequency()), nil
	}
	return 0, errors.New("mp3: frame header has invalid sample frequency")
}

// ProtectionBit returns the protection bit stored in position 16
func (f FrameHeader) ProtectionBit() int {
	return int(f&0x00010000) >> 16
}

// PaddingBit returns the padding bit stored in position 9
func (f FrameHeader) PaddingBit() int {
	return int(f&0x00000200) >> 9
}

// PrivateBit returns the private bit stored in
// position 8 - this bit may be used to store arbitrary data to be
// used by an application
func (f FrameHeader) PrivateBit() int {
	return int(f&0x00000100) >> 8
}
