package fs

import (
	"encoding/json"
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
	Parent *uuid.UUID
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
	blocks, _ := f.Slice()
	f.Blocks = blocks
	return &f, nil
}

// MakeFileFromSummary creates a File given a Summary
func MakeFileFromSummary(s *Summary) (*File, error) {
	f, err := MakeFile(s.Path)
	if err != nil {
		return nil, err
	}
	f.ID, err = uuid.FromString(s.ID)
	if err != nil {
		return nil, fmt.Errorf("Invalid ID: %s", s.ID)
	}
	if s.Parent != "" {
		parent, err := uuid.FromString(s.Parent)
		f.Parent = &parent
		if err != nil {
			return nil, fmt.Errorf("Invalid parent ID: %s", s.ID)
		}
	}
	return f, nil
}

// DeepEquals compares two files by individually comparing each byte of
// it's blocks
func (f *File) DeepEquals(f2 *File) bool {
	if f == f2 {
		return true
	}
	if len(f.Blocks) != len(f2.Blocks) {
		return false
	}
	for i := 0; i < len(f.Blocks); i++ {
		if !f.Blocks[i].DeepEquals(f2.Blocks[i]) {
			return false
		}
	}
	return true
}

// Equals makes a shallow comparison between the hashes of the blocks of a file
// It is used to compare the CONTENT of a file
func (f *File) Equals(f2 *File) bool {
	if f == f2 {
		return true
	}
	if len(f.Blocks) != len(f2.Blocks) {
		return false
	}
	for i := 0; i < len(f.Blocks); i++ {
		if !f.Blocks[i].Equals(f2.Blocks[i]) {
			return false
		}
	}
	return true
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

// Strings creates and marshals a Summary with f *File
func (f *File) String() string {
	s := MakeSummary(f)
	b, _ := json.Marshal(s)
	return string(b)
}
