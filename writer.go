package bitio

import (
	"github.com/koji-hirono/memio"
)

// A Writer provides sequential bit-level writing.
type Writer struct {
	g memio.Grower
	n int
}

// NewWriter returns a new Writer
func NewWriter(g memio.Grower) *Writer {
	return &Writer{g: g}
}

func (w *Writer) Bytes() []byte {
	return w.g.Bytes()
}

func (w *Writer) Align() {
	w.n = (w.n + 7) &^ 0x7
}

func (w *Writer) Grow(nbits int) (err error) {
	nbytes := (nbits + 7) >> 3
	return w.g.Grow(nbytes)
}

func (w *Writer) Write(p []byte) (int, error) {
	m := len(p)
	n := w.n + m*8
	err := w.Grow(n)
	if err != nil {
		return 0, err
	}
	for _, c := range p {
		w.writebit(c, 8)
	}
	return m, nil
}

// WriteByte writes a single byte.
func (w *Writer) WriteByte(c byte) error {
	n := w.n + 8
	err := w.Grow(n)
	if err != nil {
		return err
	}
	w.writebit(c, 8)
	return nil
}

// WriteBool writes the specified boolean value as a single bit.
func (w *Writer) WriteBool(b bool) error {
	n := w.n + 1
	err := w.Grow(n)
	if err != nil {
		return err
	}
	if b {
		w.writebit(0x80, 1)
	} else {
		w.writebit(0, 1)
	}
	return nil
}

// WriteBits writes nbits bits.
func (w *Writer) WriteBits(p []byte, nbits int) error {
	n := w.n + nbits
	err := w.Grow(n)
	if err != nil {
		return err
	}
	i := nbits >> 3
	for _, c := range p[:i] {
		w.writebit(c, 8)
	}
	off := nbits & 7
	if off != 0 {
		w.writebit(p[i], off)
	}
	return nil
}

func (w *Writer) WriteBitField(v uint64, nbits int) error {
	v <<= 64 - nbits
	nbytes := (nbits + 7) >> 3
	p := make([]byte, nbytes)
	for i := 0; i < nbytes; i++ {
		k := 56 - (i << 3)
		p[i] = byte(v >> k)
	}
	return w.WriteBits(p, nbits)
}

func (w *Writer) writebit(c byte, nbits int) {
	buf := w.g.Bytes()
	i := w.n >> 3
	off := w.n & 7
	buf[i] |= c >> off
	if nbits+off > 8 {
		buf[i+1] = c << (8 - off)
	}
	w.n += nbits
}
