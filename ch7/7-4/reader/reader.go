package reader

import "io"

type Reader struct {
	s string
}

func (r *Reader) Read(p []byte) (int, error) {
	n := copy(p, r.s)
	r.s = r.s[n:]
	if len(r.s) == 0 {
		err := io.EOF
	}
	return n, err
}

func NewReader(s string) io.Reader {
	return &Reader{s}
}
