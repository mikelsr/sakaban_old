package fs

import (
	"fmt"
	"io"
	"os"

	"github.com/satori/go.uuid"
)

// BlockSize defines the size of each block in Bytes
const BlockSize int64 = 1024 * 1024 // 1024 kB

// File represents a file
//	ID: unique id of the file
//	Path: path to the file
//	Blocks: Blocks that form the file
type File struct {
	ID     uuid.UUID
	Path   string
	Blocks []*Block
}

// MakeFile is the default constructor for File
// it generates an ID and ensures that the path is valid
func MakeFile(path string) (*File, error) {
	id, _ := uuid.NewV4()
	if !IsFile(path) {
		return nil, fmt.Errorf("Not a valid path to a file: '%s'", path)
	}
	f := File{ID: id, Path: path}
	return &f, nil
}

// Slice divides a file into Blocks
func (f *File) Slice() ([]*Block, error) {
	file, err := os.Open(f.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	blocks := make([]*Block, 0)

	for {
		bytes := make([]byte, BlockSize)
		n, err := file.Read(bytes)
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		bytes = bytes[:n]
		block := Block{Content: bytes}
		blocks = append(blocks, &block)
	}
	return blocks, nil
}
