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
	Files []*fs.File
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
		s.Files = append(s.Files, file)
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
