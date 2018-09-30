package fs

const (
	// BlockSize defines the size of each block in Bytes
	BlockSize int64 = 1024 * 1024 // 1024 kB

	// SummaryDir is the relative directory the summary is stored at
	SummaryDir = ".sakaban"
	// SummaryFile is the relative name of the file containing the summary
	SummaryFile = "sakaban.json"
)
