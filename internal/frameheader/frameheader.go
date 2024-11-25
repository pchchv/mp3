
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
