package comm

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"

	"bitbucket.org/mikelsr/sakaban/fs"
	"github.com/satori/go.uuid"
)

// MessageType is used to identify the message type using one byte
type MessageType byte

// Message ensures anything sent/received by peers can be dumped to or loaded
// from a byte array
type Message interface {
	Dump() []byte
	Load([]byte) error
	Recv(s *bufio.Reader) ([]byte, error)
	Size([]byte) uint64
	Type() MessageType
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
	totalLen := uint64(len(dump) + sizeOfMessage)
	// create new massage merging the base message and its size
	cDump := []byte{dump[0]}
	cDump = append(cDump, uint64ToBytes(totalLen)...)
	cDump = append(cDump, dump[1:]...)
	return cDump
}

// Load reads blockN, blockSize, fileID, content from a byte slice created
// by br.Dump()
func (bc *BlockContent) Load(msg []byte) error {
	index := 0
	headerSize := sizeOfMessageType + sizeOfMessage + sizeOfFileID + sizeOfBlockN + sizeOfBlockSize

	// parse contents of the message, extract values
	if len(msg) < headerSize || MessageType(msg[0]) != MTBlockContent {
		return errors.New("Invalid message type")
	}
	index += sizeOfMessageType

	totalSize := bc.Size(msg)
	if uint64(len(msg)) != totalSize {
		return fmt.Errorf("Invalid BlockContent dump, expected %dB got %dB", totalSize, len(msg))
	}
	index += sizeOfMessage

	// block number
	blockN := uint8(msg[index])
	index += sizeOfBlockN

	// block size
	blockSize := uint16FromBytes(msg[index : index+sizeOfBlockSize])
	index += sizeOfBlockSize

	// file id
	fileID, err := uuid.FromBytes(msg[index : index+sizeOfFileID])
	if err != nil {
		return err
	}
	index += sizeOfFileID

	// content
	content := msg[index:]

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

// Recv calls RecvMessage to receive a complete BlockContent
func (bc *BlockContent) Recv(s *bufio.Reader) ([]byte, error) {
	return RecvMessage(s, bc)
}

// Size returns the total size of the message
func (bc BlockContent) Size(msg []byte) uint64 {
	return uint64FromBytes(msg[sizeOfMessageType : sizeOfMessage+sizeOfMessage])
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
	index := 0
	headerSize := sizeOfMessageType + sizeOfBlockN + sizeOfFileID + sizeOfFilePathSize

	if len(msg) < headerSize || MessageType(msg[0]) != MTBlockRequest {
		return errors.New("Invalid message type")
	}
	index += sizeOfMessageType

	// block number
	blockN := uint8(msg[index])
	index += sizeOfBlockN

	// file id
	fileID, err := uuid.FromBytes(msg[index : index+sizeOfFileID])
	if err != nil {
		return err
	}
	index += sizeOfFileID

	// filepath size
	filePathSize := int(uint16FromBytes(msg[index : index+sizeOfFilePathSize]))
	index += sizeOfFilePathSize
	if len(msg) != index+filePathSize {
		return errors.New("Incomplete message content")
	}

	// filepath
	filePath := string(msg[index : index+filePathSize])

	// both values extracted successfully
	br.BlockN = blockN
	br.FileID = fileID
	br.FilePathSize = uint16(filePathSize)
	br.FilePath = filePath
	return nil
}

// Recv calls RecvMessage to receive a complete BlockRequest
func (br *BlockRequest) Recv(s *bufio.Reader) ([]byte, error) {
	return RecvMessage(s, br)
}

// Size returns the total size of the message
// MessageType + BlockN + UUID + FilePathSize + filePath
func (br BlockRequest) Size(msg []byte) uint64 {
	// return 1 + 1 + 2 + uint64(uint16FromBytes(msg[18:34]))
	s := sizeOfMessageType + sizeOfBlockN + sizeOfFileID
	return uint64(s) + uint64(sizeOfFilePathSize) + uint64(uint16FromBytes(msg[s:s+sizeOfFilePathSize]))
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
	dump := append([]byte{byte(MTIndexContent)}, uint64ToBytes(uint64(len(index)+sizeOfMessageType+sizeOfMessage))...)
	return append(dump, index...)
}

// Load creates a fs.Index given a MTIndexContent message
func (ic *IndexContent) Load(msg []byte) error {
	index := 0
	headerSize := sizeOfMessageType + sizeOfMessage

	if len(msg) < headerSize || MessageType(msg[0]) != MTIndexContent {
		return errors.New("Invalid message type")
	}
	index += sizeOfMessageType

	// message size
	totalSize := ic.Size(msg)
	index += sizeOfMessage

	// content of index
	if uint64(len(msg)) != totalSize {
		return fmt.Errorf("Invalid BlockContent dump, expected %dB got %dB", totalSize, len(msg))
	}
	if err := json.Unmarshal(msg[index:], &ic.Index); err != nil {
		return err
	}

	ic.MessageSize = totalSize
	return nil
}

// Recv calls RecvMessage to receive a complete IndexContent
func (ic *IndexContent) Recv(s *bufio.Reader) ([]byte, error) {
	return RecvMessage(s, ic)
}

// Size returns the total size of the message, represented in the bytes 1 to 9
func (ic IndexContent) Size(msg []byte) uint64 {
	return uint64FromBytes(msg[sizeOfMessageType : sizeOfMessageType+sizeOfMessage])
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
	if len(msg) != sizeOfMessageType || MessageType(msg[0]) != MTIndexRequest {
		return errors.New("Invalid message type")
	}
	return nil
}

// Recv calls RecvMessage to receive a complete IndexRequest
// REDUNDANT, IndexRequest has no content
func (ir *IndexRequest) Recv(s *bufio.Reader) ([]byte, error) {
	return RecvMessage(s, ir)
}

// Size returns the total size of the message
func (ir IndexRequest) Size(msg []byte) uint64 {
	return 1
}

// Type returns the type of the Message (MTIndexRequest)
func (ir IndexRequest) Type() MessageType {
	return MTIndexRequest
}
