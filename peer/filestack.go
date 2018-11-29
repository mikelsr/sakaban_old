package peer

import (
	"os"

	"bitbucket.org/mikelsr/sakaban/fs"
	uuid "github.com/satori/go.uuid"
)

// fileStack is used to store the files to be retrieved from another peer
type fileStack struct {
	files           []*fs.Summary
	providers       []*Contact
	tmpFile         *fs.File // temporary file to store blocks
	tmpFileProvider *Contact
	writeMutex      chan bool
}

func newFileStack() *fileStack {
	return &fileStack{files: make([]*fs.Summary, 0), providers: make([]*Contact, 0)}
}

// iterFile loads the next file into f.tmpFile
func (f *fileStack) iterFile() error {
	_, _, n := f.pop()
	if n != 0 {
		newFile, newProvider := f.peek()
		f.tmpFile = new(fs.File)
		f.tmpFile.ID, _ = uuid.FromString(newFile.ID)
		f.tmpFile.Parent, _ = uuid.FromString(newFile.Parent)
		f.tmpFile.Path = newFile.Path
		f.tmpFile.Blocks = make([]*fs.Block, len(newFile.Blocks))
		// load unchanged blocks from local file
		if _, err := os.Stat(newFile.Path); !os.IsNotExist(err) {
			localFile, err := fs.MakeFile(newFile.Path)
			if err != nil {
				return err
			}
			l := uint64(len(localFile.Blocks))
			for i, n := range newFile.Blocks {
				if n >= l {
					break
				}
				if n == 0 {
					f.tmpFile.Blocks[i] = localFile.Blocks[i]
				}
			}
		}
		f.tmpFileProvider = newProvider
	} else {
		f.tmpFile = nil
		f.tmpFileProvider = nil
	}
	// unlock write mutex
	f.writeMutex = make(chan bool, 1)
	return nil
}

func (f *fileStack) peek() (*fs.Summary, *Contact) {
	lenght := len(f.files)
	if lenght == 0 {
		return nil, nil
	}
	return f.files[lenght-1], f.providers[lenght-1]
}

func (f *fileStack) pop() (*fs.Summary, *Contact, int) {
	lenght := len(f.files) - 1
	if lenght == -1 {
		return nil, nil, -1
	}
	s := f.files[lenght]
	p := f.providers[lenght]
	if lenght == 0 {
		f.files = make([]*fs.Summary, 0)
		f.providers = make([]*Contact, 0)
	} else {
		f.files = f.files[:lenght]
		f.providers = f.providers[:lenght]
	}
	return s, p, lenght
}

func (f *fileStack) push(s *fs.Summary, p *Contact) {
	f.files = append(f.files, s)
	f.providers = append(f.providers, p)
}

// write file writes the current file to permanent storage
func (f *fileStack) writeFile(perm os.FileMode) error {
	f.tmpFile.Perm = perm
	if err := f.tmpFile.Write(); err != nil {
		return err
	}
	return nil
}
