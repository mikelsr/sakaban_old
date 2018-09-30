package fs

import (
	"os"
	"path/filepath"
)

const (
	// SummaryDir is the relative directory the summary is stored at
	SummaryDir = ".sakaban"
	// SummaryFile is the relative name of the file containing the summary
	SummaryFile = "sakaban.json"
)

// Scanner will be used to scan a directory and generate File structs
type Scanner struct {
	Root string
	// Summaries lightens memory usage by avoiding storing
	// files
	Summaries []*Summary
	// NewIndex will store the scanned summaries
	NewIndex *Index
	// NewIndex will store the read summaries
	OldIndex *Index
}

// MakeScanner creates a new scanner, tries to read
// OldIndex and create NewIndex
func MakeScanner(root string) (*Scanner, error) {
	s := new(Scanner)
	s.Root = root
	// Old Index
	if SummaryExists(root) {
		oldIndex, err := ReadIndex(filepath.Join(root, SummaryDir, SummaryFile))
		if err != nil {
			return nil, err
		}
		s.OldIndex = oldIndex
	} else {
		s.OldIndex, _ = MakeIndex()
	}

	// New Index
	err := s.Scan(root)
	if err != nil {
		return nil, err
	}
	s.NewIndex, _ = MakeIndex(s.Summaries...)
	return s, nil
}

// Scan runs Scanner.VisitDir in path (root folder) and each subdirectory
func (s *Scanner) Scan(path string) error {
	if err := filepath.Walk(path, s.Visit); err != nil {
		return err
	}
	return nil
}

// Visit creates a File when visiting a file and appends it to Summary.Files
func (s *Scanner) Visit(path string, f os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if f.Mode().IsRegular() {
		file, err := MakeFile(path)
		if err != nil {
			return err
		}
		summary := MakeSummary(file)
		s.Summaries = append(s.Summaries, summary)
	}
	return nil
}
