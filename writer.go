package bitio

import (
	"io"
)

// A Writer provides sequential bit-level writing.
type Writer struct {
	w   io.Writer
	buf byte
	n   int
}

// NewWriter returns a new Writer that writes to w.
func NewWriter(w io.Writer) *Writer {
	return &Writer{w: w}
}

// Flush expands any buffered data to a single byte,
// and writes to the underlying io.Writer.
func (w *Writer) Flush() error {
	if w.n == 0 {
		return nil
	}
	_, err := w.w.Write([]byte{w.buf})
	if err != nil {
		return err
	}
	w.buf = 0
	w.n = 0
	return nil
}

// Write writes len(p) bytes to the underlying io.Writer.
func (w *Writer) Write(p []byte) (int, error) {
	if w.n == 0 {
		return w.w.Write(p)
	}
	for i, c := range p {
		x := w.buf | (c >> w.n)
		_, err := w.w.Write([]byte{x})
		if err != nil {
			return i, err
		}
		w.buf = c << (8 - w.n)
	}
	return len(p), nil
}

// WriteByte writes a single byte.
func (w *Writer) WriteByte(c byte) error {
	if w.n == 0 {
		_, err := w.w.Write([]byte{c})
		return err
	}
	x := w.buf | (c >> w.n)
	_, err := w.w.Write([]byte{x})
	if err != nil {
		return err
	}
	w.buf = c << (8 - w.n)
	return nil
}

// WriteBool writes the specified boolean value as a single bit.
func (w *Writer) WriteBool(b bool) error {
	x := w.buf
	if b {
		x |= byte(0x80) >> w.n
	}
	if w.n+1 < 8 {
		w.buf = x
		w.n++
	} else {
		_, err := w.w.Write([]byte{x})
		if err != nil {
			return err
		}
		w.buf = 0
		w.n = 0
	}
	return nil
}

// WriteBits writes nbits bits to the underlying io.Writer.
func (w *Writer) WriteBits(p []byte, nbits int) (int, error) {
	n := 0
	end := nbits >> 3
	if end != 0 {
		nb, err := w.Write(p[:end])
		if err != nil {
			return nb << 3, err
		}
		n = nb
	}

	off := nbits & 7
	if off == 0 {
		return nbits, nil
	}

	c := p[end]
	x := w.buf | (c >> w.n)
	if w.n+off < 8 {
		w.buf = x
	} else {
		_, err := w.w.Write([]byte{x})
		if err != nil {
			return n << 3, err
		}
		w.buf = c << (8 - w.n)
	}
	w.n = (w.n + off) & 7
	return nbits, nil
}
