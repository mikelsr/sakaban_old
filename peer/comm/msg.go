package comm

import (
	"encoding/json"
	"errors"
	"fmt"

	"bitbucket.org/mikelsr/sakaban/fs"
	net "github.com/libp2p/go-libp2p-net"
	"github.com/satori/go.uuid"
)

// MessageType is used to identify the message type using one byte
type MessageType byte

// Message ensures anything sent/received by peers can be dumped to or loaded
// from a byte array
type Message interface {
	Dump() []byte
	Load([]byte) error
	Type() MessageType
}

// LongMessage has the necessary information to receive a full, long message
type LongMessage interface {
	Size([]byte) uint64
	Recv(s net.Stream) ([]byte, error)
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

/* Generic functions */

/* recvLongMessage reads all the content of a LongMessage from a net.String */
func recvLongMessage(s net.Stream, lm LongMessage) ([]byte, error) {
	buf := make([]byte, bufferSize)
	// receive initial bytes
	n, err := s.Read(buf)
	if err != nil {
		return nil, err
	}
	// trim received data
	buf = buf[:n]
	// extract total size of the message
	msgSize := lm.Size(buf)

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

/* Block content */

// BlockContent is used to send a fs.Block and some context, e.g. the file it
// belongs to
type BlockContent struct {
	MessageSize uint64    // total size of the message
	BlockN      uint8     // block number
	BlockSize   uint16    // block size in kB
	Content     []byte    // content of the block
	FileID      uuid.UUID // ID of the file the Block belongs to
}

// Dump creates a byte array: {MessageType, BlockNumber, BlockSize, FileID,
// Content} (28B + BlockSize * 1024B)
func (bc BlockContent) Dump() []byte {
	// create base message
	dump := append([]byte{byte(bc.Type()), byte(bc.BlockN)}, uint16ToBytes(bc.BlockSize)...)
	dump = append(dump, bc.FileID.Bytes()...)
	dump = append(dump, bc.Content...)
	// calculate size of message
	totalLen := uint64(len(dump) + 8)
	// create new massage merging the base message and its size
	cDump := []byte{dump[0]}
	cDump = append(cDump, uint64ToBytes(totalLen)...)
	cDump = append(cDump, dump[1:]...)
	return cDump
}

// Load reads blockN, blockSize, fileID, content from a byte slice created
// by br.Dump()
func (bc *BlockContent) Load(msg []byte) error {
	// parse contents of the message, extract values
	if len(msg) < 28 || MessageType(msg[0]) != MTBlockContent {
		return errors.New("Invalid message type")
	}
	totalSize := bc.Size(msg)
	if uint64(len(msg)) != totalSize {
		return fmt.Errorf("Invalid BlockContent dump, expected %dB got %dB", totalSize, len(msg))
	}

	blockN := uint8(msg[9])
	blockSize := uint16FromBytes(msg[10:12])
	fileID, err := uuid.FromBytes(msg[12:28])
	if err != nil {
		return err
	}
	content := msg[28:]

	// validate extracted values
	if len(content) > int(^uint16(0))*1024 { // bigger than MaxUint8
		return errors.New("Invalid block size")
	}

	if uint16(len(content)/1024) > blockSize {
		return errors.New("Content lenght is greater than block size")
	}

	// assign extracted values
	bc.MessageSize = totalSize
	bc.BlockN = blockN
	bc.BlockSize = blockSize
	bc.Content = content
	bc.FileID = fileID

	return nil
}

// Recv calls recvLongMessage to receive a complete BlockContent
func (bc BlockContent) Recv(s net.Stream) ([]byte, error) {
	return recvLongMessage(s, bc)
}

// Size returns the total size of the message, represented in the bytes 1 to 9
func (bc BlockContent) Size(msg []byte) uint64 {
	return uint64FromBytes(msg[1:9])
}

// Type returns the type of the Message (MTBlockContent)
func (bc BlockContent) Type() MessageType {
	return MTBlockContent
}

/* Block request */

// BlockRequest is used to ask for a block
type BlockRequest struct {
	BlockN       uint8     // block number
	FileID       uuid.UUID // ID of the file the Block belongs to
	FilePathSize uint16    // size of encoded FilePath
	FilePath     string    // relative path of the file
}

// Dump creates a byte array: {MessageType, BlockNumber, FileID} (18B)
func (br BlockRequest) Dump() []byte {
	dump := append([]byte{byte(br.Type()), byte(br.BlockN)}, br.FileID.Bytes()...)
	encodedPath := []byte(br.FilePath)
	// split string size in two bytes
	pathSize := uint16ToBytes(uint16(len(encodedPath)))
	dump = append(dump, pathSize...)
	return append(dump, encodedPath...)
}

// Load reads blockN and fileID from a byte slice created by br.Dump()
func (br *BlockRequest) Load(msg []byte) error {
	if len(msg) < 20 || MessageType(msg[0]) != MTBlockRequest {
		return errors.New("Invalid message type")
	}
	blockN := uint8(msg[1])
	fileID, err := uuid.FromBytes(msg[2:18])
	if err != nil {
		return err
	}

	filePathSize := int(uint16FromBytes(msg[18:20]))
	if len(msg) != 20+filePathSize {
		return errors.New("Incomplete message content")
	}
	filePath := string(msg[20 : 20+filePathSize])

	// both values extracted successfully
	br.BlockN = blockN
	br.FileID = fileID
	br.FilePathSize = uint16(filePathSize)
	br.FilePath = filePath
	return nil
}

// Type returns the type of the Message (MTBlockRequest)
func (br BlockRequest) Type() MessageType {
	return MTBlockRequest
}

/* Index content */

// IndexContent is used to send the fs.Index of a directory to a peer
type IndexContent struct {
	MessageSize uint64 // total size of the message
	Index       fs.Index
}

// Dump creates a byte array used to recreate (Load) the message
// The first byte contains the MessageType, the rest of them contained a
// marshalled fs.Index
func (ic IndexContent) Dump() []byte {
	index, _ := json.Marshal(ic.Index)
	dump := append([]byte{byte(MTIndexContent)}, uint64ToBytes(uint64(len(index)+9))...)
	return append(dump, index...)
}

// Load creates a fs.Index given a MTIndexContent message
func (ic *IndexContent) Load(msg []byte) error {
	if len(msg) < 9 || MessageType(msg[0]) != MTIndexContent {
		return errors.New("Invalid message type")
	}
	totalSize := ic.Size(msg)
	if uint64(len(msg)) != totalSize {
		return fmt.Errorf("Invalid BlockContent dump, expected %dB got %dB", totalSize, len(msg))
	}
	if err := json.Unmarshal(msg[9:], &ic.Index); err != nil {
		return err
	}

	ic.MessageSize = totalSize
	return nil
}

// Recv calls recvLongMessage to receive a complete IndexContent
func (ic IndexContent) Recv(s net.Stream) ([]byte, error) {
	return recvLongMessage(s, ic)
}

// Size returns the total size of the message, represented in the bytes 1 to 9
func (ic IndexContent) Size(msg []byte) uint64 {
	return uint64FromBytes(msg[1:9])
}

// Type returns the type of the Message (MTIndexContent)
func (ic IndexContent) Type() MessageType {
	return MTIndexContent
}

/* Index request */

// IndexRequest is used to ask a Peer for the index of it's assigned directory
type IndexRequest struct{}

// Dump creates a byte array used to recreate (Load) the message
func (ir IndexRequest) Dump() []byte {
	return []byte{byte(ir.Type())}
}

// Load creates a IndexRequest given the content bytes
func (ir *IndexRequest) Load(msg []byte) error {
	// The only content of the message is the message type (one byte)
	if len(msg) != 1 || MessageType(msg[0]) != MTIndexRequest {
		return errors.New("Invalid message type")
	}
	return nil
}

// Type returns the type of the Message (MTIndexRequest)
func (ir IndexRequest) Type() MessageType {
	return MTIndexRequest
}
