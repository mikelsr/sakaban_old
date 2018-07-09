package fs

import (
	"fmt"
	"os"
)

// IsFile checks if a path exists and contains a file, not a directory
func IsFile(path string) bool {
	file, err := os.Stat(path)
	// Check that file exists
	if os.IsNotExist(err) {
		return false
	}
	// Check that file is not a directory
	return file.Mode().IsRegular()
}

// CalcBlockN calculates the number of blocks to
// be extracted from a file
func CalcBlockN(f os.FileInfo) int {
	n := f.Size() / BlockSize
	if f.Size()-BlockSize*n == 0 {
		return int(n)
	}
	return int(n + 1)
}

// ProjectPath returns the directory this project is supposed to be at
func ProjectPath() string {
	return fmt.Sprintf("%s/src/bitbucket.org/mikelsr/sakaban", os.Getenv("GOPATH"))
}
