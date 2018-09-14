package comm

const (
	/* Message types (MT) */

	// MTFileContent is a group of blocks forming a file
	MTFileContent MessageType = iota
	// MTFileRequest is used to ask for FileContent
	MTFileRequest
	// MTIndexContent is the summary of a directory
	MTIndexContent
	// MTIndexRequest is used to ask for a IndexContent
	MTIndexRequest
)
