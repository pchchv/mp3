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
