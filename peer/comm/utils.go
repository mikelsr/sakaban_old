package comm

import (
	"bufio"
	"encoding/binary"
	"errors"
)

// EmptyMessageFromMessageType returns an empty message given a message type
func EmptyMessageFromMessageType(msgType MessageType) (Message, error) {
	var msg Message
	switch msgType {
	case MTBlockContent:
		msg = &BlockContent{}
	case MTBlockRequest:
		msg = &BlockRequest{}
	case MTIndexContent:
		msg = &IndexContent{}
	case MTIndexRequest:
		msg = &IndexRequest{}
	default:
		return nil, errors.New("Invalid message type")
	}
	return msg, nil
}

// MessageTypeFromBytes reads the MessageType from the first element of a
// byte slice
func MessageTypeFromBytes(bytes []byte) (*MessageType, error) {
	if len(bytes) < 1 {
		return nil, errors.New("Invalid byte slice size")
	}

	messageType := MessageType(bytes[0])
	if minMessageType <= messageType && messageType <= maxMessageType {
		return &messageType, nil
	}
	return nil, errors.New("Unknown MessageType")
}

// RecvMessage reads all the content of a Message from a net.String
// UNTESTED
func RecvMessage(s *bufio.Reader, msg Message) ([]byte, error) {

	// FIXME: some test failed due to `make` blocking execution?!?!?!
	buf := make([]byte, bufferSize)

	// receive initial bytes
	n, err := s.Read(buf)
	if err != nil {
		return nil, err
	}

	// trim received data
	buf = buf[:n]
	// extract total size of the message
	msgSize := msg.Size(buf)

	// receive complete message
	for uint64(len(buf)) < msgSize {
		recv := make([]byte, bufferSize)
		n, err := s.Read(recv)
		if err != nil {
			return nil, err
		}
		recv = recv[:n]
		buf = append(buf, recv...)
	}
	return buf, nil
}

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
