package peer

import (
	"sync"

	"bitbucket.org/mikelsr/sakaban/fs"
)

// fileStack is used to store the files to be retrieved from another peer
type fileStack struct {
	files      []*fs.Summary
	writeMutex sync.Mutex
	tmpFile    *fs.File // temporary file to store blocks
}

func newFileStack() *fileStack {
	return &fileStack{files: make([]*fs.Summary, 0)}
}

func (f *fileStack) peek() *fs.Summary {
	lenght := len(f.files)
	if lenght == 0 {
		return nil
	}
	return f.files[lenght-1]
}

func (f *fileStack) pop() (*fs.Summary, int) {
	lenght := len(f.files) - 1
	if lenght == -1 {
		return nil, -1
	}
	s := f.files[lenght]
	f.files = f.files[:lenght]
	return s, lenght
}

func (f *fileStack) push(s *fs.Summary) {
	f.files = append(f.files, s)
}
