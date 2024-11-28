package mp3

import (
	"io"

	"github.com/pchchv/mp3/internal/consts"
	"github.com/pchchv/mp3/internal/frame"
	"github.com/pchchv/mp3/internal/frameheader"
)

const invalidLength = -1

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

func (d *Decoder) readFrame() (err error) {
	d.frame, _, err = frame.Read(d.source, d.source.pos, d.frame)
	if err != nil {
		if err == io.EOF {
			return io.EOF
		}

		if _, ok := err.(*consts.UnexpectedEOF); ok {
			return io.EOF
		}

		return err
	}

	d.buf = append(d.buf, d.frame.Decode()...)
	return nil
}

func (d *Decoder) ensureFrameStartsAndLength() error {
	if d.length != invalidLength {
		return nil
	}

	if _, ok := d.source.reader.(io.Seeker); !ok {
		return nil
	}

	// keep the current position
	pos, err := d.source.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}

	if err = d.source.rewind(); err != nil {
		return err
	}

	if err = d.source.skipTags(); err != nil {
		return err
	}

	var l int64
	for {
		h, pos, err := frameheader.Read(d.source, d.source.pos)
		if err != nil {
			if err == io.EOF {
				break
			}

			if _, ok := err.(*consts.UnexpectedEOF); ok {
				break
			}

			return err
		}

		d.frameStarts = append(d.frameStarts, pos)
		d.bytesPerFrame = int64(h.BytesPerFrame())
		l += d.bytesPerFrame

		framesize, err := h.FrameSize()
		if err != nil {
			return err
		}

		buf := make([]byte, framesize-4)
		if _, err := d.source.ReadFull(buf); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}
	d.length = l

	if _, err := d.source.Seek(pos, io.SeekStart); err != nil {
		return err
	}

	return nil
}
