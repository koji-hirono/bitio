package bitio

import (
	"bytes"
	"testing"
)

var bitFieldTestCases = []struct {
	name  string
	v     uint64
	nbits int
	b     []byte
}{
	{
		name:  "1bit",
		v:     0b1,
		nbits: 1,
		b:     []byte{0b1000_0000},
	},
	{
		name:  "2bit",
		v:     0b01,
		nbits: 2,
		b:     []byte{0b0100_0000},
	},
	{
		name:  "3bit",
		v:     0b101,
		nbits: 3,
		b:     []byte{0b1010_0000},
	},
	{
		name:  "4bit",
		v:     0b0110,
		nbits: 4,
		b:     []byte{0b0110_0000},
	},
	{
		name:  "5bit",
		v:     0b10111,
		nbits: 5,
		b:     []byte{0b1011_1000},
	},
	{
		name:  "6bit",
		v:     0b101101,
		nbits: 6,
		b:     []byte{0b1011_0100},
	},
	{
		name:  "7bit",
		v:     0b1011001,
		nbits: 7,
		b:     []byte{0b1011_0010},
	},
	{
		name:  "8bit",
		v:     0b10110011,
		nbits: 8,
		b:     []byte{0b1011_0011},
	},
	{
		name:  "9bit",
		v:     0b101010011,
		nbits: 9,
		b:     []byte{0b1010_1001, 0b1000_0000},
	},
	{
		name:  "16bit",
		v:     0b1010_1001_1100_0110,
		nbits: 16,
		b:     []byte{0b1010_1001, 0b1100_0110},
	},
	{
		name:  "32bit",
		v:     0xc1b2a3d4,
		nbits: 32,
		b:     []byte{0xc1, 0xb2, 0xa3, 0xd4},
	},
	{
		name:  "64bit",
		v:     0xc1b2a3d4_e590f687,
		nbits: 64,
		b:     []byte{0xc1, 0xb2, 0xa3, 0xd4, 0xe5, 0x90, 0xf6, 0x87},
	},
}

func TestPutBitField(t *testing.T) {
	for _, tt := range bitFieldTestCases {
		t.Run(tt.name, func(t *testing.T) {
			nbytes := (tt.nbits + 7) >> 3
			b := make([]byte, nbytes)
			PutBitField(b, tt.v, tt.nbits)
			if !bytes.Equal(b, tt.b) {
				t.Errorf("want %x; but got %x\n", tt.b, b)
			}
		})
	}
}

func TestBitField(t *testing.T) {
	for _, tt := range bitFieldTestCases {
		t.Run(tt.name, func(t *testing.T) {
			v := BitField(tt.b, tt.nbits)
			if v != tt.v {
				t.Errorf("want %x; but got %x\n", tt.v, v)
			}
		})
	}
}
