package fs

import (
	"fmt"
	"os"
	"reflect"
)

// Index stores multiple Summary structs:
//	Current files indexed by path
//	Parent files indexed by ID
//	Deleted files indexed by ID
type Index struct {
	Files   map[string]*Summary `json:"files"`
	Parents map[string]*Summary `json:"parents"`
	// TODO: is anything other than the ID of Deletions used?
	Deletions map[string]*Summary `json:"deletions"`
}

// MakeIndex creates an Index from a slice of summaries
func MakeIndex(summaries ...*Summary) (*Index, error) {
	i := new(Index)
	i.Files = make(map[string]*Summary)
	i.Parents = make(map[string]*Summary)
	i.Deletions = make(map[string]*Summary)
	for _, s := range summaries {
		if _, found := i.Files[s.Path]; found {
			// repeated path
			return nil, os.ErrExist
		}
		i.Files[s.Path] = s
	}
	return i, nil
}

// Add adds a new set of Summary to Index.Files
func (i *Index) Add(summaries ...*Summary) error {
	for _, s := range summaries {
		if _, found := i.Files[s.Path]; found {
			return os.ErrExist
		}
		i.Files[s.Path] = s
	}
	return nil
}

// AddParent adds a new set of Summary to Index.Parents
func (i *Index) AddParent(summaries ...*Summary) error {
	for _, s := range summaries {
		if _, found := i.Parents[s.ID]; found {
			return os.ErrExist
		}
		i.Parents[s.ID] = s
	}
	return nil
}

// Contains compares the hashes of the blocks of a file with the summaries in
// Index.Files and returns (path to the file, found/not found)
func (i *Index) Contains(s *Summary) (string, bool) {
	for path, s2 := range i.Files {
		if s.Equals(s2) {
			return path, true
		}
	}
	return "", false
}

// Delete removes a set of Summary from Index.Files
func (i *Index) Delete(summaries ...*Summary) error {
	for _, s := range summaries {
		if _, found := i.Files[s.Path]; !found {
			return os.ErrNotExist
		}
		delete(i.Files, s.Path)
	}
	return nil
}

// DeleteParent removes a set of Summary from Index.Parents
// Not used
func (i *Index) DeleteParent(summaries ...*Summary) error {
	for _, s := range summaries {
		if _, found := i.Parents[s.ID]; !found {
			return os.ErrNotExist
		}
		delete(i.Parents, s.ID)
	}
	return nil
}

// Update compares two IndexSummaries from the same local directory
// and returns the resulting Index
func (i *Index) Update(newIndex *Index) *Index {
	u, _ := MakeIndex()
	// look for old files
Lookup:
	for path, s := range i.Files {
		ns, found := newIndex.Files[path]
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
			for _, ns := range newIndex.Files {
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
	for path, s := range newIndex.Files {
		if _, found := u.Files[path]; !found {
			u.Add(s)
		}
	}
	return u
}

// Merge compares a summary of a local and a remote directory
// This function should return the same summary switching s1 and s2
func Merge(is1 *Index, is2 *Index) (*Index, error) {
	m, _ := MakeIndex()

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