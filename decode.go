package mp3

import "github.com/pchchv/mp3/internal/frame"

// Decoder is a MP3-decoded stream.
// Decoder decodes its underlying source on the fly.
type Decoder struct {
	source        *source
	sampleRate    int
	length        int64
	frameStarts   []int64
	buf           []byte
	frame         *frame.Frame
	pos           int64
	bytesPerFrame int64
}

// Length returns the total size in bytes.
// Length returns -1 when the total size is not available
// e.g. when the given source is not io.Seeker.
func (d *Decoder) Length() int64 {
	return d.length
}

// SampleRate returns the sample rate like 44100.
// Note that the sample rate is retrieved from the first frame.
func (d *Decoder) SampleRate() int {
	return d.sampleRate
}
