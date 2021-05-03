package bitio

import (
	"io"
)

// A Reader provides sequential bit-level access.
type Reader struct {
	r   io.Reader
	buf byte
	n   int
}

// NewReader returns a new Reader that reads from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{r: r}
}

// DiscardBuffer skips any buffered data.
func (r *Reader) DiscardBuffer() {
	r.n = 0
	r.buf = 0
}

// Read reads data into p.
// It returns the number of bytes read into p.
// To read exactly len(p) bytes, use io.ReadFull(r, p).
func (r *Reader) Read(p []byte) (int, error) {
	if r.n == 0 {
		return r.r.Read(p)
	}
	n := len(p)
	for i := 0; i < n; i++ {
		c, err := r.ReadByte()
		if err != nil {
			return i, err
		}
		p[i] = c
	}
	return n, nil
}

// ReadByte reads and returns a single byte.
// If no byte is available, returns an error.
func (r *Reader) ReadByte() (byte, error) {
	var buf [1]byte
	_, err := r.r.Read(buf[:])
	if err != nil {
		return 0, err
	}
	if r.n == 0 {
		return buf[0], nil
	}
	c := r.buf
	c |= buf[0] >> r.n
	r.buf = buf[0] << (8 - r.n)
	return c, nil
}

// ReadBool reads a single bit and returns a boolean value.
// If a single bit is not available, returns an error.
func (r *Reader) ReadBool() (bool, error) {
	var b bool
	if r.n == 0 {
		var buf [1]byte
		n, _ := r.r.Read(buf[:])
		if n == 0 {
			return b, io.EOF
		}
		r.buf = buf[0]
		r.n = 8
	}

	if r.buf&0x80 == 0x80 {
		b = true
	}
	r.buf <<= 1
	r.n--
	return b, nil
}

// ReadBits reads nbits bits.
func (r *Reader) ReadBits(p []byte, nbits int) (int, error) {
	var err error
	end := nbits >> 3
	if end != 0 {
		n, e := r.Read(p[:end])
		if n != end {
			return n << 3, e
		}
		err = e
	}

	off := nbits & 7
	if off == 0 {
		return nbits, err
	}

	p[end] = r.buf

	if off <= r.n {
		p[end] >>= 8 - off
		p[end] <<= 8 - off
		r.buf <<= off
		r.n -= off
		return nbits, err
	}

	n := r.n
	r.buf = 0
	r.n = 0

	var buf [1]byte
	_, e := r.r.Read(buf[:])
	if e != nil {
		return end<<3 + n, e
	}
	p[end] |= buf[0] >> n
	p[end] >>= 8 - off
	p[end] <<= 8 - off
	r.buf = buf[0] << (off - n)
	r.n = 8 - off + n

	return nbits, nil
}
