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
