package supervisor

import (
	"fmt"
	"os"
	"path/filepath"

	"bitbucket.org/mikelsr/sakaban/fs"
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
	// NewSummary lightens memory usage by avoiding storing
	// files
	NewSummary []*fs.FileSummary
	// NewIndexedSummary will store the scanned summaries
	NewIndexedSummary *fs.IndexedSummary
	// NewIndexedSummary will store the read summaries
	OldIndexedSummary *fs.IndexedSummary
}

// MakeScanner creates a new scanner, tries to read
// OldIndexedSummary and create NewIndexedSummary
func MakeScanner(root string) (*Scanner, error) {
	s := new(Scanner)
	s.Root = root
	// Old IndexedSummary
	if SummaryExists(root) {
		oldIndex, err := ReadIndexedSummary(filepath.Join(root, SummaryDir, SummaryFile))
		if err != nil {
			return nil, err
		}
		s.OldIndexedSummary = oldIndex
	} else {
		s.OldIndexedSummary, _ = fs.MakeIndexedSummary()
	}

	// New IndexedSummary
	err := s.Scan(root)
	if err != nil {
		return nil, err
	}
	s.NewIndexedSummary, _ = fs.MakeIndexedSummary(s.NewSummary...)
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
		file, err := fs.MakeFile(path)
		if err != nil {
			return err
		}
		summary := fs.MakeFileSummary(file)
		s.NewSummary = append(s.NewSummary, summary)
	}
	return nil
}

// SummaryExists checks wheter the summary file exists
func SummaryExists(root string) bool {
	f, err := os.Stat(fmt.Sprintf("%s/%s/%s", root, SummaryDir, SummaryFile))
	if err != nil {
		return false
	}
	return !f.Mode().IsDir()
}
