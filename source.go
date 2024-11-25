package mp3

import "io"

type source struct {
	reader io.Reader
	buf    []byte
	pos    int64
}
