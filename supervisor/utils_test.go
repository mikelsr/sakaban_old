package supervisor

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"bitbucket.org/mikelsr/sakaban/fs"
)

// TestReadIndexedSummary reads an IndexedSummary from a valid
// and an invalid file
func TestReadIndexedSummary(t *testing.T) {
	s, _ := MakeScanner(filepath.Join(fs.ProjectPath(), "res"))
	filename := filepath.Join(testDir, SummaryFile)
	WriteIndexedSummary(s.NewIndexedSummary, filename)

	// read existing, valid file
	_, err := ReadIndexedSummary(filename)
	if err != nil {
		t.FailNow()
	}

	// read existing, invalid file
	ioutil.WriteFile(filename, []byte{42}, 0777)
	_, err = ReadIndexedSummary(filename)
	if err == nil {
		t.FailNow()
	}

	// read non exisitng file
	_, err = ReadIndexedSummary("")
	if err == nil {
		t.FailNow()
	}
}

// TestWriteIndexedSummary writes an IndexedSummary to a valid
// and an invalid diretory
func TestWriteIndexedSummary(t *testing.T) {
	s := new(Scanner)
	s.Scan(testDir)
	var err error
	s.NewIndexedSummary, err = fs.MakeIndexedSummary(s.NewSummary...)
	if err != nil {
		t.Fatal(err)
		// t.Fail()
	}

	// write summary to valid path and file
	err = WriteIndexedSummary(s.NewIndexedSummary,
		filepath.Join(testDir, SummaryFile))
	if err != nil {
		t.FailNow()
	}

	// write summary to directory with no write permissions
	os.Mkdir(filepath.Join(testDir, "noperm"), 0000)
	err = WriteIndexedSummary(s.NewIndexedSummary,
		filepath.Join(testDir, "noperm", SummaryFile))
	if err == nil {
		t.FailNow()
	}
}
