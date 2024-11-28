package mp3

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func BenchmarkDecode(b *testing.B) {
	buf, err := os.ReadFile("example/classic.mp3")
	if err != nil {
		b.Fatal(err)
	}

	src := bytes.NewReader(buf)
	for i := 0; i < b.N; i++ {
		if _, err := src.Seek(0, io.SeekStart); err != nil {
			b.Fatal(err)
		}

		d, err := NewDecoder(src)
		if err != nil {
			b.Fatal(err)
		}

		if _, err := io.ReadAll(d); err != nil {
			b.Fatal(err)
		}
	}
}
