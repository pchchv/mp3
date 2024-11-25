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

func (s *source) Unread(buf []byte) {
	s.buf = append(s.buf, buf...)
	s.pos -= int64(len(buf))
}

func (s *source) rewind() error {
	if _, err := s.Seek(0, io.SeekStart); err != nil {
		return err
	}

	s.pos = 0
	s.buf = nil
	return nil
}

func (s *source) skipTags() error {
	buf := make([]byte, 3)
	if _, err := s.ReadFull(buf); err != nil {
		return err
	}

	switch string(buf) {
	case "TAG":
		buf = make([]byte, 125)
		if _, err := s.ReadFull(buf); err != nil {
			return err
		}
	case "ID3":
		// skip version (2 bytes) and flag (1 byte)
		buf := make([]byte, 3)
		if _, err := s.ReadFull(buf); err != nil {
			return err
		}

		buf = make([]byte, 4)
		n, err := s.ReadFull(buf)
		if err != nil {
			return err
		} else if n != 4 {
			return nil
		}

		size := (uint32(buf[0]) << 21) | (uint32(buf[1]) << 14) |
			(uint32(buf[2]) << 7) | uint32(buf[3])
		buf = make([]byte, size)
		if _, err := s.ReadFull(buf); err != nil {
			return err
		}
	default:
		s.Unread(buf)
	}

	return nil
}
