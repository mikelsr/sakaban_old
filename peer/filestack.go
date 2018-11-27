package peer

import (
	"os"

	"bitbucket.org/mikelsr/sakaban/fs"
	uuid "github.com/satori/go.uuid"
)

// fileStack is used to store the files to be retrieved from another peer
type fileStack struct {
	files      []*fs.Summary
	tmpFile    *fs.File // temporary file to store blocks
	writeMutex chan bool
}

func newFileStack() *fileStack {
	return &fileStack{files: make([]*fs.Summary, 0)}
}

// iterFile loads the next file into f.tmpFile
func (f *fileStack) iterFile() {
	_, n := f.pop()
	if n != 0 {
		newFile := f.peek()
		f.tmpFile = new(fs.File)
		f.tmpFile.ID, _ = uuid.FromString(newFile.ID)
		f.tmpFile.Parent, _ = uuid.FromString(newFile.Parent)
		f.tmpFile.Path = newFile.Path
		f.tmpFile.Blocks = make([]*fs.Block, len(newFile.Blocks))
	} else {
		f.tmpFile = nil
	}
	// unlock write mutex
	f.writeMutex = make(chan bool, 1)
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
	if lenght == 0 {
		f.files = make([]*fs.Summary, 0)
	} else {
		f.files = f.files[:lenght]
	}
	return s, lenght
}

func (f *fileStack) push(s *fs.Summary) {
	f.files = append(f.files, s)
}

// write file writes the current file to permanent storage
func (f *fileStack) writeFile(perm os.FileMode) error {
	f.tmpFile.Perm = perm
	if err := f.tmpFile.Write(); err != nil {
		return err
	}
	return nil
}
