package bitio

import (
	"bytes"
	"errors"
	"testing"
)

func TestWriter_Flush(t *testing.T) {
	cases := []struct {
		name string
		buf  byte
		n    int
		out  []byte
	}{
		{
			name: "empty",
			buf:  0,
			n:    0,
			out:  []byte{},
		},
		{
			name: "1B",
			buf:  0b1000_0000,
			n:    1,
			out:  []byte{0b1000_0000},
		},
		{
			name: "0B",
			buf:  0b0000_0000,
			n:    1,
			out:  []byte{0b0000_0000},
		},
		{
			name: "11B",
			buf:  0b1100_0000,
			n:    2,
			out:  []byte{0b1100_0000},
		},
		{
			name: "10B",
			buf:  0b1000_0000,
			n:    2,
			out:  []byte{0b1000_0000},
		},
		{
			name: "01B",
			buf:  0b0100_0000,
			n:    2,
			out:  []byte{0b0100_0000},
		},
		{
			name: "111B",
			buf:  0b1110_0000,
			n:    3,
			out:  []byte{0b1110_0000},
		},
		{
			name: "1111B",
			buf:  0b1111_0000,
			n:    4,
			out:  []byte{0b1111_0000},
		},
		{
			name: "1_1111B",
			buf:  0b1111_1000,
			n:    5,
			out:  []byte{0b1111_1000},
		},
		{
			name: "11_1111B",
			buf:  0b1111_1100,
			n:    6,
			out:  []byte{0b1111_1100},
		},
		{
			name: "111_1111B",
			buf:  0b1111_1110,
			n:    7,
			out:  []byte{0b1111_1110},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBuffer([]byte{})
			w := NewWriter(b)
			w.buf = tt.buf
			w.n = tt.n
			err := w.Flush()
			if err != nil {
				t.Fatal(err)
			}
			out := b.Bytes()
			if !bytes.Equal(out, tt.out) {
				t.Errorf("out want %x; but got %x\n", tt.out, out)
			}
			if w.buf != 0 {
				t.Errorf("w.buf want 0; but got %x\n", w.buf)
			}
			if w.n != 0 {
				t.Errorf("w.n want 0; but got %v\n", w.n)
			}
		})
	}
}

func TestWriter_WriteByte(t *testing.T) {
	cases := []struct {
		name string
		buf  byte
		n    int
		c    byte
		out  []byte
		pbuf byte
		pn   int
	}{
		{
			name: "Offset 0 Aligned",
			buf:  0,
			n:    0,
			c:    0b1110_0101,
			out:  []byte{0b1110_0101},
			pbuf: 0,
			pn:   0,
		},
		{
			name: "Offset 1 Unaligned",
			buf:  0,
			n:    1,
			c:    0b1110_0101,
			out:  []byte{0b0111_0010},
			pbuf: 0b1000_0000,
			pn:   1,
		},
		{
			name: "Offset 2 Unaligned",
			buf:  0,
			n:    2,
			c:    0b1110_0101,
			out:  []byte{0b0011_1001},
			pbuf: 0b0100_0000,
			pn:   2,
		},
		{
			name: "Offset 3 Unaligned",
			buf:  0,
			n:    3,
			c:    0b1110_0101,
			out:  []byte{0b0001_1100},
			pbuf: 0b1010_0000,
			pn:   3,
		},
		{
			name: "Offset 4 Unaligned",
			buf:  0,
			n:    4,
			c:    0b1110_0101,
			out:  []byte{0b0000_1110},
			pbuf: 0b0101_0000,
			pn:   4,
		},
		{
			name: "Offset 5 Unaligned",
			buf:  0,
			n:    5,
			c:    0b1110_0101,
			out:  []byte{0b0000_0111},
			pbuf: 0b0010_1000,
			pn:   5,
		},
		{
			name: "Offset 6 Unaligned",
			buf:  0,
			n:    6,
			c:    0b1110_0101,
			out:  []byte{0b0000_0011},
			pbuf: 0b1001_0100,
			pn:   6,
		},
		{
			name: "Offset 7 Unaligned",
			buf:  0,
			n:    7,
			c:    0b1110_0101,
			out:  []byte{0b0000_0001},
			pbuf: 0b1100_1010,
			pn:   7,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBuffer([]byte{})
			w := NewWriter(b)
			w.buf = tt.buf
			w.n = tt.n
			err := w.WriteByte(tt.c)
			if err != nil {
				t.Fatal(err)
			}
			out := b.Bytes()
			if !bytes.Equal(out, tt.out) {
				t.Errorf("out want %x; but got %x\n", tt.out, out)
			}
			if w.buf != tt.pbuf {
				t.Errorf("w.buf want %x; but got %x\n", tt.pbuf, w.buf)
			}
			if w.n != tt.pn {
				t.Errorf("w.n want %v; but got %v\n", tt.pn, w.n)
			}
		})
	}
}

func TestWriter_WriteBool(t *testing.T) {
	cases := []struct {
		name string
		buf  byte
		n    int
		b    bool
		out  []byte
		pbuf byte
		pn   int
	}{
		{
			name: "Offset 0 Aligned",
			buf:  0,
			n:    0,
			b:    true,
			out:  []byte{},
			pbuf: 0b1000_0000,
			pn:   1,
		},
		{
			name: "Offset 1 Unaligned",
			buf:  0,
			n:    1,
			b:    true,
			out:  []byte{},
			pbuf: 0b0100_0000,
			pn:   2,
		},
		{
			name: "Offset 2 Unaligned",
			buf:  0,
			n:    2,
			b:    true,
			out:  []byte{},
			pbuf: 0b0010_0000,
			pn:   3,
		},
		{
			name: "Offset 3 Unaligned",
			buf:  0,
			n:    3,
			b:    true,
			out:  []byte{},
			pbuf: 0b0001_0000,
			pn:   4,
		},
		{
			name: "Offset 4 Unaligned",
			buf:  0,
			n:    4,
			b:    true,
			out:  []byte{},
			pbuf: 0b0000_1000,
			pn:   5,
		},
		{
			name: "Offset 5 Unaligned",
			buf:  0,
			n:    5,
			b:    true,
			out:  []byte{},
			pbuf: 0b0000_0100,
			pn:   6,
		},
		{
			name: "Offset 6 Unaligned",
			buf:  0,
			n:    6,
			b:    true,
			out:  []byte{},
			pbuf: 0b0000_0010,
			pn:   7,
		},
		{
			name: "Offset 7 Unaligned",
			buf:  0,
			n:    7,
			b:    true,
			out:  []byte{0b0000_0001},
			pbuf: 0,
			pn:   0,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBuffer([]byte{})
			w := NewWriter(b)
			w.buf = tt.buf
			w.n = tt.n
			err := w.WriteBool(tt.b)
			if err != nil {
				t.Fatal(err)
			}
			out := b.Bytes()
			if !bytes.Equal(out, tt.out) {
				t.Errorf("out want %x; but got %x\n", tt.out, out)
			}
			if w.buf != tt.pbuf {
				t.Errorf("w.buf want %x; but got %x\n", tt.pbuf, w.buf)
			}
			if w.n != tt.pn {
				t.Errorf("w.n want %v; but got %v\n", tt.pn, w.n)
			}
		})
	}
}

func TestWriter_Write(t *testing.T) {
	cases := []struct {
		name string
		buf  byte
		n    int
		p    []byte
		out  []byte
		pbuf byte
		pn   int
	}{
		{
			name: "Offset 0 Aligned",
			buf:  0,
			n:    0,
			p:    []byte{0b1100_0101, 0b0001_1010, 0b0010_1011, 0b0100_1001},
			out:  []byte{0b1100_0101, 0b0001_1010, 0b0010_1011, 0b0100_1001},
			pbuf: 0,
			pn:   0,
		},
		{
			name: "Offset 1 Unaligned",
			buf:  0b1000_0000,
			n:    1,
			p:    []byte{0b1100_0101, 0b0001_1010, 0b0010_1011, 0b0100_1001},
			out:  []byte{0b1110_0010, 0b1000_1101, 0b0001_0101, 0b1010_0100},
			pbuf: 0b1000_0000,
			pn:   1,
		},
		{
			name: "Offset 2 Unaligned",
			buf:  0b0100_0000,
			n:    2,
			p:    []byte{0b1100_0101, 0b0001_1010, 0b0010_1011, 0b0100_1001},
			out:  []byte{0b0111_0001, 0b0100_0110, 0b1000_1010, 0b1101_0010},
			pbuf: 0b0100_0000,
			pn:   2,
		},
		{
			name: "Offset 3 Unaligned",
			buf:  0b1010_0000,
			n:    3,
			p:    []byte{0b1100_0101, 0b0001_1010, 0b0010_1011, 0b0100_1001},
			out:  []byte{0b1011_1000, 0b1010_0011, 0b0100_0101, 0b0110_1001},
			pbuf: 0b0010_0000,
			pn:   3,
		},
		{
			name: "Offset 4 Unaligned",
			buf:  0b0101_0000,
			n:    4,
			p:    []byte{0b1100_0101, 0b0001_1010, 0b0010_1011, 0b0100_1001},
			out:  []byte{0b0101_1100, 0b0101_0001, 0b1010_0010, 0b1011_0100},
			pbuf: 0b1001_0000,
			pn:   4,
		},
		{
			name: "Offset 5 Unaligned",
			buf:  0b1010_1000,
			n:    5,
			p:    []byte{0b1100_0101, 0b0001_1010, 0b0010_1011, 0b0100_1001},
			out:  []byte{0b1010_1110, 0b0010_1000, 0b1101_0001, 0b0101_1010},
			pbuf: 0b0100_1000,
			pn:   5,
		},
		{
			name: "Offset 6 Unaligned",
			buf:  0b0101_0100,
			n:    6,
			p:    []byte{0b1100_0101, 0b0001_1010, 0b0010_1011, 0b0100_1001},
			out:  []byte{0b0101_0111, 0b0001_0100, 0b0110_1000, 0b1010_1101},
			pbuf: 0b0010_0100,
			pn:   6,
		},
		{
			name: "Offset 7 Unaligned",
			buf:  0b1010_1010,
			n:    7,
			p:    []byte{0b1100_0101, 0b0001_1010, 0b0010_1011, 0b0100_1001},
			out:  []byte{0b1010_1011, 0b1000_1010, 0b0011_0100, 0b0101_0110},
			pbuf: 0b1001_0010,
			pn:   7,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBuffer([]byte{})
			w := NewWriter(b)
			w.buf = tt.buf
			w.n = tt.n
			m, err := w.Write(tt.p)
			if err != nil {
				t.Fatal(err)
			}
			if m != len(tt.p) {
				t.Errorf("m want %x; but got %x\n", len(tt.p), m)
			}
			out := b.Bytes()
			if !bytes.Equal(out, tt.out) {
				t.Errorf("out want %x; but got %x\n", tt.out, out)
			}
			if w.buf != tt.pbuf {
				t.Errorf("w.buf want %x; but got %x\n", tt.pbuf, w.buf)
			}
			if w.n != tt.pn {
				t.Errorf("w.n want %v; but got %v\n", tt.pn, w.n)
			}
		})
	}
}

func TestWriter_WriteBits(t *testing.T) {
	cases := []struct {
		name  string
		buf   byte
		n     int
		v     uint64
		p     []byte
		nbits int
		on    int
		out   []byte
		pbuf  byte
		pn    int
	}{
		{
			name:  "Offset 0 Aligned",
			buf:   0,
			n:     0,
			p:     []byte{0b0010_0101, 0b0001_1010},
			nbits: 16,
			on:    16,
			out:   []byte{0b0010_0101, 0b0001_1010},
			pbuf:  0,
			pn:    0,
		},
		{
			name:  "Offset 0 Aligned",
			buf:   0,
			n:     0,
			p:     []byte{0b1001_0100, 0b0110_1000, 0b1010_1100},
			nbits: 22,
			on:    22,
			out:   []byte{0b1001_0100, 0b0110_1000},
			pbuf:  0b1010_1100,
			pn:    6,
		},
		{
			name:  "Offset 1 Unaligned",
			buf:   0b1000_0000,
			n:     1,
			p:     []byte{0b1001_0100, 0b0110_1000, 0b1010_1100},
			nbits: 22,
			on:    22,
			out:   []byte{0b1100_1010, 0b0011_0100},
			pbuf:  0b0101_0110,
			pn:    7,
		},
		{
			name:  "Offset 2 Unaligned",
			buf:   0b0100_0000,
			n:     2,
			p:     []byte{0b1001_0100, 0b0110_1000, 0b1010_1100},
			nbits: 22,
			on:    22,
			out:   []byte{0b0110_0101, 0b0001_1010, 0b0010_1011},
			pbuf:  0,
			pn:    0,
		},
		{
			name:  "Offset 3 Unaligned",
			buf:   0b1010_0000,
			n:     3,
			p:     []byte{0b1001_0100, 0b0110_1000, 0b1010_1100},
			nbits: 22,
			on:    22,
			out:   []byte{0b1011_0010, 0b1000_1101, 0b0001_0101},
			pbuf:  0b1000_0000,
			pn:    1,
		},
		{
			name:  "Offset 4 Unaligned",
			buf:   0b0101_0000,
			n:     4,
			p:     []byte{0b1001_0100, 0b0110_1000, 0b1010_1100},
			nbits: 22,
			on:    22,
			out:   []byte{0b0101_1001, 0b0100_0110, 0b1000_1010},
			pbuf:  0b1100_0000,
			pn:    2,
		},
		{
			name:  "Offset 5 Unaligned",
			buf:   0b1010_1000,
			n:     5,
			p:     []byte{0b1001_0100, 0b0110_1000, 0b1010_1100},
			nbits: 22,
			on:    22,
			out:   []byte{0b1010_1100, 0b1010_0011, 0b0100_0101},
			pbuf:  0b0110_0000,
			pn:    3,
		},
		{
			name:  "Offset 6 Unaligned",
			buf:   0b0101_0100,
			n:     6,
			p:     []byte{0b1001_0100, 0b0110_1000, 0b1010_1100},
			nbits: 22,
			on:    22,
			out:   []byte{0b0101_0110, 0b0101_0001, 0b1010_0010},
			pbuf:  0b1011_0000,
			pn:    4,
		},
		{
			name:  "Offset 7 Unaligned",
			buf:   0b1010_1010,
			n:     7,
			p:     []byte{0b1001_0100, 0b0110_1000, 0b1010_1100},
			nbits: 22,
			on:    22,
			out:   []byte{0b1010_1011, 0b0010_1000, 0b1101_0001},
			pbuf:  0b0101_1000,
			pn:    5,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewBuffer([]byte{})
			w := NewWriter(b)
			w.buf = tt.buf
			w.n = tt.n
			on, err := w.WriteBits(tt.p, tt.nbits)
			if err != nil {
				t.Fatal(err)
			}
			if on != tt.on {
				t.Errorf("on want %v; but got %v\n", tt.on, on)
			}
			out := b.Bytes()
			if !bytes.Equal(out, tt.out) {
				t.Errorf("out want %x; but got %x\n", tt.out, out)
			}
			if w.buf != tt.pbuf {
				t.Errorf("w.buf want %x; but got %x\n", tt.pbuf, w.buf)
			}
			if w.n != tt.pn {
				t.Errorf("w.n want %v; but got %v\n", tt.pn, w.n)
			}
		})
	}
}

type FixedBuffer struct {
	buf []byte
}

var ErrTooLarge = errors.New("too large")

func NewFixedBuffer(buf []byte) *FixedBuffer {
	return &FixedBuffer{buf: buf}
}

func (b *FixedBuffer) Bytes() []byte {
	return b.buf
}

func (b *FixedBuffer) Write(p []byte) (int, error) {
	n := len(p)
	l := len(b.buf)
	m := cap(b.buf) - l
	if n > m {
		n = m
	}
	b.buf = b.buf[:l+n]
	copy(b.buf[l:], p[:n])
	if n < len(p) {
		return n, ErrTooLarge
	}
	return n, nil
}

func TestWriter_Flush_Abnormal(t *testing.T) {
	t.Run("3bit", func(t *testing.T) {
		b := NewFixedBuffer([]byte{})
		w := NewWriter(b)
		w.buf = 0b1010_0000
		w.n = 3
		err := w.Flush()
		if err != ErrTooLarge {
			t.Errorf("err want %v; but got %v\n", ErrTooLarge, err)
		}
		out := b.Bytes()
		if !bytes.Equal(out, []byte{}) {
			t.Errorf("out want %x; but got %x\n", []byte{}, out)
		}
		if w.buf != 0b1010_0000 {
			t.Errorf("w.buf want %x; but got %x\n", 0b1010_0000, w.buf)
		}
		if w.n != 3 {
			t.Errorf("w.n want %v; but got %v\n", 1, w.n)
		}
	})
}

func TestWriter_WriteByte_Abnormal(t *testing.T) {
	t.Run("", func(t *testing.T) {
		b := NewFixedBuffer([]byte{})
		w := NewWriter(b)
		w.buf = 0b1000_0000
		w.n = 1
		c := byte(0xc5)
		err := w.WriteByte(c)
		if err != ErrTooLarge {
			t.Errorf("err want %v; but got %v\n", ErrTooLarge, err)
		}
		out := b.Bytes()
		if !bytes.Equal(out, []byte{}) {
			t.Errorf("out want %x; but got %x\n", []byte{}, out)
		}
		if w.buf != 0b1000_0000 {
			t.Errorf("w.buf want %x; but got %x\n", 0b1000_0000, w.buf)
		}
		if w.n != 1 {
			t.Errorf("w.n want %v; but got %v\n", 1, w.n)
		}
	})
}

func TestWriter_WriteBool_Abnormal(t *testing.T) {
	t.Run("", func(t *testing.T) {
		b := NewFixedBuffer([]byte{})
		w := NewWriter(b)
		w.buf = 0b1001_0010
		w.n = 7
		err := w.WriteBool(true)
		if err != ErrTooLarge {
			t.Errorf("err want %v; but got %v\n", ErrTooLarge, err)
		}
		out := b.Bytes()
		if !bytes.Equal(out, []byte{}) {
			t.Errorf("out want %x; but got %x\n", []byte{}, out)
		}
		if w.buf != 0b1001_0010 {
			t.Errorf("w.buf want %x; but got %x\n", 0b1001_0010, w.buf)
		}
		if w.n != 7 {
			t.Errorf("w.n want %v; but got %v\n", 7, w.n)
		}
	})
}

func TestWriter_Write_Abnormal(t *testing.T) {
	t.Run("buf 1B write 2byte", func(t *testing.T) {
		buf := make([]byte, 0, 1)
		b := NewFixedBuffer(buf)
		w := NewWriter(b)
		w.buf = 0b1000_0000
		w.n = 1
		p := []byte{0b1100_0101, 0b1010_0110}
		on, err := w.Write(p)
		if err != ErrTooLarge {
			t.Errorf("err want %v; but got %v\n", ErrTooLarge, err)
		}
		if on != 1 {
			t.Errorf("on want %v; but got %v\n", 1, on)
		}
		out := b.Bytes()
		if !bytes.Equal(out, []byte{0b1110_0010}) {
			t.Errorf("out want %x; but got %x\n", []byte{0b1110_0010}, out)
		}
		if w.buf != 0b1000_0000 {
			t.Errorf("w.buf want %x; but got %x\n", 0b1000_0000, w.buf)
		}
		if w.n != 1 {
			t.Errorf("w.n want %v; but got %v\n", 1, w.n)
		}
	})
}

func TestWriter_WriteBits_Abnormal(t *testing.T) {
	t.Run("buf 1B write 15bit", func(t *testing.T) {
		buf := make([]byte, 0, 1)
		b := NewFixedBuffer(buf)
		w := NewWriter(b)
		w.buf = 0b1000_0000
		w.n = 1
		p := []byte{0b1100_0101, 0b1010_0110}
		on, err := w.WriteBits(p, 15)
		if err != ErrTooLarge {
			t.Errorf("err want %v; but got %v\n", ErrTooLarge, err)
		}
		if on != 8 {
			t.Errorf("on want %v; but got %v\n", 8, on)
		}
		out := b.Bytes()
		if !bytes.Equal(out, []byte{0b1110_0010}) {
			t.Errorf("out want %x; but got %x\n", []byte{0b1110_0010}, out)
		}
		if w.buf != 0b1000_0000 {
			t.Errorf("w.buf want %x; but got %x\n", 0b1000_0000, w.buf)
		}
		if w.n != 1 {
			t.Errorf("w.n want %v; but got %v\n", 1, w.n)
		}
	})
	t.Run("buf 1B write 16bit", func(t *testing.T) {
		buf := make([]byte, 0, 1)
		b := NewFixedBuffer(buf)
		w := NewWriter(b)
		w.buf = 0b1000_0000
		w.n = 1
		p := []byte{0b1100_0101, 0b1010_0111}
		on, err := w.WriteBits(p, 16)
		if err != ErrTooLarge {
			t.Errorf("err want %v; but got %v\n", ErrTooLarge, err)
		}
		if on != 8 {
			t.Errorf("on want %v; but got %v\n", 8, on)
		}
		out := b.Bytes()
		if !bytes.Equal(out, []byte{0b1110_0010}) {
			t.Errorf("out want %x; but got %x\n", []byte{0b1110_0010}, out)
		}
		if w.buf != 0b1000_0000 {
			t.Errorf("w.buf want %x; but got %x\n", 0b1000_0000, w.buf)
		}
		if w.n != 1 {
			t.Errorf("w.n want %v; but got %v\n", 1, w.n)
		}
	})
}
