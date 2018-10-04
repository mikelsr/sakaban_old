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

// uint64FromBytes converts a byte slice into an unsigned 64 bit integer
// LITTLE ENDIAN
func uint64FromBytes(b []byte) uint64 {
	if len(b) < 8 {
		return uint64(0)
	}
	return binary.LittleEndian.Uint64(b[0:8])
}

// uint64ToBytes converts a 64 bit unsigned integer into a byte slice
// LITTLE ENDIAN
func uint64ToBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, n)
	return b
}
