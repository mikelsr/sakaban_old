package comm

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
