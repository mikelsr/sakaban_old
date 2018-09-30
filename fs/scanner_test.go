package fs

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// testDir will contain the files generated for this tests
var testDir string

// TestMain will create and delete the testing directory
// before and after running tests, respectively
func TestMain(m *testing.M) {
	rand.Seed(time.Now().UTC().UnixNano())
	testDir = filepath.Join(ProjectPath(), "test",
		fmt.Sprintf("sakaban-test-%d", rand.Intn(1e8)))
	os.MkdirAll(testDir, 0770)
	defer os.RemoveAll(testDir)
	m.Run()
}

func TestMakeScanner(t *testing.T) {
	filename := filepath.Join(testDir, "MakeScanner", SummaryDir, SummaryFile)
	unitTestDir := filepath.Join(testDir, "MakeScanner")
	os.MkdirAll(filepath.Join(unitTestDir, SummaryDir), 0755)
	// scanner with no old indexed summary
	s, err := MakeScanner(filepath.Join(ProjectPath(), "res"))
	if err != nil {
		t.FailNow()
	}
	// scanner with old indexed summary
	WriteIndex(*s.NewIndex, filename)
	s, err = MakeScanner(unitTestDir)
	if err != nil || len(s.OldIndex.Files) < 1 {
		t.FailNow()
	}

	// existing but incorrect old indexed summary
	ioutil.WriteFile(filename, []byte{42}, 0755)
	_, err = MakeScanner(unitTestDir)
	if err == nil {
		t.FailNow()
	}

	// create scanner on non-readable folder
	os.Chmod(unitTestDir, 0000)
	_, err = MakeScanner(unitTestDir)
	if err == nil {
		t.FailNow()
	}
}

// TestScanner_Scan runs Scanner.Scan in the repository resource dir and
// compares the number of generated files with the number of files in the dir
// It is also used to test Scanner.Visit
func TestScanner_Scan(t *testing.T) {
	resFolder := filepath.Join(ProjectPath(), "res")
	unitTestDir := filepath.Join(testDir, "Scanner_Scan")
	scanner := new(Scanner)
	scanner.Scan(resFolder)
	resFiles, _ := ioutil.ReadDir(resFolder)
	if len(scanner.Summaries) != len(resFiles) {
		t.FailNow()
	}

	// scan folder with no read permission
	os.MkdirAll(unitTestDir, 0000)
	err := scanner.Scan(unitTestDir)
	if err == nil {
		t.FailNow()
	}
}

// TestSummaryExists ensures SummaryExists properly detects the summary
func TestSummaryExists(t *testing.T) {
	if SummaryExists(testDir) {
		t.FailNow()
	}
	err := os.MkdirAll(filepath.Join(testDir, SummaryDir), 0755)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.OpenFile(filepath.Join(testDir, SummaryDir, SummaryFile),
		os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		t.Fatal(err)
	}
	if !SummaryExists(testDir) {
		t.FailNow()
	}
}
