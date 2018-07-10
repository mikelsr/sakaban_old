package fs

import (
	"os"
	"reflect"
)

// IndexedSummary stores multiple Summary structs indexed by path
// TODO: fix path redundancy
type IndexedSummary struct {
	Files   map[string]*Summary `json:"files"`
	Parents []*Summary          `json:"parents"`
}

// MakeIndexedSummary creates an IndexedSummary from a slice of summaries
func MakeIndexedSummary(summaries ...*Summary) (*IndexedSummary, error) {
	is := new(IndexedSummary)
	is.Files = make(map[string]*Summary)
	is.Parents = make([]*Summary, 0)
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

// Update compares two IndexSummaries and returns the resulting IndexedSummary
func (is *IndexedSummary) Update(newIS *IndexedSummary) *IndexedSummary {
	u, _ := MakeIndexedSummary()
	// look for old files
	for path, s := range is.Files {
		ns, found := newIS.Files[path]
		// file may have been updated
		if found {
			if s.Equals(ns) { // File is equal
				u.Add(ns)
			} else { // File has been updated
				// TODO: Allow record of child and parents in the same path
				u.Add(&Summary{ID: ns.ID, Parent: s.ID, Path: path, Blocks: ns.Blocks})
				u.Parents = append(u.Parents, s)
			}
		} else {
			// comparing contents is slow
			for _, ns := range newIS.Files {
				// file has been moved
				if reflect.DeepEqual(s.Blocks, ns.Blocks) {
					u.Add(&Summary{ID: ns.ID, Parent: s.ID, Path: ns.Path, Blocks: ns.Blocks})
					u.Parents = append(u.Parents, s)
					break
				}
			}
		}
	}
	// add missing (newly created) files
	for path, s := range newIS.Files {
		if _, found := u.Files[path]; !found {
			u.Add(s)
		}
	}
	return u
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
