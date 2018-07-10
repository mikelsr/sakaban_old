package fs

import (
	"testing"

	"github.com/satori/go.uuid"
)

// TestMakeIndexedSummary gives a valid and an invalid set to the
// IndexedSummary constructor
func TestMakeIndexedSummary(t *testing.T) {
	f, _ := MakeFile(muffinPath)
	s1 := MakeSummary(f)
	s2 := *s1
	s2.Path = "/s2/Path"
	is, err := MakeIndexedSummary(s1, &s2)
	if err != nil {
		t.FailNow()
	}
	if !is.Files[s1.Path].Equals(s1) {
		t.FailNow()
	}
	s2.Path = s1.Path
	_, err = MakeIndexedSummary(s1, &s2)
	if err == nil {
		t.FailNow()
	}
}

// TestMakeSummary checks that a Summary is built properly from a File
func TestMakeSummary(t *testing.T) {
	f, _ := MakeFile(muffinPath)
	s := MakeSummary(f)
	if s.ID != f.ID.String() || s.Parent != "" {
		t.FailNow()
	}

	parent, _ := uuid.NewV4()
	f.Parent = &parent
	s = MakeSummary(f)
	if s.Parent != parent.String() {
		t.FailNow()
	}
}

// TestIndexedSummary_Add adds a new and a repeated summary to the
// IndexedSummary
func TestIndexedSummary_Add(t *testing.T) {
	f, _ := MakeFile(muffinPath)
	s := MakeSummary(f)
	is, _ := MakeIndexedSummary()
	// new addition
	err := is.Add(s)
	if err != nil {
		t.FailNow()
	}
	// repeated addition
	err = is.Add(s)
	if err == nil {
		t.FailNow()
	}
}

// TestIndexedSummary_AddParent adds a new and a repeated parent to the
// IndexedSummary
func TestIndexedSummary_AddParent(t *testing.T) {
	f, _ := MakeFile(muffinPath)
	s := MakeSummary(f)
	is, _ := MakeIndexedSummary()
	// new addition
	err := is.AddParent(s)
	if err != nil {
		t.FailNow()
	}
	// repeated addition
	err = is.AddParent(s)
	if err == nil {
		t.FailNow()
	}
}

// TestIndexedSummary_Delete deletes an existing and a nonexisting summary
// from the IndexedSummary
func TestIndexedSummary_Delete(t *testing.T) {
	f, _ := MakeFile(muffinPath)
	s := MakeSummary(f)
	is, _ := MakeIndexedSummary(s)
	err := is.Delete(s)
	if err != nil {
		t.FailNow()
	}
	err = is.Delete(s)
	if err == nil {
		t.FailNow()
	}
}

// TestIndexedSummary_DeleteParent deletes an existing and a nonexisting parent
// from the IndexedSummary
func TestIndexedSummary_DeleteParent(t *testing.T) {
	f, _ := MakeFile(muffinPath)
	s := MakeSummary(f)
	is, _ := MakeIndexedSummary()
	is.AddParent(s)
	err := is.DeleteParent(s)
	if err != nil {
		t.FailNow()
	}
	err = is.DeleteParent(s)
	if err == nil {
		t.FailNow()
	}
}

// TestIndexedSummary_Update creates and updates an IndexedSummary,
// checking the operations: change, move, delete, keep, create
func TestIndexedSummary_Update(t *testing.T) {
	is1, _ := MakeIndexedSummary()
	is1.Add(&Summary{ID: "f1.0", Path: "/f1", Blocks: []uint64{1}},
		&Summary{ID: "f2.0", Path: "/f2", Blocks: []uint64{2}},
		&Summary{ID: "f3.0", Path: "/f3", Blocks: []uint64{3}},
		&Summary{ID: "f4.0", Path: "/f4", Blocks: []uint64{4}})
	is2, _ := MakeIndexedSummary()
	is2.Files = make(map[string]*Summary)
	is2.Add(&Summary{ID: "f1.1", Path: "/f1", Blocks: []uint64{11}}, // change
		&Summary{ID: "f2.2", Path: "/n2", Blocks: []uint64{2}}, // move
		&Summary{ID: "f4.0", Path: "/f4", Blocks: []uint64{4}}, // keep
		&Summary{ID: "f5.0", Path: "/f5", Blocks: []uint64{4}}) // create

	is3 := is1.Update(is2)

	// change
	if is3.Files["/f1"].Parent != is1.Files["/f1"].ID {
		t.FailNow()
	}
	// move
	if is3.Files["/n2"].Parent != is1.Files["/f2"].ID ||
		is3.Files["/n2"].Path == is1.Files["/f2"].Path {
		t.FailNow()
	}
	// delete
	if _, found := is3.Files["/f3"]; found {
		t.FailNow()
	}
	// keep
	if !is3.Files["/f4"].Equals(is1.Files["/f4"]) {
		t.FailNow()
	}
	// create
	if _, found := is3.Files["/f5"]; !found {
		t.FailNow()
	}
}