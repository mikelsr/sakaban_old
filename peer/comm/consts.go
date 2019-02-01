package comm

import (
	"unsafe"

	uuid "github.com/satori/go.uuid"
)

const (
	/* Message types (MT) */

	// MTBlockContent is a group of blocks forming a file
	MTBlockContent MessageType = iota
	// MTBlockRequest is used to ask for BlockContent
	MTBlockRequest
	// MTIndexContent is the summary of a directory
	MTIndexContent
	// MTIndexRequest is used to ask for a IndexContent
	MTIndexRequest
)

const minMessageType MessageType = MTBlockContent
const maxMessageType MessageType = MTIndexRequest

const (
	/* other constats */
	bufferSize = 1024 * 1024 * 2 // recv buffer size
	/* sizes of fields */
	sizeOfBlockN       = int(unsafe.Sizeof(uint8(0)))
	sizeOfBlockSize    = int(unsafe.Sizeof(uint16(0)))
	sizeOfFileID       = uuid.Size
	sizeOfFilePathSize = int(unsafe.Sizeof(uint16(0)))
	sizeOfMessage      = int(unsafe.Sizeof(uint64(0)))
	sizeOfMessageType  = 1
)
