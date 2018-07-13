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
	// Summaries lightens memory usage by avoiding storing
	// files
	Summaries []*fs.Summary
	// NewIndex will store the scanned summaries
	NewIndex *fs.Index
	// NewIndex will store the read summaries
	OldIndex *fs.Index
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
		s.OldIndex, _ = fs.MakeIndex()
	}

	// New Index
	err := s.Scan(root)
	if err != nil {
		return nil, err
	}
	s.NewIndex, _ = fs.MakeIndex(s.Summaries...)
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
		summary := fs.MakeSummary(file)
		s.Summaries = append(s.Summaries, summary)
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
