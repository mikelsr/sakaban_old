package comm

import "encoding/binary"

// uint16FromBytes converts a byte slice into an unsigned 16 bit integer
// LITTLE ENDIAN
func uint16FromBytes(b []byte) uint16 {
	if len(b) < 2 {
		return uint16(0)
	}
	return binary.LittleEndian.Uint16(b[0:2])
}

// uint16ToBytes converts a 16 bit unsigned integer into a byte slice
// LITTLE ENDIAN
func uint16ToBytes(n uint16) []byte {
	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, n)
	return b
}
