package fs

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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
func commonRoot(s1 *Summary, s2 *Summary, parents map[string]*Summary) bool {

	parentSet := make(map[string]bool)

	// quick case
	if s1.Parent == s2.Parent {
		return true
	}

	if s1.Parent == "" || s2.Parent == "" {
		return false
	}

	// store the line of s1 in a "set"
	if _, found := parents[s1.Parent]; !found {
		return false
	}
	p := parents[s1.Parent].Parent
	for p != "" {
		parentSet[p] = true
		if _, found := parents[s1.Parent]; !found {
			return false
		}
		p = parents[p].Parent
	}

	p = s2.Parent
	for p != "" {
		if _, found := parentSet[p]; found {
			return true
		}
		if _, found := parents[p]; !found {
			return false
		}
		p = parents[p].Parent
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

// isDescendant iterates the line of one of the summaries to find out if the other
// one is a descendant
func isDescendant(descendant *Summary, ancestor *Summary, line map[string]*Summary) bool {

	if descendant.Parent == "" {
		return false
	}

	if descendant.Parent == ancestor.ID {
		return true
	}

	if s, found := line[descendant.Parent]; found {
		return isDescendant(s, ancestor, line)
	}

	return false
}

// mergeSummaryMap creates a map with keys/values from all the maps
// if ignoreCollisions is set to false, an error will be returned when an
// existing key is added to the map
func mergeSummaryMap(ignoreCollisions bool, maps ...map[string]*Summary) (map[string]*Summary, error) {

	if len(maps) < 1 {
		return nil, errors.New("No maps were passed")
	}

	m := make(map[string]*Summary)
	for k, v := range maps[0] {
		m[k] = v
	}
	for i := 1; i < len(maps); i++ {
		for k, v := range maps[i] {
			if s, conflict := m[k]; !ignoreCollisions && conflict && !v.Equals(s) {
				return nil, fmt.Errorf("Error merging key: %s", k)
			}
			m[k] = v
		}
	}

	return m, nil
}

// ProjectPath returns the directory this project is supposed to be at
func ProjectPath() string {
	return fmt.Sprintf("%s/src/bitbucket.org/mikelsr/sakaban", os.Getenv("GOPATH"))
}

// ReadIndex creates an Index given a path
// to a file containing a valid json
func ReadIndex(filename string) (*Index, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var is Index
	err = json.Unmarshal(file, &is)
	if err != nil {
		return nil, err
	}
	return &is, nil
}

// SummaryExists checks wheter the summary file exists
func SummaryExists(root string) bool {
	f, err := os.Stat(fmt.Sprintf("%s/%s/%s", root, SummaryDir, SummaryFile))
	if err != nil {
		return false
	}
	return !f.Mode().IsDir()
}

// WriteIndex writes an Index as a JSON in
// the file specified by the path
func WriteIndex(index Index, filename string) error {
	content, _ := json.Marshal(index)
	err := ioutil.WriteFile(filename, content, 0755)
	if err != nil {
		return err
	}
	return nil
}
