package comm

import (
	"encoding/json"
	"errors"

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
	Type() MessageType
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
// Content} (19B + BlockSize * 1024B)
func (bc BlockContent) Dump() []byte {
	dump := append([]byte{byte(bc.Type()), byte(bc.BlockN)}, uint16ToBytes(bc.BlockSize)...)
	dump = append(dump, bc.FileID.Bytes()...)
	dump = append(dump, bc.Content...)
	totalLen := uint64(len(dump) + 8)
	cDump := []byte{dump[0]}
	cDump = append(cDump, uint64ToBytes(totalLen)...)
	cDump = append(cDump, dump[1:]...)
	return cDump
}

// Load reads blockN, blockSize, fileID, content from a byte slice created
// by br.Dump()
func (bc *BlockContent) Load(msg []byte) error {
	if len(msg) < 28 || MessageType(msg[0]) != MTBlockContent {
		return errors.New("Invalid message type")
	}
	totalSize := uint64FromBytes(msg[1:9])
	blockN := uint8(msg[9])
	blockSize := uint16FromBytes(msg[10:12])
	fileID, err := uuid.FromBytes(msg[12:28])
	if err != nil {
		return err
	}
	content := msg[28:]

	if len(content) > int(^uint16(0))*1024 { // bigger than MaxUint8
		return errors.New("Invalid block size")
	}

	if uint16(len(content)/1024) > blockSize {
		return errors.New("Content lenght is greater than block size")
	}
	bc.MessageSize = totalSize
	bc.BlockN = blockN
	bc.BlockSize = blockSize
	bc.Content = content
	bc.FileID = fileID

	return nil
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
	index fs.Index
}

// Dump creates a byte array used to recreate (Load) the message
// The first byte contains the MessageType, the rest of them contained a
// marshalled fs.Index
func (ic IndexContent) Dump() []byte {
	index, _ := json.Marshal(ic.index)
	return append([]byte{byte(MTIndexContent)}, index...)
}

// Load creates a fs.Index given a MTIndexContent message
func (ic *IndexContent) Load(msg []byte) error {
	if len(msg) < 2 || MessageType(msg[0]) != MTIndexContent {
		return errors.New("Invalid message type")
	}
	return json.Unmarshal(msg[1:], &ic.index)
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
