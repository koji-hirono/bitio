package bitio

import (
	"bytes"
	"testing"

	"github.com/koji-hirono/memio"
)

func TestWriter_WriteByte(t *testing.T) {
	t.Run("1byte", func(t *testing.T) {
		c := byte(0xcd)
		want := []byte{0xcd}
		g := memio.NewVar(nil)
		w := NewWriter(g)
		err := w.WriteByte(c)
		if err != nil {
			t.Fatal(err)
		}
		out := w.Bytes()
		if !bytes.Equal(out, want) {
			t.Errorf("want %x; but got %x\n", want, out)
		}
	})
	t.Run("out of memory", func(t *testing.T) {
		c := byte(0xcd)
		g := memio.NewFixed(nil)
		w := NewWriter(g)
		err := w.WriteByte(c)
		if err != memio.ErrNoMem {
			t.Errorf("want %v; but got %v\n", memio.ErrNoMem, err)
		}
	})
}

func TestWriter_WriteBool(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		b := true
		want := []byte{0x80}
		g := memio.NewVar(nil)
		w := NewWriter(g)
		err := w.WriteBool(b)
		if err != nil {
			t.Fatal(err)
		}
		out := w.Bytes()
		if !bytes.Equal(out, want) {
			t.Errorf("want %x; but got %x\n", want, out)
		}
	})
	t.Run("false", func(t *testing.T) {
		b := false
		want := []byte{0x00}
		g := memio.NewVar(nil)
		w := NewWriter(g)
		err := w.WriteBool(b)
		if err != nil {
			t.Fatal(err)
		}
		out := w.Bytes()
		if !bytes.Equal(out, want) {
			t.Errorf("want %x; but got %x\n", want, out)
		}
	})
	t.Run("out of memory", func(t *testing.T) {
		b := true
		g := memio.NewFixed(nil)
		w := NewWriter(g)
		err := w.WriteBool(b)
		if err != memio.ErrNoMem {
			t.Errorf("want %v; but got %v\n", memio.ErrNoMem, err)
		}
	})
}

func TestWriter_Write(t *testing.T) {
	t.Run("8byte", func(t *testing.T) {
		p := []byte{0xcd, 0xef, 0x8b, 0x76, 0x13, 0x54}
		want := []byte{0xcd, 0xef, 0x8b, 0x76, 0x13, 0x54}
		g := memio.NewVar(nil)
		w := NewWriter(g)
		_, err := w.Write(p)
		if err != nil {
			t.Fatal(err)
		}
		out := w.Bytes()
		if !bytes.Equal(out, want) {
			t.Errorf("want %x; but got %x\n", want, out)
		}
	})
	t.Run("out of memory", func(t *testing.T) {
		p := []byte{0xcd, 0xef, 0x8b, 0x76, 0x13, 0x54}
		g := memio.NewFixed(nil)
		w := NewWriter(g)
		_, err := w.Write(p)
		if err != memio.ErrNoMem {
			t.Errorf("want %v; but got %v\n", memio.ErrNoMem, err)
		}
	})
}

func TestWriter_WriteBits(t *testing.T) {
	t.Run("46bit", func(t *testing.T) {
		p := []byte{0xcd, 0xef, 0x8b, 0x76, 0x13, 0x54}
		nbits := 46
		want := []byte{0xcd, 0xef, 0x8b, 0x76, 0x13, 0x54}
		g := memio.NewVar(nil)
		w := NewWriter(g)
		err := w.WriteBits(p, nbits)
		if err != nil {
			t.Fatal(err)
		}
		out := w.Bytes()
		if !bytes.Equal(out, want) {
			t.Errorf("want %x; but got %x\n", want, out)
		}
	})
	t.Run("out of memory", func(t *testing.T) {
		p := []byte{0xcd, 0xef, 0x8b, 0x76, 0x13, 0x54}
		nbits := 46
		g := memio.NewFixed(nil)
		w := NewWriter(g)
		err := w.WriteBits(p, nbits)
		if err != memio.ErrNoMem {
			t.Errorf("want %v; but got %v\n", memio.ErrNoMem, err)
		}
	})
}

func TestWriter_WriteBitField(t *testing.T) {
	v := uint64(1)
	nbits := 1
	want := []byte{0x80}
	g := memio.NewVar(nil)
	w := NewWriter(g)
	err := w.WriteBitField(v, nbits)
	if err != nil {
		t.Fatal(err)
	}
	out := w.Bytes()
	if !bytes.Equal(out, want) {
		t.Errorf("want %x; but got %x\n", want, out)
	}
}

func TestWriter_WriteBitField2(t *testing.T) {
	g := memio.NewVar(nil)
	w := NewWriter(g)
	t.Run("1bit", func(t *testing.T) {
		v := uint64(1)
		nbits := 1
		err := w.WriteBitField(v, nbits)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("3bit", func(t *testing.T) {
		v := uint64(5)
		nbits := 3
		err := w.WriteBitField(v, nbits)
		if err != nil {
			t.Fatal(err)
		}
	})
	want := []byte{0b1101_0000}
	out := w.Bytes()
	if !bytes.Equal(out, want) {
		t.Errorf("want %x; but got %x\n", want, out)
	}
}

func TestWriter_WriteBitField3(t *testing.T) {
	g := memio.NewVar(nil)
	w := NewWriter(g)
	t.Run("6bit", func(t *testing.T) {
		v := uint64(0b0011_0101)
		nbits := 6
		err := w.WriteBitField(v, nbits)
		if err != nil {
			t.Fatal(err)
		}
	})
	t.Run("3bit", func(t *testing.T) {
		v := uint64(0b0000_0111)
		nbits := 3
		err := w.WriteBitField(v, nbits)
		if err != nil {
			t.Fatal(err)
		}
	})
	want := []byte{0b1101_0111, 0b1000_0000}
	out := w.Bytes()
	if !bytes.Equal(out, want) {
		t.Errorf("want %x; but got %x\n", want, out)
	}
}

func TestWriter_WriteBitField4(t *testing.T) {
	g := memio.NewVar(nil)
	w := NewWriter(g)
	t.Run("6bit", func(t *testing.T) {
		v := uint64(0b0011_0101)
		nbits := 6
		err := w.WriteBitField(v, nbits)
		if err != nil {
			t.Fatal(err)
		}
	})
	w.Align()
	t.Run("3bit", func(t *testing.T) {
		v := uint64(0b0000_0111)
		nbits := 3
		err := w.WriteBitField(v, nbits)
		if err != nil {
			t.Fatal(err)
		}
	})
	want := []byte{0b1101_0100, 0b1110_0000}
	out := w.Bytes()
	if !bytes.Equal(out, want) {
		t.Errorf("want %x; but got %x\n", want, out)
	}
}

func TestWriter_Grow_ErrNoMem(t *testing.T) {
	buf := make([]byte, 1)
	g := memio.NewFixed(buf)
	w := NewWriter(g)
	nbits := 16
	err := w.Grow(nbits)
	if err != memio.ErrNoMem {
		t.Errorf("wantErr %v; but got %v\n", memio.ErrNoMem, err)
	}
}
