package bitio

// Convert bit sequences into 64-bit unsigned integers.

// PutBitField encodes a uint64 into p.
// If the buffer is too small, PutBitField will panic.
func PutBitField(p []byte, v uint64, nbits int) {
	v <<= 64 - nbits
	nbytes := (nbits + 7) >> 3
	for i := 0; i < nbytes; i++ {
		k := 56 - (i << 3)
		p[i] = byte(v >> k)
	}
}

// BitField decodes a uint64 from p.
// If the buffer is too small, BitField will panic.
func BitField(p []byte, nbits int) uint64 {
	nbytes := (nbits + 7) >> 3
	v := uint64(0)
	for i := 0; i < nbytes; i++ {
		k := 56 - (i << 3)
		v |= uint64(p[i]) << k
	}
	v >>= 64 - nbits
	return v
}
