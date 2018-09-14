package comm

import (
	"encoding/json"
	"errors"

	"bitbucket.org/mikelsr/sakaban/fs"
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

/* Index content */

// IndexContent is used to send the fs.Index of a directory to a peer
type IndexContent struct {
	index fs.Index
}

// Dump creates a byte array used to recreate (Load) the message
// The first byte contains the MessageType, the rest of them contained a
// marshalled fs.Index
func (ic IndexContent) Dump() []byte {
	index, _ := json.Marshal(ic)
	return append([]byte{byte(MTIndexContent)}, index...)
}

// Load creates a fs.Index given a MTIndexContent message
func (ic *IndexContent) Load(msg []byte) error {
	if len(msg) < 2 {
		return errors.New("Invalid message type")
	}
	return json.Unmarshal(msg[1:], ic)
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
