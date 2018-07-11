package fs

import (
	"fmt"
	"os"
)

// CalcBlockN calculates the number of blocks to
// be extracted from a file
func CalcBlockN(f os.FileInfo) int {
	n := f.Size() / BlockSize
	if f.Size()-BlockSize*n == 0 {
		return int(n)
	}
	return int(n + 1)
}

// commonRoot iteratively checks if the files have a common ancestor
// Time complexity (quick case aside): best -> O(n), worst -> O(n+m)
// Space complexity: best == worst -> O(n)
func commonRoot(s1 *Summary, s2 *Summary, p1 map[string]*Summary, p2 map[string]*Summary) bool {
	parents := make(map[string]bool)

	// quick case
	if s1.Parent == s2.Parent {
		return true
	}

	// store the line of s1 in a "set"
	p := p1[s1.Parent].Parent
	for p != "" {
		parents[p] = true
		p = p1[p].Parent
	}

	p = s2.Parent
	for p != "" {
		if _, found := parents[p]; found {
			return true
		}
		p = p2[p].Parent
	}

	return false
}

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

// ProjectPath returns the directory this project is supposed to be at
func ProjectPath() string {
	return fmt.Sprintf("%s/src/bitbucket.org/mikelsr/sakaban", os.Getenv("GOPATH"))
}
