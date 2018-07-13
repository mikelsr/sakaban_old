package fs

import (
	"github.com/satori/go.uuid"
)

// Summary is used to marshal/unmarshal Files to/from JSON files
type Summary struct {
	ID     string   `json:"id"`
	Parent string   `json:"parent"`
	Path   string   `json:"path"`
	Blocks []uint64 `json:"blocks"`
}

// MakeSummary creates a marshable Summary from a File
func MakeSummary(f *File) *Summary {
	var parent string
	if f.Parent == uuid.Nil {
		parent = ""
	} else {
		parent = f.Parent.String()
	}
	s := Summary{ID: f.ID.String(), Parent: parent, Path: f.Path}
	s.Blocks = make([]uint64, len(f.Blocks))
	for i, b := range f.Blocks {
		s.Blocks[i] = b.Hash()
	}
	return &s
}

// Equals is used to compare both the CONTENT of a Summary
func (s *Summary) Equals(s2 *Summary) bool {
	if len(s.Blocks) != len(s2.Blocks) {
		return false
	}
	for i, b := range s.Blocks {
		if b != s2.Blocks[i] {
			return false
		}
	}
	return true
}
