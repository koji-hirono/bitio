package bitio

import (
	"bytes"
	"io"
	"testing"
)

func TestReader_ReadByte(t *testing.T) {
	want := byte(0xcd)
	buf := []byte{0xcd}
	r := NewReader(buf)
	c, err := r.ReadByte()
	if err != nil {
		t.Fatal(err)
	}
	if c != want {
		t.Errorf("want %x; but got %x\n", want, c)
	}
}

func TestReader_ReadBool(t *testing.T) {
	want := true
	buf := []byte{0x80}
	r := NewReader(buf)
	c, err := r.ReadBool()
	if err != nil {
		t.Fatal(err)
	}
	if c != want {
		t.Errorf("want %v; but got %v\n", want, c)
	}
}

func TestReader_AtomicRead(t *testing.T) {
	want := []byte{0xcd, 0xef, 0x89, 0x31}
	buf := []byte{0xcd, 0xef, 0x89, 0x31}
	p := make([]byte, 4)
	r := NewReader(buf)
	err := r.AtomicRead(p)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(p, want) {
		t.Errorf("want %x; but got %x\n", want, p)
	}
}

func TestReader_Read(t *testing.T) {
	want := []byte{0xcd, 0xef, 0x89, 0x31}
	buf := []byte{0xcd, 0xef, 0x89, 0x31}
	p := make([]byte, 8)
	r := NewReader(buf)
	n, err := r.Read(p)
	if n != 4 {
		t.Errorf("want %v; but got %v\n", 4, n)
	}
	if !bytes.Equal(p[:n], want) {
		t.Errorf("want %x; but got %x\n", want, p[:n])
	}
	if err != nil && err != io.EOF {
		t.Errorf("want nil or EOF; but got %v\n", err)
	}
}

func TestReader_Read2(t *testing.T) {
	want := []byte{0xcd, 0xef, 0x89}
	buf := []byte{0xcd, 0xef, 0x89, 0x31}
	p := make([]byte, 3)
	r := NewReader(buf)
	n, err := r.Read(p)
	if n != 3 {
		t.Errorf("want %v; but got %v\n", 3, n)
	}
	if !bytes.Equal(p[:n], want) {
		t.Errorf("want %x; but got %x\n", want, p[:n])
	}
	if err != nil && err != io.EOF {
		t.Errorf("want nil or EOF; but got %v\n", err)
	}
}

func TestReader_ReadBits(t *testing.T) {
	want := []byte{0xcd, 0xef, 0x89, 0x34}
	buf := []byte{0xcd, 0xef, 0x89, 0x34}
	nbits := 30
	p := make([]byte, 4)
	r := NewReader(buf)
	err := r.ReadBits(p, nbits)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(p, want) {
		t.Errorf("want %x; but got %x\n", want, p)
	}
}

func TestReader_ReadBitField(t *testing.T) {
	want := uint64(0b0011_0011_0111_1011_1110_0010_0100_1100)
	nbits := 30
	buf := []byte{0xcd, 0xef, 0x89, 0x31}
	r := NewReader(buf)
	v, err := r.ReadBitField(nbits)
	if err != nil {
		t.Fatal(err)
	}
	if v != want {
		t.Errorf("want %x; but got %x\n", want, v)
	}
}

func TestReader_ReadBitField2(t *testing.T) {
	//              aaab bbbb    bbbb cccc    cccc cccc    cccc cccc
	buf := []byte{0b1100_1101, 0b1110_1111, 0b1000_1001, 0b0011_0001}
	r := NewReader(buf)
	t.Run("3bit", func(t *testing.T) {
		want := uint64(0b110)
		nbits := 3
		v, err := r.ReadBitField(nbits)
		if err != nil {
			t.Fatal(err)
		}
		if v != want {
			t.Errorf("want %x; but got %x\n", want, v)
		}
	})
	t.Run("9bit", func(t *testing.T) {
		want := uint64(0b0_1101_1110)
		nbits := 9
		v, err := r.ReadBitField(nbits)
		if err != nil {
			t.Fatal(err)
		}
		if v != want {
			t.Errorf("want %x; but got %x\n", want, v)
		}
	})
	t.Run("9bit", func(t *testing.T) {
		want := uint64(0b1111_1000_1001_0011_0001)
		nbits := 20
		v, err := r.ReadBitField(nbits)
		if err != nil {
			t.Fatal(err)
		}
		if v != want {
			t.Errorf("want %x; but got %x\n", want, v)
		}
	})
}

func TestReader_ReadBitField3(t *testing.T) {
	//              aaa_ ____    bbbb bbbb    b
	buf := []byte{0b1100_1101, 0b1110_1111, 0b1000_1001, 0b0011_0001}
	r := NewReader(buf)
	t.Run("3bit", func(t *testing.T) {
		want := uint64(0b110)
		nbits := 3
		v, err := r.ReadBitField(nbits)
		if err != nil {
			t.Fatal(err)
		}
		if v != want {
			t.Errorf("want %x; but got %x\n", want, v)
		}
	})
	r.Align()
	t.Run("9bit", func(t *testing.T) {
		want := uint64(0b1110_1111_1)
		nbits := 9
		v, err := r.ReadBitField(nbits)
		if err != nil {
			t.Fatal(err)
		}
		if v != want {
			t.Errorf("want %x; but got %x\n", want, v)
		}
	})
}

func TestReader_ReadBool_ErrUnexpectedEOF(t *testing.T) {
	buf := []byte{}
	r := NewReader(buf)
	_, err := r.ReadBool()
	if err != io.ErrUnexpectedEOF {
		t.Errorf("want %v; but got %v\n", io.ErrUnexpectedEOF, err)
	}
}

func TestReader_ReadBits_ErrUnexpectedEOF(t *testing.T) {
	buf := []byte{}
	r := NewReader(buf)
	p := make([]byte, 4)
	nbits := 28
	err := r.ReadBits(p, nbits)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("want %v; but got %v\n", io.ErrUnexpectedEOF, err)
	}
}

func TestReader_ReadBitField_ErrUnexpectedEOF(t *testing.T) {
	buf := []byte{}
	r := NewReader(buf)
	nbits := 28
	_, err := r.ReadBitField(nbits)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("want %v; but got %v\n", io.ErrUnexpectedEOF, err)
	}
}

func TestReader_AtomicRead_ErrUnexpectedEOF(t *testing.T) {
	buf := []byte{}
	r := NewReader(buf)
	p := make([]byte, 4)
	err := r.AtomicRead(p)
	if err != io.ErrUnexpectedEOF {
		t.Errorf("want %v; but got %v\n", io.ErrUnexpectedEOF, err)
	}
}
