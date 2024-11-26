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

// Mode returns the channel mode, stored in position 7,6
func (f FrameHeader) Mode() consts.Mode {
	return consts.Mode((f & 0x000000c0) >> 6)
}

// UseIntensityStereo returns a boolean value indicating whether the
// frame uses intensity stereo.
func (f FrameHeader) UseIntensityStereo() bool {
	if f.Mode() != consts.ModeJointStereo {
		return false
	}
	return f.modeExtension()&0x1 != 0
}

// UseMSStereo returns a boolean value indicating whether the
// frame uses middle/side stereo.
func (f FrameHeader) UseMSStereo() bool {
	if f.Mode() != consts.ModeJointStereo {
		return false
	}
	return f.modeExtension()&0x2 != 0
}

// Copyright returns whether or not
// this recording is copywritten - stored in position 3
func (f FrameHeader) Copyright() int {
	return int(f&0x00000008) >> 3
}

// OriginalOrCopy returns whether or not
// this is an Original recording or a copy of one - stored in position 2
func (f FrameHeader) OriginalOrCopy() int {
	return int(f&0x00000004) >> 2
}

// Emphasis returns emphasis - the emphasis indication is here to
// tell the decoder that the file must be de-emphasized - stored in position 0,1
func (f FrameHeader) Emphasis() int {
	return int(f&0x00000003) >> 0
}

func (f FrameHeader) BytesPerFrame() int {
	return consts.SamplesPerGr * f.Granules() * 4
}

func (f FrameHeader) Granules() int {
	return consts.GranulesMpeg1 >> uint(f.LowSamplingFrequency()) // MPEG2 uses only 1 granule
}

func (f FrameHeader) NumberOfChannels() int {
	if f.Mode() == consts.ModeSingleChannel {
		return 1
	}
	return 2
}

// IsValid returns a boolean value indicating whether the header is valid or not.
func (f FrameHeader) IsValid() bool {
	const sync = 0xffe00000
	if (f & sync) != sync {
		return false
	}

	if f.ID() == consts.VersionReserved {
		return false
	}

	if f.BitrateIndex() == 15 {
		return false
	}

	if f.SamplingFrequency() == consts.SamplingFrequencyReserved {
		return false
	}

	if f.Layer() == consts.LayerReserved {
		return false
	}

	if f.Emphasis() == 2 {
		return false
	}

	return true
}

// modeExtension returns the mode_extension -
// for use with Joint Stereo -
// stored in position 4,5
func (f FrameHeader) modeExtension() int {
	return int(f&0x00000030) >> 4
}
