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

func TestSummary_Diff(t *testing.T) {
	s1 := Summary{Blocks: []uint64{1, 3, 2}}
	s2 := Summary{Blocks: []uint64{1, 2, 3, 4}}
	expectedDiff := []uint64{0, 2, 3, 4}
	diff, change := s1.Diff(&s2)
	if !change {
		t.FailNow()
	}
	for i, block := range diff {
		if block != expectedDiff[i] {
			t.FailNow()
		}
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
