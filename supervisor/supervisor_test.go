package supervisor

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

var testDir = fmt.Sprintf("%s/tmp/sakaban-test-%d", os.Getenv("HOME"), rand.Intn(1000))

func TestScan(t *testing.T) {
	if SummaryExists(testDir) {
		t.FailNow()
	}
	err := os.MkdirAll(fmt.Sprintf("%s/%s", testDir, SummaryDir), 0666)
	defer os.RemoveAll(testDir)

	if err != nil && err == os.ErrPermission {
		t.SkipNow()
	}
	_, err = os.Create(fmt.Sprintf("%s/%s/%s", testDir, SummaryDir, SummaryFile))
	if err != nil && err == os.ErrPermission {
		t.SkipNow()
	}
	if !SummaryExists(testDir) {
		t.FailNow()
	}
}
