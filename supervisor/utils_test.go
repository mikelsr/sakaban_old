package supervisor

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"bitbucket.org/mikelsr/sakaban/fs"
)

// TestReadIndex reads an Index from a valid
// and an invalid file
func TestReadIndex(t *testing.T) {
	s, _ := MakeScanner(filepath.Join(fs.ProjectPath(), "res"))
	filename := filepath.Join(testDir, SummaryFile)
	WriteIndex(s.NewIndex, filename)

	// read existing, valid file
	_, err := ReadIndex(filename)
	if err != nil {
		t.FailNow()
	}

	// read existing, invalid file
	ioutil.WriteFile(filename, []byte{42}, 0755)
	_, err = ReadIndex(filename)
	if err == nil {
		t.FailNow()
	}

	// read non exisitng file
	_, err = ReadIndex("")
	if err == nil {
		t.FailNow()
	}
}

// TestWriteIndex writes an Index to a valid
// and an invalid diretory
func TestWriteIndex(t *testing.T) {
	s := new(Scanner)
	s.Scan(testDir)
	var err error
	s.NewIndex, err = fs.MakeIndex(s.Summaries...)
	if err != nil {
		t.Fatal(err)
		// t.Fail()
	}

	// write summary to valid path and file
	err = WriteIndex(s.NewIndex,
		filepath.Join(testDir, SummaryFile))
	if err != nil {
		t.FailNow()
	}

	// write summary to directory with no write permissions
	os.Mkdir(filepath.Join(testDir, "noperm"), 0000)
	err = WriteIndex(s.NewIndex,
		filepath.Join(testDir, "noperm", SummaryFile))
	if err == nil {
		t.FailNow()
	}
}
