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

// FileSummary is used to marshal/unmarshal Files to/from JSON files
type FileSummary struct {
	ID     string   `json:"id"`
	Parent string   `json:"parent"`
	Path   string   `json:"path"`
	Blocks []uint64 `json:"blocks"`
}

// IndexedSummary stores multiple FileSummary structs indexed by path
// TODO: remove path redundancy
type IndexedSummary struct {
	Files map[string]*FileSummary `json:"files"`
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

// MakeFileFromSummary creates a File given a FileSummary
func MakeFileFromSummary(fSum *FileSummary) (*File, error) {
	f, err := MakeFile(fSum.Path)
	if err != nil {
		return nil, err
	}
	f.ID, err = uuid.FromString(fSum.ID)
	if err != nil {
		return nil, fmt.Errorf("Invalid ID: %s", fSum.ID)
	}
	if fSum.Parent != "" {
		parent, err := uuid.FromString(fSum.Parent)
		f.Parent = &parent
		if err != nil {
			return nil, fmt.Errorf("Invalid parent ID: %s", fSum.ID)
		}
	}
	return f, nil
}

// MakeFileSummary creates a marshable FileSummary from a File
func MakeFileSummary(f *File) *FileSummary {
	var parent string
	if f.Parent == nil {
		parent = ""
	} else {
		parent = f.Parent.String()
	}
	fSum := FileSummary{ID: f.ID.String(), Parent: parent, Path: f.Path}
	fSum.Blocks = make([]uint64, len(f.Blocks))
	for i, b := range f.Blocks {
		fSum.Blocks[i] = b.Hash()
	}
	return &fSum
}

// MakeIndexedSummary creates an IndexedSummary from a slice of summaries
func MakeIndexedSummary(summaries ...*FileSummary) (*IndexedSummary, error) {
	is := new(IndexedSummary)
	is.Files = make(map[string]*FileSummary)
	for _, s := range summaries {
		if _, found := is.Files[s.Path]; found {
			// repeated path
			return nil, os.ErrExist
		}
		is.Files[s.Path] = s
	}
	return is, nil
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

// Strings creates and marshals a FileSummary with f *File
func (f *File) String() string {
	fSum := MakeFileSummary(f)
	s, _ := json.Marshal(fSum)
	return string(s)
}

// Equals is used to compare both the CONTENT of a FileSummary
func (fSum *FileSummary) Equals(fSum2 *FileSummary) bool {
	if len(fSum.Blocks) != len(fSum2.Blocks) {
		return false
	}
	for i, b := range fSum.Blocks {
		if b != fSum2.Blocks[i] {
			return false
		}
	}
	return true
}

// Is compares the ID, PARENT and CONTENT of a FileSummary
func (fSum *FileSummary) Is(fSum2 *FileSummary) bool {
	return fSum.ID == fSum2.ID && fSum.Parent == fSum2.Parent && fSum.Equals(fSum2)
}

// Add adds a new set of FileSummary to IndexedSummary.Files
func (is *IndexedSummary) Add(summaries ...*FileSummary) error {
	for _, s := range summaries {
		if _, found := is.Files[s.Path]; found {
			return os.ErrExist
		}
		is.Files[s.Path] = s
	}
	return nil
}

// Delete removes a set of FileSummary from IndexedSummary.Files
func (is *IndexedSummary) Delete(summaries ...*FileSummary) error {
	for _, s := range summaries {
		if _, found := is.Files[s.Path]; !found {
			return os.ErrNotExist
		}
		delete(is.Files, s.Path)
	}
	return nil
}
