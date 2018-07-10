package fs

import (
	"encoding/json"
	"os"
)

// IndexedSummary stores multiple Summary structs indexed by path
// TODO: fix path redundancy
type IndexedSummary struct {
	Files map[string]*Summary `json:"files"`
}

// MakeIndexedSummary creates an IndexedSummary from a slice of summaries
func MakeIndexedSummary(summaries ...*Summary) (*IndexedSummary, error) {
	is := new(IndexedSummary)
	is.Files = make(map[string]*Summary)
	for _, s := range summaries {
		if _, found := is.Files[s.Path]; found {
			// repeated path
			return nil, os.ErrExist
		}
		is.Files[s.Path] = s
	}
	return is, nil
}

// Add adds a new set of Summary to IndexedSummary.Files
func (is *IndexedSummary) Add(summaries ...*Summary) error {
	for _, s := range summaries {
		if _, found := is.Files[s.Path]; found {
			return os.ErrExist
		}
		is.Files[s.Path] = s
	}
	return nil
}

// Delete removes a set of Summary from IndexedSummary.Files
func (is *IndexedSummary) Delete(summaries ...*Summary) error {
	for _, s := range summaries {
		if _, found := is.Files[s.Path]; !found {
			return os.ErrNotExist
		}
		delete(is.Files, s.Path)
	}
	return nil
}

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
	if f.Parent == nil {
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

// Strings creates and marshals a Summary with f *File
func (f *File) String() string {
	s := MakeSummary(f)
	b, _ := json.Marshal(s)
	return string(b)
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

// Is compares the ID, PARENT and CONTENT of a Summary
func (s *Summary) Is(s2 *Summary) bool {
	return s.ID == s2.ID && s.Parent == s2.Parent && s.Equals(s2)
}
