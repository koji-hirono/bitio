package bitio

import (
	"io"
)

// A Reader provides sequential bit-level access.
type Reader struct {
	buf []byte
	n   int
}

// NewReader returns a new Reader.
func NewReader(buf []byte) *Reader {
	return &Reader{buf: buf}
}

// Align skips any data
func (r *Reader) Align() {
	r.n = (r.n + 7) &^ 7
}

// Read reads data into p.
// It returns the number of bytes read into p.
// To read exactly len(p) bytes, use io.ReadFull(r, p).
func (r *Reader) Read(p []byte) (int, error) {
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

func (r *Reader) AtomicRead(p []byte) error {
	m := len(p)
	n := ((r.n + m*8) + 7) >> 3
	if n > len(r.buf) {
		return io.ErrUnexpectedEOF
	}
	for i := 0; i < m; i++ {
		p[i] = r.readbit(8)
	}
	return nil
}

// ReadByte reads and returns a single byte.
// If no byte is available, returns an error.
func (r *Reader) ReadByte() (byte, error) {
	n := ((r.n + 8) + 7) >> 3
	if n > len(r.buf) {
		return 0, io.EOF
	}
	return r.readbit(8), nil
}

// ReadBool reads a single bit and returns a boolean value.
// If a single bit is not available, returns an error.
func (r *Reader) ReadBool() (bool, error) {
	n := ((r.n + 1) + 7) >> 3
	if n > len(r.buf) {
		return false, io.ErrUnexpectedEOF
	}
	return r.readbit(1) == byte(0x80), nil
}

// ReadBits reads nbits bits.
func (r *Reader) ReadBits(p []byte, nbits int) error {
	n := ((r.n + nbits) + 7) >> 3
	if n > len(r.buf) {
		return io.ErrUnexpectedEOF
	}
	m := nbits >> 3
	for i := 0; i < m; i++ {
		p[i] = r.readbit(8)
	}
	off := nbits & 7
	if off != 0 {
		p[m] = r.readbit(off)
	}
	return nil
}

func (r *Reader) ReadBitField(nbits int) (uint64, error) {
	nbytes := (nbits + 7) >> 3
	p := make([]byte, nbytes)
	err := r.ReadBits(p, nbits)
	if err != nil {
		return 0, err
	}
	v := uint64(0)
	for i := 0; i < nbytes; i++ {
		k := 56 - (i << 3)
		v |= uint64(p[i]) << k
	}
	v >>= 64 - nbits
	return v, nil
}

func (r *Reader) readbit(nbits int) byte {
	i := r.n >> 3
	off := r.n & 7
	c := r.buf[i] << off
	if off+nbits > 8 {
		c |= r.buf[i+1] >> (8 - off)
	}
	r.n += nbits
	mask := ^byte(0) << (8 - nbits)
	c &= mask
	return c
}
