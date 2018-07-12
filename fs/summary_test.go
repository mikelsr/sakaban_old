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
		t.Fatal(err)
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
	f.Parent = parent
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
		t.Fatal(err)
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
		t.Fatal(err)
	}
	// repeated addition
	err = is.AddParent(s)
	if err == nil {
		t.FailNow()
	}
}

// TestIndexedSummary_Contains checks that an IndexedSummary contains a
// summary and doesn't contain another
func TestIndexedSummary_Contains(t *testing.T) {
	id, _ := uuid.NewV4()
	f1 := &File{ID: id, Path: "1", Blocks: []*Block{&Block{Content: []byte{0, 1}}}}
	id, _ = uuid.NewV4()
	f2 := &File{ID: id, Path: "2", Blocks: []*Block{&Block{Content: []byte{0, 1}}}}
	id, _ = uuid.NewV4()
	f3 := &File{ID: id, Path: "3", Blocks: []*Block{&Block{Content: []byte{1, 0}}}}

	s1 := MakeSummary(f1)
	s2 := MakeSummary(f2)
	s3 := MakeSummary(f3)

	is, _ := MakeIndexedSummary(s1)

	if path, found := is.Contains(s2); !found || path != s1.Path {
		t.FailNow()
	}

	if path, found := is.Contains(s3); found || path != "" {
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
		t.Fatal(err)
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
		t.Fatal(err)
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
	if _, found := is3.Deletions["f3.0"]; !found {
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

// TestMerge checks that the following merge operations are successfully carried
// out:
//	No changes
//	Two branches of a same file
//	Edit a file in one branch, move it in the other
//	Delete a file
//	Delete a file in one branch, edit it in another
//	Delete a file in one branch, move it in another
func TestMerge(t *testing.T) {
	testMerge1(t) // no changes
	testMerge2(t) // two branches of a same file
	testMerge3(t) // edit a file in one branch, move it in the other
	testMerge4(t) // delete a file
	testMerge5(t) // delete a file in one branch, edit it in another
	testMerge6(t) // delete a file in one branch, move it in another
}

// testMerge1 merges the same file
func testMerge1(t *testing.T) {
	id, _ := uuid.NewV4()
	f1_0ab := &File{ID: id, Path: "/path_1", Blocks: []*Block{&Block{Content: []byte{0}}}}
	s1_0ab := MakeSummary(f1_0ab)

	is1, _ := MakeIndexedSummary(s1_0ab)
	is2, _ := MakeIndexedSummary(s1_0ab)

	// equal summaries
	is3, err := Merge(is1, is2)
	if err != nil {
		t.Fatal(err)
	}
	if len(is3.Files) != 1 || len(is3.Parents) != 0 || len(is3.Deletions) != 0 {
		t.FailNow()
	}

	if !is3.Files[s1_0ab.Path].Equals(s1_0ab) {
		t.FailNow()
	}
}

// testMerge2 merges different branches of a same file
func testMerge2(t *testing.T) {
	id, _ := uuid.NewV4()
	f1_0ab := &File{ID: id, Path: "/path_1", Blocks: []*Block{&Block{Content: []byte{0}}}}
	s1_0ab := MakeSummary(f1_0ab)

	id, _ = uuid.NewV4()
	f1_1a := &File{ID: id, Parent: f1_0ab.ID, Path: "/path_1", Blocks: []*Block{&Block{Content: []byte{0, 1}}}}
	id, _ = uuid.NewV4()
	f1_2a := &File{ID: id, Parent: f1_1a.ID, Path: "/path_1", Blocks: []*Block{&Block{Content: []byte{0, 1}}}}
	id, _ = uuid.NewV4()
	f1_1b := &File{ID: id, Parent: f1_0ab.ID, Path: "/path_1", Blocks: []*Block{&Block{Content: []byte{0, 2}}}}

	s1_1a := MakeSummary(f1_1a)
	s1_2a := MakeSummary(f1_2a)
	s1_1b := MakeSummary(f1_1b)

	is1, _ := MakeIndexedSummary(s1_2a)
	is2, _ := MakeIndexedSummary(s1_1b)
	is1.AddParent(s1_0ab, s1_1a)
	is2.AddParent(s1_0ab)

	is3, err := Merge(is1, is2)
	if err != nil {
		t.Fatal(err)
	}
	if len(is3.Files) != 2 || len(is3.Parents) != 2 || len(is3.Deletions) != 0 {
		t.FailNow()
	}
}

// testMerge3 merges an edited and a moved branch of a file
func testMerge3(t *testing.T) {
	id, _ := uuid.NewV4()
	f1_0ab := &File{ID: id, Path: "/path_1", Blocks: []*Block{&Block{Content: []byte{0}}}}
	s1_0ab := MakeSummary(f1_0ab)

	id, _ = uuid.NewV4()
	f1_1a := &File{ID: id, Parent: f1_0ab.ID, Path: "/path_1", Blocks: []*Block{&Block{Content: []byte{0, 1}}}}
	id, _ = uuid.NewV4()
	f1_2a := &File{ID: id, Parent: f1_1a.ID, Path: "/path_2", Blocks: []*Block{&Block{Content: []byte{0, 1}}}}
	id, _ = uuid.NewV4()
	f1_1b := &File{ID: id, Parent: f1_0ab.ID, Path: "/path_1", Blocks: []*Block{&Block{Content: []byte{0, 2}}}}

	s1_1a := MakeSummary(f1_1a)
	s1_2a := MakeSummary(f1_2a)
	s1_1b := MakeSummary(f1_1b)

	is1, _ := MakeIndexedSummary(s1_2a)
	is2, _ := MakeIndexedSummary(s1_1b)
	is1.AddParent(s1_0ab, s1_1a)
	is2.AddParent(s1_0ab)

	is3, err := Merge(is1, is2)
	if err != nil {
		t.Fatal(err)
	}
	if len(is3.Files) != 2 || len(is3.Parents) != 2 || len(is3.Deletions) != 0 {
		t.FailNow()
	}
	if _, found := is3.Files[s1_2a.Path]; !found {
		t.FailNow()
	}
}

// testMerge4 creates a file in both branches and deletes it in one of them
func testMerge4(t *testing.T) {
	id, _ := uuid.NewV4()
	f1 := &File{ID: id, Path: "/path_1", Blocks: []*Block{&Block{Content: []byte{0}}}}
	s1 := MakeSummary(f1)

	is1, _ := MakeIndexedSummary(s1)
	is2, _ := MakeIndexedSummary()
	// TODO: deletion function for IndexedSummary struct
	is2.Deletions[s1.ID] = s1

	is3, err := Merge(is1, is2)
	if err != nil {
		t.Fatal(err)
	}
	if len(is3.Files) != 0 || len(is3.Parents) != 0 || len(is3.Deletions) != 1 {
		t.FailNow()
	}
}

// testMerge5 deletes a file in one branch, edits it in another
func testMerge5(t *testing.T) {
	id, _ := uuid.NewV4()
	f1_0 := &File{ID: id, Path: "/path_1", Blocks: []*Block{&Block{Content: []byte{0}}}}
	id, _ = uuid.NewV4()
	f1_1 := &File{ID: id, Parent: f1_0.ID, Path: "/path_1", Blocks: []*Block{&Block{Content: []byte{1}}}}
	id, _ = uuid.NewV4()
	f1_2 := &File{ID: id, Parent: f1_1.ID, Path: "/path_1", Blocks: []*Block{&Block{Content: []byte{2}}}}
	s1_0 := MakeSummary(f1_0)
	s1_1 := MakeSummary(f1_1)
	s1_2 := MakeSummary(f1_2)

	is1, _ := MakeIndexedSummary(s1_2)
	is2, _ := MakeIndexedSummary()

	is1.AddParent(s1_0, s1_1)
	is2.Deletions[s1_0.ID] = s1_0

	is3, err := Merge(is1, is2)
	if err != nil {
		t.Fatal(err)
	}
	if len(is3.Files) != 1 || len(is3.Parents) != 2 || len(is3.Deletions) != 0 {
		t.FailNow()
	}
}

// testMerge6 deletes a file in one branch, moves it in another
func testMerge6(t *testing.T) {
	id, _ := uuid.NewV4()
	f1_0 := &File{ID: id, Path: "/path_1", Blocks: []*Block{&Block{Content: []byte{0}}}}
	id, _ = uuid.NewV4()
	f1_1 := &File{ID: id, Parent: f1_0.ID, Path: "/path_2", Blocks: []*Block{&Block{Content: []byte{0}}}}
	id, _ = uuid.NewV4()
	f1_2 := &File{ID: id, Parent: f1_1.ID, Path: "/path_3", Blocks: []*Block{&Block{Content: []byte{0}}}}
	s1_0 := MakeSummary(f1_0)
	s1_1 := MakeSummary(f1_1)
	s1_2 := MakeSummary(f1_2)

	is1, _ := MakeIndexedSummary(s1_2)
	is2, _ := MakeIndexedSummary()

	is1.AddParent(s1_0, s1_1)
	is2.Deletions[s1_0.ID] = s1_0

	is3, err := Merge(is1, is2)
	if err != nil {
		t.Fatal(err)
	}
	if len(is3.Files) != 1 || len(is3.Parents) != 2 || len(is3.Deletions) != 0 {
		t.FailNow()
	}
}
