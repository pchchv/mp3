package mp3

import (
	"errors"
	"io"
)

type source struct {
	reader io.Reader
	buf    []byte
	pos    int64
}

func (s *source) Seek(position int64, whence int) (n int64, err error) {
	seeker, ok := s.reader.(io.Seeker)
	if !ok {
		return 0, errors.New("mp3: source must be io.Seeker")
	}

	s.buf = nil
	if n, err = seeker.Seek(position, whence); err != nil {
		return 0, err
	} else {
		s.pos = n
	}

	return n, nil
}

func (s *source) ReadFull(buf []byte) (n int, err error) {
	var read int
	if s.buf != nil {
		read = copy(buf, s.buf)
		if len(s.buf) > read {
			s.buf = s.buf[read:]
		} else {
			s.buf = nil
		}

		if len(buf) == read {
			return read, nil
		}
	}

	if n, err = io.ReadFull(s.reader, buf[read:]); err == io.ErrUnexpectedEOF {
		// allow if all data can't be read
		err = io.EOF
	}

	s.pos += int64(n)
	return n + read, err
}
