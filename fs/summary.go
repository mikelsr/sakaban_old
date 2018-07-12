package fs

import (
	"fmt"
	"os"
	"reflect"

	"github.com/satori/go.uuid"
)

// IndexedSummary stores multiple Summary structs indexed by path
// TODO: fix path redundancy
type IndexedSummary struct {
	Files   map[string]*Summary `json:"files"`
	Parents map[string]*Summary `json:"parents"`
	// TODO: is anything other than the ID of Deletions used?
	Deletions map[string]*Summary `json:"deletions"`
}

// MakeIndexedSummary creates an IndexedSummary from a slice of summaries
func MakeIndexedSummary(summaries ...*Summary) (*IndexedSummary, error) {
	is := new(IndexedSummary)
	is.Files = make(map[string]*Summary)
	is.Parents = make(map[string]*Summary)
	is.Deletions = make(map[string]*Summary)
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

// AddParent adds a new set of Summary to IndexedSummary.Parents
func (is *IndexedSummary) AddParent(summaries ...*Summary) error {
	for _, s := range summaries {
		if _, found := is.Parents[s.ID]; found {
			return os.ErrExist
		}
		is.Parents[s.ID] = s
	}
	return nil
}

// Contains compares the hashes of the blocks of a file with the summaries in
// IndexedSummary.Files and returns (path to the file, found/not found)
func (is *IndexedSummary) Contains(s *Summary) (string, bool) {
	for path, s2 := range is.Files {
		if s.Equals(s2) {
			return path, true
		}
	}
	return "", false
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

// DeleteParent removes a set of Summary from IndexedSummary.Parents
// Not used
func (is *IndexedSummary) DeleteParent(summaries ...*Summary) error {
	for _, s := range summaries {
		if _, found := is.Parents[s.ID]; !found {
			return os.ErrNotExist
		}
		delete(is.Parents, s.ID)
	}
	return nil
}

// Update compares two IndexSummaries from the same local directory
// and returns the resulting IndexedSummary
func (is *IndexedSummary) Update(newIS *IndexedSummary) *IndexedSummary {
	u, _ := MakeIndexedSummary()
	// look for old files
Lookup:
	for path, s := range is.Files {
		ns, found := newIS.Files[path]
		// file may have been updated
		if found {
			if s.Equals(ns) { // File is equal
				u.Add(ns)
			} else { // File has been updated
				// TODO: Allow record of child and parents in the same path
				u.Add(&Summary{ID: ns.ID, Parent: s.ID, Path: path, Blocks: ns.Blocks})
				u.AddParent(s)
			}
		} else {
			// comparing contents is slow
			for _, ns := range newIS.Files {
				// file has been moved
				if reflect.DeepEqual(s.Blocks, ns.Blocks) {
					u.Add(&Summary{ID: ns.ID, Parent: s.ID, Path: ns.Path, Blocks: ns.Blocks})
					u.AddParent(s)
					continue Lookup
				}
			}
			// file deleted
			u.Deletions[s.ID] = s
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

// Merge compares a summary of a local and a remote directory
// This function should return the same summary switching s1 and s2
func Merge(is1 *IndexedSummary, is2 *IndexedSummary) (*IndexedSummary, error) {
	m, _ := MakeIndexedSummary()

	// merge parents
	parents, err := mergeSummaryMap(false, is1.Parents, is2.Parents)
	if err != nil {
		return nil, err
	}
	m.Parents = parents

	// merge deletions
	m.Deletions, _ = mergeSummaryMap(true, is1.Deletions, is2.Deletions)

	// filter misdeletions out
	for id := range m.Parents {
		if _, found := m.Deletions[id]; found {
			delete(m.Deletions, id)
		}
	}

	for path, s := range is1.Files {
		if ns, found := is2.Files[path]; found {
			// same file
			if s.ID == ns.ID {
				m.Add(s)
				continue
			}

			if isDescendant(s, ns, is1.Parents) {
				m.Add(ns)
				continue
			}
			if isDescendant(ns, s, is2.Parents) {
				m.Add(s)
				continue
			}
			// branches of the same file
			if commonRoot(s, ns, m.Parents) {
				s1 := *s
				s2 := *ns
				s1.Path = fmt.Sprintf("%s_%s", s.Path, s.ID)
				s2.Path = fmt.Sprintf("%s_%s", ns.Path, ns.ID)
				m.Add(&s1, &s2)
				continue
			}
		} else {
			// file is now a parent
			if _, parent := m.Parents[s.ID]; parent {
				continue
			}
			// file is now deleted
			if _, deleted := m.Deletions[s.ID]; deleted {
				continue
			}
			m.Add(s)
		}
	}

	for path, s := range is2.Files {
		// file has been merged/added already
		if _, found := is1.Files[path]; found {
			continue
		}
		// file is now a parent
		if _, parent := m.Parents[s.ID]; parent {
			continue
		}
		// file is now deleted
		if _, deleted := m.Deletions[s.ID]; deleted {
			continue
		}
		m.Add(s)
	}

	return m, nil
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

// Is compares the ID, PARENT and CONTENT of a Summary
func (s *Summary) Is(s2 *Summary) bool {
	return s.ID == s2.ID && s.Parent == s2.Parent && s.Equals(s2)
}
