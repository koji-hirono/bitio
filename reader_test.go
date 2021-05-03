package bitio

import (
	"bufio"
	"bytes"
	"io"
	"testing"
)

func TestReader_DiscardBuffer(t *testing.T) {
	cases := []struct {
		name string
		src  []byte
		buf  byte
		n    int
	}{
		{
			name: "buf 1B",
			src:  []byte{0xe5, 0xc6},
			buf:  0b1000_0000,
			n:    1,
		},
		{
			name: "buf 11B",
			src:  []byte{0xe5, 0xc6},
			buf:  0b1100_0000,
			n:    2,
		},
		{
			name: "buf 101B",
			src:  []byte{0xe5, 0xc6},
			buf:  0b1010_0000,
			n:    3,
		},
		{
			name: "buf 0101B",
			src:  []byte{0xe5, 0xc6},
			buf:  0b0101_0000,
			n:    4,
		},
		{
			name: "buf 10101B",
			src:  []byte{0xe5, 0xc6},
			buf:  0b1010_1000,
			n:    5,
		},
		{
			name: "buf 010101B",
			src:  []byte{0xe5, 0xc6},
			buf:  0b0101_0100,
			n:    6,
		},
		{
			name: "buf 1010101B",
			src:  []byte{0xe5, 0xc6},
			buf:  0b1010_1010,
			n:    7,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewReader(tt.src)
			r := NewReader(b)
			r.buf = tt.buf
			r.n = tt.n
			r.DiscardBuffer()
			if r.buf != 0 {
				t.Errorf("r.buf want 0; but got %x\n", r.buf)
			}
			if r.n != 0 {
				t.Errorf("r.n want 0; but got %v\n", r.n)
			}
		})
	}
}

func TestReader_ReadByte(t *testing.T) {
	cases := []struct {
		name string
		src  []byte
		buf  byte
		n    int
		c    byte
		pbuf byte
		pn   int
	}{
		{
			name: "buf empty",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0,
			n:    0,
			c:    0b1110_0101,
			pbuf: 0,
			pn:   0,
		},
		{
			name: "buf 0B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b0000_0000,
			n:    1,
			c:    0b0111_0010,
			pbuf: 0b1000_0000,
			pn:   1,
		},
		{
			name: "buf 1B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b1000_0000,
			n:    1,
			c:    0b1111_0010,
			pbuf: 0b1000_0000,
			pn:   1,
		},
		{
			name: "buf 00B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b0000_0000,
			n:    2,
			c:    0b0011_1001,
			pbuf: 0b0100_0000,
			pn:   2,
		},
		{
			name: "buf 01B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b0100_0000,
			n:    2,
			c:    0b0111_1001,
			pbuf: 0b0100_0000,
			pn:   2,
		},
		{
			name: "buf 10B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b1000_0000,
			n:    2,
			c:    0b1011_1001,
			pbuf: 0b0100_0000,
			pn:   2,
		},
		{
			name: "buf 11B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b1100_0000,
			n:    2,
			c:    0b1111_1001,
			pbuf: 0b0100_0000,
			pn:   2,
		},
		{
			name: "buf 101B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b1010_0000,
			n:    3,
			c:    0b1011_1100,
			pbuf: 0b1010_0000,
			pn:   3,
		},
		{
			name: "buf 0101B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b0101_0000,
			n:    4,
			c:    0b0101_1110,
			pbuf: 0b0101_0000,
			pn:   4,
		},
		{
			name: "buf 10101B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b1010_1000,
			n:    5,
			c:    0b1010_1111,
			pbuf: 0b0010_1000,
			pn:   5,
		},
		{
			name: "buf 010101B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b0101_0100,
			n:    6,
			c:    0b0101_0111,
			pbuf: 0b1001_0100,
			pn:   6,
		},
		{
			name: "buf 1010101B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b1010_1010,
			n:    7,
			c:    0b1010_1011,
			pbuf: 0b1100_1010,
			pn:   7,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewReader(tt.src)
			r := NewReader(b)
			r.buf = tt.buf
			r.n = tt.n
			c, err := r.ReadByte()
			if err != nil {
				t.Fatal(err)
			}
			if c != tt.c {
				t.Errorf("c want %x; but got %x\n", tt.c, c)
			}
			if r.buf != tt.pbuf {
				t.Errorf("r.buf want %x; but got %x\n", tt.pbuf, r.buf)
			}
			if r.n != tt.pn {
				t.Errorf("r.n want %v; but got %v\n", tt.pn, r.n)
			}
		})
	}

	t.Run("0bit EOF", func(t *testing.T) {
		b := bytes.NewReader([]byte{})
		r := NewReader(b)
		r.buf = 0
		r.n = 0
		_, err := r.ReadByte()
		if err != io.EOF {
			t.Errorf("want io.EOF; but got %v\n", err)
		}
		if r.buf != 0 {
			t.Errorf("r.buf want 0; but got %x\n", r.buf)
		}
		if r.n != 0 {
			t.Errorf("r.n want 0; but got %v\n", r.n)
		}
	})
	t.Run("4bit EOF", func(t *testing.T) {
		b := bytes.NewReader([]byte{})
		r := NewReader(b)
		r.buf = 0b0101_0000
		r.n = 4
		_, err := r.ReadByte()
		if err != io.EOF {
			t.Errorf("want io.EOF; but got %v\n", err)
		}
		if r.buf != 0b0101_0000 {
			t.Errorf("r.buf want 0b0101_0000; but got %x\n", r.buf)
		}
		if r.n != 4 {
			t.Errorf("r.n want 4; but got %v\n", r.n)
		}
	})
	t.Run("7bit EOF", func(t *testing.T) {
		b := bytes.NewReader([]byte{})
		r := NewReader(b)
		r.buf = 0b1100_1010
		r.n = 7
		_, err := r.ReadByte()
		if err != io.EOF {
			t.Errorf("want io.EOF; but got %v\n", err)
		}
		if r.buf != 0b1100_1010 {
			t.Errorf("r.buf want 0b1100_1010; but got %x\n", r.buf)
		}
		if r.n != 7 {
			t.Errorf("r.n want 7; but got %v\n", r.n)
		}
	})
}

func TestReader_ReadBool(t *testing.T) {
	cases := []struct {
		name string
		src  []byte
		buf  byte
		n    int
		b    bool
		pbuf byte
		pn   int
	}{
		{
			name: "src empty buf 1010011B",
			src:  []byte{},
			buf:  0b1010_0110,
			n:    7,
			b:    true,
			pbuf: 0b0100_1100,
			pn:   6,
		},
		{
			name: "src empty buf 1B",
			src:  []byte{},
			buf:  0b1000_0000,
			n:    1,
			b:    true,
			pbuf: 0,
			pn:   0,
		},
		{
			name: "buf empty",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0,
			n:    0,
			b:    true,
			pbuf: 0b1100_1010,
			pn:   7,
		},
		{
			name: "buf 0B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b0000_0000,
			n:    1,
			b:    false,
			pbuf: 0,
			pn:   0,
		},
		{
			name: "buf 1B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b1000_0000,
			n:    1,
			b:    true,
			pbuf: 0,
			pn:   0,
		},
		{
			name: "buf 00B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b0000_0000,
			n:    2,
			b:    false,
			pbuf: 0b0000_0000,
			pn:   1,
		},
		{
			name: "buf 01B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b0100_0000,
			n:    2,
			b:    false,
			pbuf: 0b1000_0000,
			pn:   1,
		},
		{
			name: "buf 10B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b1000_0000,
			n:    2,
			b:    true,
			pbuf: 0b0000_0000,
			pn:   1,
		},
		{
			name: "buf 11B",
			src:  []byte{0b1110_0101, 0b1100_0110},
			buf:  0b1100_0000,
			n:    2,
			b:    true,
			pbuf: 0b1000_0000,
			pn:   1,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewReader(tt.src)
			r := NewReader(b)
			r.buf = tt.buf
			r.n = tt.n
			bb, err := r.ReadBool()
			if err != nil {
				t.Fatal(err)
			}
			if bb != tt.b {
				t.Errorf("b want %v; but got %v\n", tt.b, bb)
			}
			if r.buf != tt.pbuf {
				t.Errorf("r.buf want %x; but got %x\n", tt.pbuf, r.buf)
			}
			if r.n != tt.pn {
				t.Errorf("r.n want %v; but got %v\n", tt.pn, r.n)
			}
		})
	}

	t.Run("0bit EOF", func(t *testing.T) {
		b := bytes.NewReader([]byte{})
		r := NewReader(b)
		r.buf = 0
		r.n = 0
		_, err := r.ReadBool()
		if err != io.EOF {
			t.Errorf("want io.EOF; but got %v\n", err)
		}
		if r.buf != 0 {
			t.Errorf("r.buf want 0; but got %x\n", r.buf)
		}
		if r.n != 0 {
			t.Errorf("r.n want 0; but got %v\n", r.n)
		}
	})
}

func TestReader_Read(t *testing.T) {
	cases := []struct {
		name     string
		src      []byte
		buf      byte
		n        int
		nbytes   int
		p        []byte
		on       int
		pbuf     byte
		pn       int
		allowEof bool
	}{
		{
			name: "buf 0bit read 4byte",
			src: []byte{
				0b1110_0000, 0b1100_0110, 0b0111_0011, 0b0001_0101,
			},
			buf:    0,
			n:      0,
			nbytes: 4,
			p: []byte{
				0b1110_0000, 0b1100_0110, 0b0111_0011, 0b0001_0101,
			},
			on:   4,
			pbuf: 0,
			pn:   0,
		},
		{
			name: "buf 0bit read 5byte",
			src: []byte{
				0b1110_0000, 0b1100_0110, 0b0111_0011, 0b0001_0101,
			},
			buf:    0,
			n:      0,
			nbytes: 5,
			p: []byte{
				0b1110_0000, 0b1100_0110, 0b0111_0011, 0b0001_0101,
			},
			on:       4,
			pbuf:     0,
			pn:       0,
			allowEof: true,
		},
		{
			name: "buf 1bit read 4byte",
			src: []byte{
				0b1110_0000, 0b1100_0110, 0b0111_0011, 0b0001_0101,
			},
			buf:    0b1000_0000,
			n:      1,
			nbytes: 4,
			p: []byte{
				0b1111_0000, 0b0110_0011, 0b0011_1001, 0b1000_1010,
			},
			on:   4,
			pbuf: 0b1000_0000,
			pn:   1,
		},
		{
			name: "buf 1bit read 5byte",
			src: []byte{
				0b1110_0000, 0b1100_0110, 0b0111_0011, 0b0001_0101,
			},
			buf:    0b1000_0000,
			n:      1,
			nbytes: 5,
			p: []byte{
				0b1111_0000, 0b0110_0011, 0b0011_1001, 0b1000_1010,
			},
			on:       4,
			pbuf:     0b1000_0000,
			pn:       1,
			allowEof: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewReader(tt.src)
			r := NewReader(b)
			r.buf = tt.buf
			r.n = tt.n
			p := make([]byte, tt.nbytes)
			on, err := r.Read(p)
			if on != tt.on {
				t.Errorf("on want %v; but got %v\n", tt.on, on)
			}
			if !bytes.Equal(p[:on], tt.p) {
				t.Errorf("p want %x; but got %x\n", tt.p, p[:on])
			}
			if err != nil {
				if err != io.EOF || !tt.allowEof {
					t.Fatal(err)
				}
			}
			if r.buf != tt.pbuf {
				t.Errorf("r.buf want %x; but got %x\n", tt.pbuf, r.buf)
			}
			if r.n != tt.pn {
				t.Errorf("r.n want %v; but got %v\n", tt.pn, r.n)
			}
		})
	}
}

func TestReader_ReadBits(t *testing.T) {
	cases := []struct {
		name     string
		src      []byte
		buf      byte
		n        int
		p        []byte
		nbits    int
		on       int
		pbuf     byte
		pn       int
		allowEof bool
	}{
		{
			name:  "buf 0bit read 1bit",
			src:   []byte{0b1110_0101, 0b1100_0110},
			buf:   0,
			n:     0,
			p:     []byte{0b1000_0000},
			nbits: 1,
			on:    1,
			pbuf:  0b1100_1010,
			pn:    7,
		},
		{
			name:  "buf 0bit read 2bit",
			src:   []byte{0b1110_0101, 0b1100_0110},
			buf:   0,
			n:     0,
			p:     []byte{0b1100_0000},
			nbits: 2,
			on:    2,
			pbuf:  0b1001_0100,
			pn:    6,
		},
		{
			name:  "buf 0bit read 3bit",
			src:   []byte{0b1110_0101, 0b1100_0110},
			buf:   0,
			n:     0,
			p:     []byte{0b1110_0000},
			nbits: 3,
			on:    3,
			pbuf:  0b0010_1000,
			pn:    5,
		},
		{
			name:  "buf 0bit read 4bit",
			src:   []byte{0b1110_0101, 0b1100_0110},
			buf:   0,
			n:     0,
			p:     []byte{0b1110_0000},
			nbits: 4,
			on:    4,
			pbuf:  0b0101_0000,
			pn:    4,
		},
		{
			name:  "buf 0bit read 5bit",
			src:   []byte{0b1110_0101, 0b1100_0110},
			buf:   0,
			n:     0,
			p:     []byte{0b1110_0000},
			nbits: 5,
			on:    5,
			pbuf:  0b1010_0000,
			pn:    3,
		},
		{
			name:  "buf 0bit read 6bit",
			src:   []byte{0b1110_0101, 0b1100_0110},
			buf:   0,
			n:     0,
			p:     []byte{0b1110_0100},
			nbits: 6,
			on:    6,
			pbuf:  0b0100_0000,
			pn:    2,
		},
		{
			name:  "buf 0bit read 7bit",
			src:   []byte{0b1110_0101, 0b1100_0110},
			buf:   0,
			n:     0,
			p:     []byte{0b1110_0100},
			nbits: 7,
			on:    7,
			pbuf:  0b1000_0000,
			pn:    1,
		},
		{
			name:  "buf 0bit read 8bit",
			src:   []byte{0b1110_0101, 0b1100_0110},
			buf:   0,
			n:     0,
			p:     []byte{0b1110_0101},
			nbits: 8,
			on:    8,
			pbuf:  0,
			pn:    0,
		},
		{
			name:     "buf 0bit read 16bit",
			src:      []byte{0b1110_0101, 0b1100_0110},
			buf:      0,
			n:        0,
			p:        []byte{0b1110_0101, 0b1100_0110},
			nbits:    16,
			on:       16,
			pbuf:     0,
			pn:       0,
			allowEof: true,
		},
		{
			name:     "buf 0bit read 17bit",
			src:      []byte{0b1110_0101, 0b1100_0110},
			buf:      0,
			n:        0,
			p:        []byte{0b1110_0101, 0b1100_0110},
			nbits:    17,
			on:       16,
			pbuf:     0,
			pn:       0,
			allowEof: true,
		},
		{
			name:     "buf 0bit read 32bit",
			src:      []byte{0b1110_0101, 0b1100_0110},
			buf:      0,
			n:        0,
			p:        []byte{0b1110_0101, 0b1100_0110},
			nbits:    32,
			on:       16,
			pbuf:     0,
			pn:       0,
			allowEof: true,
		},
		{
			name:  "buf 1bit read 16bit",
			src:   []byte{0b1001_0100, 0b0110_1000},
			buf:   0b1000_0000,
			n:     1,
			p:     []byte{0b1100_1010, 0b0011_0100},
			nbits: 16,
			on:    16,
			pbuf:  0b0000_0000,
			pn:    1,
		},
		{
			name:  "buf 5bit read 20bit",
			src:   []byte{0b1001_0100, 0b0110_1001},
			buf:   0b1010_1000,
			n:     5,
			p:     []byte{0b1010_1100, 0b1010_0011, 0b0100_0000},
			nbits: 20,
			on:    20,
			pbuf:  0b1000_0000,
			pn:    1,
		},
		{
			name:     "buf 6bit read 22bit",
			src:      []byte{0b1001_0100, 0b0110_1000},
			buf:      0b1010_1100,
			n:        6,
			p:        []byte{0b1010_1110, 0b0101_0001, 0b1010_0000},
			nbits:    22,
			on:       22,
			pbuf:     0,
			pn:       0,
			allowEof: true,
		},
		{
			name:     "buf 6bit read 23bit",
			src:      []byte{0b1001_0100, 0b0110_1000},
			buf:      0b1010_1100,
			n:        6,
			p:        []byte{0b1010_1110, 0b0101_0001, 0b1010_0000},
			nbits:    23,
			on:       22,
			pbuf:     0,
			pn:       0,
			allowEof: true,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			b := bytes.NewReader(tt.src)
			r := NewReader(b)
			r.buf = tt.buf
			r.n = tt.n
			nbytes := (tt.nbits + 7) >> 3
			p := make([]byte, nbytes)
			on, err := r.ReadBits(p, tt.nbits)
			if err != nil {
				if err != io.EOF || !tt.allowEof {
					t.Fatal(err)
				}
			}
			if on != tt.on {
				t.Errorf("on want %v; but got %v\n", tt.on, on)
			}
			onbytes := (on + 7) >> 3
			if !bytes.Equal(p[:onbytes], tt.p) {
				t.Errorf("p want %x; but got %x\n", tt.p, p[:onbytes])
			}
			if r.buf != tt.pbuf {
				t.Errorf("r.buf want %x; but got %x\n", tt.pbuf, r.buf)
			}
			if r.n != tt.pn {
				t.Errorf("r.n want %v; but got %v\n", tt.pn, r.n)
			}
		})
	}

	t.Run("empty", func(t *testing.T) {
		b := bytes.NewReader([]byte{})
		r := NewReader(b)
		nbits := 1
		nbytes := (nbits + 7) >> 3
		p := make([]byte, nbytes)
		on, err := r.ReadBits(p, nbits)
		if err != io.EOF {
			t.Errorf("err want io.EOF; but got %v\n", err)
		}
		if on != 0 {
			t.Errorf("on want 0; but got %v\n", on)
		}
		onbytes := (on + 7) >> 3
		if !bytes.Equal(p[:onbytes], []byte{}) {
			t.Errorf("p want empty; but got %x\n", p[:onbytes])
		}
		if r.buf != 0 {
			t.Errorf("r.buf want 0; but got %x\n", r.buf)
		}
		if r.n != 0 {
			t.Errorf("r.n want 0; but got %v\n", r.n)
		}
	})
}

func TestReader_ReadBits_FromLimitedBuffer(t *testing.T) {
	data := []byte{
		0x11, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77, 0x88,
		0x99, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff, 0x0f,
		0x0e, 0x0d, 0x0c, 0x0b, 0x0a, 0x09,
	}
	b := bytes.NewReader(data)
	bs := bufio.NewReaderSize(b, 16)
	r := NewReader(bs)
	t.Run("ReadByte", func(t *testing.T) {
		got, err := r.ReadByte()
		if err != nil {
			t.Fatal(err)
		}
		want := byte(0x11)
		if got != want {
			t.Errorf("want %x; but got %x\n", want, got)
		}
	})
	t.Run("ReadBits", func(t *testing.T) {
		want := data[1:]
		got := make([]byte, len(want))
		n, err := r.ReadBits(got, len(want)*8)
		if err != nil {
			t.Fatal(err)
		}
		if n > len(want)*8 {
			t.Errorf("want <= %v; but got %v\n", len(want)*8, n)
		}
		nbytes := (n + 7) >> 3
		if !bytes.Equal(got[:nbytes], want[:nbytes]) {
			t.Errorf("want %x; but got %x\n", want[:nbytes], got[:nbytes])
		}
	})
}
