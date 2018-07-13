package fs

import (
	"testing"

	"github.com/satori/go.uuid"
)

// TestMakeSummary checks that a Summary is built properly from a File
func TestMakeSummary(t *testing.T) {
	f, _ := MakeFile(muffinPath)
	s := MakeSummary(f)
	if s.ID != f.ID.String() || s.Parent != "" {
		t.FailNow()
	}

	parent, _ := uuid.NewV4()
	f.Parent = parent
	s = MakeSummary(f)
	if s.Parent != parent.String() {
		t.FailNow()
	}
}
