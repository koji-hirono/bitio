package bitio_test

import (
	"fmt"
	"math/bits"

	"github.com/koji-hirono/bitio"
	"github.com/koji-hirono/memio"
)

// ASN.1 Aligned/Unaligned PER

// 11.5 Encoding of a constrained whole number
func encodeConstrainedWholeNumber(w *bitio.Writer, n, lb, ub int64, aligned bool) {
	x := uint64(n - lb)
	R := uint64(ub - lb + 1)
	switch {
	case R == 1:
		// empty
	case !aligned || R < 256:
		// Encoding as a non-negative-binary-integer
		nbits := bits.Len64(R)
		w.WriteBitField(x, nbits)
	case R == 256:
		w.WriteBitField(x, 8)
	case R < 65536:
		w.WriteBitField(x, 16)
	default:
		L := (bits.Len64(x) + 7) >> 3
		encodeLength(w, uint64(L), 0, nil, aligned)
		// Encoding as a non-negative-binary-integer
		nbits := L << 3
		w.WriteBitField(x, nbits)
	}
}

// 11.6 Encoding of a normally small non-negative whole number
func encodeSmallNonNegativeWholeNumber(w *bitio.Writer, n uint64, aligned bool) {
	if n < 64 {
		w.WriteBool(false)
		w.WriteBitField(n, 6)
	} else {
		w.WriteBool(true)
		encodeSemiConstrainedWholeNumber(w, int64(n), 0, aligned)
	}
}

// 11.7 Encoding of a semi-constrained whole number
func encodeSemiConstrainedWholeNumber(w *bitio.Writer, n, lb int64, aligned bool) {
	x := uint64(n - lb)
	L := (bits.Len64(x) + 7) >> 3
	encodeLength(w, uint64(L), 0, nil, aligned)

	// Encoding as a non-negative-binary-integer
	nbits := L << 3
	w.WriteBitField(x, nbits)
}

// 11.8 Encoding of an unconstrained whole number
func encodeUnconstrainedWholeNumber(w *bitio.Writer, n int64, aligned bool) {
	x := uint64(n)
	if n < 0 {
		x = ^x
	}
	L := (bits.Len64(x) + 7) >> 3
	encodeLength(w, uint64(L), 0, nil, aligned)

	// Encoding as a 2's-complement-binary-integer
	nbits := L << 3
	w.WriteBitField(x, nbits)
}

// 11.9 General rules for encodinga length determinant
func encodeLength(w *bitio.Writer, L, lb uint64, ub *uint64, aligned bool) uint64 {
	if ub != nil && *ub < 65536 && lb == *ub {
		return 0
	}

	if !aligned {
		if ub != nil && *ub < 65536 {
			x := L - lb
			R := *ub - lb + 1
			nbits := bits.Len64(R)
			w.WriteBitField(x, nbits)
			return 0
		}
	}

	if ub != nil && *ub < 65536 {
		encodeConstrainedWholeNumber(w, int64(L), int64(lb), int64(*ub), aligned)
		return 0
	}

	switch {
	case L >= 4*16*1024:
		// 1100 0100
		w.WriteByte(0b1100_0100)
		return L - 4*16*1024
	case L >= 3*16*1024:
		// 1100 0011
		w.WriteByte(0b1100_0011)
		return L - 3*16*1024
	case L >= 2*16*1024:
		// 1100 0010
		w.WriteByte(0b1100_0010)
		return L - 2*16*1024
	case L >= 1*16*1024:
		// 1100 0001
		w.WriteByte(0b1100_0001)
		return L - 1*16*1024
	case L >= 128:
		// 10LL LLLL LLLL LLLL
		w.WriteByte(byte(L>>8) | 1<<7)
		w.WriteByte(byte(L))
		return 0
	default:
		// 0LLL LLLL
		w.WriteByte(byte(L))
		return 0
	}
}

// 13 Encoding the integer type
func encodeInteger(w *bitio.Writer, n int64, lb *int64, ub *int64, ext, aligned bool) {
	if ext {
		if (lb != nil && n < *lb) || (ub != nil && n >= *ub) {
			w.WriteBool(true)
			encodeUnconstrainedWholeNumber(w, n, aligned)
			return
		} else {
			w.WriteBool(false)
		}
	}

	switch {
	case lb != nil && ub != nil:
		encodeConstrainedWholeNumber(w, n, *lb, *ub, aligned)
	case lb != nil:
		encodeSemiConstrainedWholeNumber(w, n, *lb, aligned)
	default:
		encodeUnconstrainedWholeNumber(w, n, aligned)
	}
}

// 16 Encoding the bitstring type
func encodeBitString(w *bitio.Writer, p []byte, n, lb uint64, ub *uint64, ext, aligned bool) {
	if ext {
		if n < lb || (ub != nil && n >= *ub) {
			w.WriteBool(true)
		} else {
			w.WriteBool(false)
		}
	}

	if ub != nil && *ub == lb {
		if n <= 16 {
			w.WriteBits(p, int(*ub))
			return
		}
		if n < 64*1024 {
			w.WriteBits(p, int(n))
			if aligned {
				w.Align()
			}
			return
		}
	}

	encodeLength(w, n, lb, ub, aligned)
	w.WriteBits(p, int(n))
}

// TS 38.473 F1AP F1AP-IEs
//
// SCellIndex ::= INTEGER (1..31, ...)
//
func Example_encodeSCellIndex() {
	g := memio.NewVar([]byte{})
	w := bitio.NewWriter(g)

	n := int64(23)
	lb := int64(1)
	ub := int64(31)
	encodeInteger(w, n, &lb, &ub, true, true)

	fmt.Printf("%x\n", w.Bytes())

	// Output:
	// 58
}

// TS 38.331 RRC NR-RRC-Definitions
//
// NG-5G-S-TMSI ::= BIT STRING (SIZE (48))
//
func Example_encodeNG5GSTMSI() {
	g := memio.NewVar([]byte{})
	w := bitio.NewWriter(g)

	lb := uint64(48)
	ub := uint64(48)
	p := []byte{0x11, 0x22, 0x33, 0x44, 0x55, 0x66}
	nbits := uint64(48)

	encodeBitString(w, p, nbits, lb, &ub, false, false)

	fmt.Printf("%x\n", w.Bytes())

	// Output:
	// 112233445566
}
