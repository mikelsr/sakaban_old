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

// TestSummary_Equals checks that true is returned when the content is the same
// and false when it is different
func TestSummary_Equals(t *testing.T) {
	id, _ := uuid.NewV4()
	f1 := &File{ID: id, Path: "/path_1", Blocks: []*Block{&Block{[]byte{0}}}}
	s1 := MakeSummary(f1)

	id, _ = uuid.NewV4()
	f1.ID = id
	f1.Path = "/path_2"
	s2 := MakeSummary(f1)

	id, _ = uuid.NewV4()
	f2 := &File{ID: id, Path: "/path_3", Blocks: []*Block{&Block{[]byte{1}}}}
	s3 := MakeSummary(f2)

	if !s1.Equals(s2) || s1.Equals(s3) {
		t.FailNow()
	}
}
