package supervisor

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"

	"bitbucket.org/mikelsr/sakaban/fs"
)

// testDir will contain the files generated for this tests
var testDir = filepath.Join(fs.ProjectPath(), "test", fmt.Sprintf("sakaban-test-%d", rand.Intn(1e8)))

// TestScanner_Scan runs Scanner.Scan in the repository resource dir and
// compares the number of generated files with the number of files in the dir
// It is also used to test Scanner.Visit
func TestScanner_Scan(t *testing.T) {
	resFolder := filepath.Join(fs.ProjectPath(), "res")
	scanner := new(Scanner)
	scanner.Scan(resFolder)
	resFiles, _ := ioutil.ReadDir(resFolder)
	if len(scanner.Files) != len(resFiles) {
		t.FailNow()
	}
}

// TestSummaryExists ensures SummaryExists properly detects the summary
func TestSummaryExists(t *testing.T) {
	if SummaryExists(testDir) {
		t.FailNow()
	}
	err := os.MkdirAll(filepath.Join(testDir, SummaryDir), 0777)
	defer os.RemoveAll(testDir)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.OpenFile(filepath.Join(testDir, SummaryDir, SummaryFile), os.O_RDWR|os.O_CREATE, 0777)
	if err != nil {
		t.Fatal(err)
	}
	if !SummaryExists(testDir) {
		t.FailNow()
	}
}
