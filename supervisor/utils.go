package supervisor

import (
	"encoding/json"
	"io/ioutil"

	"bitbucket.org/mikelsr/sakaban/fs"
)

// ReadIndex creates an Index given a path
// to a file containing a valid json
func ReadIndex(filename string) (*fs.Index, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var is fs.Index
	err = json.Unmarshal(file, &is)
	if err != nil {
		return nil, err
	}
	return &is, nil
}

// WriteIndex writes an Index as a JSON in
// the file specified by the path
func WriteIndex(index *fs.Index, filename string) error {
	content, _ := json.Marshal(index)
	err := ioutil.WriteFile(filename, content, 0755)
	if err != nil {
		return err
	}
	return nil
}
