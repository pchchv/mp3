
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
