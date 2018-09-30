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
	blockN    uint8     // block number
	blockSize uint8     // block size in kB
	content   []byte    // content of the block
	fileID    uuid.UUID // ID of the file the Block belongs to
}

// Dump creates a byte array: {MessageType, BlockNumber, BlockSize, FileID,
// Content} (19B + BlockSize * 1024B)
func (bc BlockContent) Dump() []byte {
	return append(append([]byte{byte(bc.Type()), byte(bc.blockN),
		byte(bc.blockSize)}, bc.fileID.Bytes()...), bc.content...)
}

// Load reads blockN, blockSize, fileID, content from a byte slice created
// by br.Dump()
func (bc *BlockContent) Load(msg []byte) error {
	if len(msg) < 18 || MessageType(msg[0]) != MTBlockContent {
		return errors.New("Invalid message type")
	}
	blockN := uint8(msg[1])
	blockSize := uint8(msg[2])
	fileID, err := uuid.FromBytes(msg[3:19])
	if err != nil {
		return err
	}
	content := msg[19:]

	if len(content) > int(^uint8(0))*1024 { // bigger than MaxUint8
		return errors.New("Invalid block size")
	}
	if uint8(len(content)/1024) != blockSize {
		return errors.New("Block size and content length do not match")
	}
	bc.blockN = blockN
	bc.blockSize = blockSize
	bc.content = content
	bc.fileID = fileID

	return nil
}

// Type returns the type of the Message (MTBlockContent)
func (bc BlockContent) Type() MessageType {
	return MTBlockContent
}

/* Block request */

// BlockRequest is used to ask for a block
type BlockRequest struct {
	BlockN uint8     // block number
	FileID uuid.UUID // ID of the file the Block belongs to
}

// Dump creates a byte array: {MessageType, BlockNumber, FileID} (18B)
func (br BlockRequest) Dump() []byte {
	return append([]byte{byte(br.Type()), byte(br.BlockN)}, br.FileID.Bytes()...)
}

// Load reads blockN and fileID from a byte slice created by br.Dump()
func (br *BlockRequest) Load(msg []byte) error {
	if len(msg) != 18 || MessageType(msg[0]) != MTBlockRequest {
		return errors.New("Invalid message type")
	}
	blockN := uint8(msg[1])
	fileID, err := uuid.FromBytes(msg[2:])
	if err != nil {
		return err
	}

	// both values extracted successfully
	br.BlockN = blockN
	br.FileID = fileID
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
