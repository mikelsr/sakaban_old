package fs

import (
	"testing"

	"github.com/satori/go.uuid"
)

// TestCommonRoot creates a tree of file changes and verifies that the common
// ancestor is correctly identified
func TestCommonRoot(t *testing.T) {
	// line 1:	f1	-> f2	-> f3
	//			|->f4	-> f5
	// line 2:	f6
	id, _ := uuid.NewV1()
	f1 := &File{ID: id, Blocks: make([]*Block, 0)}
	id, _ = uuid.NewV1()
	f2 := &File{ID: id, Parent: f1.ID}
	id, _ = uuid.NewV1()
	f3 := &File{ID: id, Parent: f2.ID}
	id, _ = uuid.NewV1()
	f4 := &File{ID: id, Parent: f1.ID}
	id, _ = uuid.NewV1()
	f5 := &File{ID: id, Parent: f4.ID}
	id, _ = uuid.NewV1()
	f6 := &File{ID: id}

	s1 := MakeSummary(f1)
	s2 := MakeSummary(f2)
	s3 := MakeSummary(f3)
	s4 := MakeSummary(f4)
	s5 := MakeSummary(f5)
	s6 := MakeSummary(f6)

	is1, _ := MakeIndexedSummary(s3)
	is1.AddParent(s1, s2)
	is2, _ := MakeIndexedSummary(s5)
	is2.AddParent(s1, s4)

	// they have a corrent ancestor
	if !commonRoot(s3, s5, is1.Parents, is2.Parents) {
		t.FailNow()
	}

	is1.Add(s6)
	is2.Add(s6)

	// they don't have a corrent ancestor
	if commonRoot(s3, s6, is1.Parents, is2.Parents) {
		t.FailNow()
	}
}
