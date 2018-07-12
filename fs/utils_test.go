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
	id, _ := uuid.NewV4()
	f1 := &File{ID: id, Blocks: make([]*Block, 0)}
	id, _ = uuid.NewV4()
	f2 := &File{ID: id, Parent: f1.ID}
	id, _ = uuid.NewV4()
	f3 := &File{ID: id, Parent: f2.ID}
	id, _ = uuid.NewV4()
	f4 := &File{ID: id, Parent: f1.ID}
	id, _ = uuid.NewV4()
	f5 := &File{ID: id, Parent: f4.ID}
	id, _ = uuid.NewV4()
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

	parents, _ := mergeSummaryMap(true, is1.Parents, is2.Parents)

	// they have a corrent ancestor
	if !commonRoot(s3, s5, parents) {
		t.FailNow()
	}

	is1.Add(s6)
	is2.Add(s6)

	// they don't have a corrent ancestor
	if commonRoot(s3, s6, parents) {
		t.FailNow()
	}
}

// TestDescendant creates summaries and lines and checks that the verification
// is correct
func TestDescendant(t *testing.T) {
	// f1 -> f2 -> f3
	id, _ := uuid.NewV4()
	f1 := &File{ID: id}
	id, _ = uuid.NewV4()
	f2 := &File{ID: id, Parent: f1.ID}
	id, _ = uuid.NewV4()
	f3 := &File{ID: id, Parent: f2.ID}

	s1 := MakeSummary(f1)
	s2 := MakeSummary(f2)
	s3 := MakeSummary(f3)

	parents := make(map[string]*Summary)
	parents[s1.ID] = s1
	parents[s2.ID] = s2

	if !isDescendant(s3, s2, parents) {
		t.FailNow()
	}
	if !isDescendant(s3, s1, parents) {
		t.FailNow()
	}
	if !isDescendant(s2, s1, parents) {
		t.FailNow()
	}
	if isDescendant(s1, s2, parents) {
		t.FailNow()
	}
	if isDescendant(s2, s3, parents) {
		t.FailNow()
	}

}

func TestMergeSummaryMap(t *testing.T) {

	// no maps
	if _, err := mergeSummaryMap(true); err == nil {
		t.FailNow()
	}

	id, _ := uuid.NewV4()
	f1 := &File{ID: id, Blocks: []*Block{&Block{Content: []byte{0, 1}}}}
	f1_0 := &File{ID: id, Blocks: []*Block{&Block{Content: []byte{0, 1}}}}
	f1_1 := &File{ID: id, Blocks: []*Block{&Block{Content: []byte{1, 0}}}}
	id, _ = uuid.NewV4()
	f2 := &File{ID: id, Blocks: []*Block{&Block{Content: []byte{0, 1}}}}

	s1 := MakeSummary(f1)
	s1_0 := MakeSummary(f1_0)
	s1_1 := MakeSummary(f1_1)
	s2 := MakeSummary(f2)

	m1 := make(map[string]*Summary)
	m2 := make(map[string]*Summary)

	// Merge with two equal summaries
	m1[s1.ID] = s1
	m2[s1_0.ID] = s1_0
	if _, err := mergeSummaryMap(false, m1, m2); err != nil {
		t.FailNow()
	}

	// Merge with two different summaries with the same ID
	m2[s1_1.ID] = s1_1
	// 	conflict
	if _, err := mergeSummaryMap(false, m1, m2); err == nil {
		t.FailNow()
	}
	//	ignore conflict
	if _, err := mergeSummaryMap(true, m1, m2); err != nil {
		t.FailNow()
	}
	delete(m2, s1_1.ID)

	m2[s2.ID] = s2
	if _, err := mergeSummaryMap(false, m1, m2); err != nil {
		t.FailNow()
	}

}
