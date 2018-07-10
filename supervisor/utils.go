package supervisor

import (
	"encoding/json"
	"io/ioutil"

	"bitbucket.org/mikelsr/sakaban/fs"
)

// ReadIndexedSummary creates an IndexedSummary given a path
// to a file containing a valid json
func ReadIndexedSummary(filename string) (*fs.IndexedSummary, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var is fs.IndexedSummary
	err = json.Unmarshal(file, &is)
	if err != nil {
		return nil, err
	}
	return &is, nil
}

// WriteIndexedSummary writes an IndexedSummary as a JSON in
// the file specified by the path
func WriteIndexedSummary(summary *fs.IndexedSummary, filename string) error {
	content, _ := json.Marshal(summary)
	err := ioutil.WriteFile(filename, content, 0777)
	if err != nil {
		return err
	}
	return nil
}
