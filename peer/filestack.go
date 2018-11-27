package peer

import (
	"sync"

	"bitbucket.org/mikelsr/sakaban/fs"
	uuid "github.com/satori/go.uuid"
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

// iterFile loads the next file into f.tmpFile
func (f *fileStack) iterFile() {
	f.writeMutex.Lock()
	defer f.writeMutex.Unlock()
	_, n := f.pop()
	if n != -1 {
		newFile := f.peek()
		f.tmpFile = new(fs.File)
		f.tmpFile.ID, _ = uuid.FromString(newFile.ID)
		f.tmpFile.Parent, _ = uuid.FromString(newFile.Parent)
		f.tmpFile.Path = newFile.Path
		f.tmpFile.Blocks = make([]*fs.Block, len(newFile.Blocks))
	} else {
		f.tmpFile = nil
	}
}

// write file writes the current file to permanent storage
func (f *fileStack) writeFile() error {
	f.writeMutex.Lock()
	defer f.writeMutex.Unlock()
	if err := f.tmpFile.Write(); err != nil {
		return err
	}
	return nil
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
