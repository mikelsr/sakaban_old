package fs

import (
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
)

// TestMakeIndex gives a valid and an invalid set to the
// Index constructor
func TestMakeIndex(t *testing.T) {
	s1 := &Summary{ID: "id", Path: "/path", Blocks: []uint64{0}}
	s2 := *s1
	s2.Path = "/s2/path"
	i, err := MakeIndex(s1, &s2)
	if err != nil {
		t.Fatal(err)
	}
	if !i.Files[s1.Path].Equals(s1) {
		t.FailNow()
	}
	s2.Path = s1.Path
	_, err = MakeIndex(s1, &s2)
	if err == nil {
		t.FailNow()
	}
}

// TestIndex_Add adds a new and a repeated summary to the
// Index
func TestIndex_Add(t *testing.T) {
	s := &Summary{ID: "id", Path: "/path", Blocks: []uint64{0}}
	i, _ := MakeIndex()
	// new addition
	err := i.Add(s)
	if err != nil {
		t.Fatal(err)
	}
	// repeated addition
	err = i.Add(s)
	if err == nil {
		t.FailNow()
	}
}

// TestIndex_AddDeletion adds a new and a repeated deletion to the
// Index
func TestIndex_AddDeletion(t *testing.T) {
	s := &Summary{ID: "id", Path: "/path", Blocks: []uint64{0}}
	i, _ := MakeIndex()
	// new addition
	err := i.AddDeletion(s)
	if err != nil {
		t.Fatal(err)
	}
	// repeated addition
	err = i.AddDeletion(s)
	if err == nil {
		t.FailNow()
	}
}

// TestIndex_AddParent adds a new and a repeated parent to the
// Index
func TestIndex_AddParent(t *testing.T) {
	s := &Summary{ID: "id", Path: "/path", Blocks: []uint64{0}}
	i, _ := MakeIndex()
	// new addition
	err := i.AddParent(s)
	if err != nil {
		t.Fatal(err)
	}
	// repeated addition
	err = i.AddParent(s)
	if err == nil {
		t.FailNow()
	}
}

func TestIndex_Compare(t *testing.T) {
	id1, _ := uuid.NewV4()
	id2, _ := uuid.NewV4()
	id3, _ := uuid.NewV4()
	sum1 := Summary{ID: id1.String(), Path: "/1", Blocks: []uint64{1, 2, 3}}
	sum2 := Summary{ID: id2.String(), Path: "/2", Blocks: []uint64{4, 5, 6}}
	sum3 := Summary{ID: id3.String(), Path: "/3", Blocks: []uint64{7, 8, 9}}

	sum2_1 := Summary{ID: id2.String(), Path: "/2", Blocks: []uint64{5, 4, 6}}

	index1, _ := MakeIndex(&sum1, &sum2, &sum3)
	index2, _ := MakeIndex(&sum1, &sum2_1)
	index2.Deletions[sum3.ID] = &sum3

	expected := new(Comparison)
	expected.Additions = make(map[string]*Summary)
	expected.Additions["/2"] = &Summary{Blocks: []uint64{5, 4, 0}}
	expected.Deletions = []string{sum3.Path}

	comparison := index1.Compare(index2)

	if !reflect.DeepEqual(expected, comparison) {
		t.FailNow()
	}
}

// TestIndex_Contains checks that an Index contains a
// summary and doesn't contain another
func TestIndex_Contains(t *testing.T) {
	id, _ := uuid.NewV4()
	f1 := &File{ID: id, Path: "1", Blocks: []*Block{&Block{Content: []byte{0, 1}}}}
	id, _ = uuid.NewV4()
	f2 := &File{ID: id, Path: "2", Blocks: []*Block{&Block{Content: []byte{0, 1}}}}
	id, _ = uuid.NewV4()
	f3 := &File{ID: id, Path: "3", Blocks: []*Block{&Block{Content: []byte{1, 0}}}}

	s1 := MakeSummary(f1)
	s2 := MakeSummary(f2)
	s3 := MakeSummary(f3)

	i, _ := MakeIndex(s1)

	if path, found := i.Contains(s2); !found || path != s1.Path {
		t.FailNow()
	}

	if path, found := i.Contains(s3); found || path != "" {
		t.FailNow()
	}
}

// TestIndex_Delete deletes an existing and a nonexisting summary
// from the Index
func TestIndex_Delete(t *testing.T) {
	s := &Summary{ID: "id", Path: "/path", Blocks: []uint64{0}}
	i, _ := MakeIndex(s)
	err := i.Delete(s)
	if err != nil {
		t.Fatal(err)
	}
	err = i.Delete(s)
	if err == nil {
		t.FailNow()
	}
}

// TestIndex_DeleteDeletion deletes an existing and a nonexisting deletion
// from the Index
func TestIndex_DeleteDeletion(t *testing.T) {
	s := &Summary{ID: "id", Path: "/path", Blocks: []uint64{0}}
	i, _ := MakeIndex()
	i.AddDeletion(s)
	err := i.DeleteDeletion(s)
	if err != nil {
		t.Fatal(err)
	}
	err = i.DeleteDeletion(s)
	if err == nil {
		t.FailNow()
	}
}

// TestIndex_DeleteParent deletes an existing and a nonexisting parent
// from the Index
func TestIndex_DeleteParent(t *testing.T) {
	s := &Summary{ID: "id", Path: "/path", Blocks: []uint64{0}}
	i, _ := MakeIndex()
	i.AddParent(s)
	err := i.DeleteParent(s)
	if err != nil {
		t.Fatal(err)
	}
	err = i.DeleteParent(s)
	if err == nil {
		t.FailNow()
	}
}

// TestEquals compares indices with different and equal attributes
func TestIndex_Equals(t *testing.T) {
	s1 := &Summary{ID: "f1.0", Path: "/f1", Blocks: []uint64{1}}
	s2 := &Summary{ID: "f2.0", Path: "/f2", Blocks: []uint64{2}}
	s3 := &Summary{ID: "f3.0", Path: "/f3", Blocks: []uint64{3}}

	i1, _ := MakeIndex(s1)
	i1.AddParent(s2)
	i1.AddDeletion(s3)

	i2, _ := MakeIndex()

	// Different number of files/deletions/parents
	if i1.Equals(i2) {
		t.FailNow()
	}

	// different files
	i2.Add(s2)
	i2.AddParent(s1)
	i2.AddDeletion(s1)
	if i1.Equals(i2) {
		t.FailNow()
	}

	// different parents
	i2.Delete(s2)
	i2.Add(s1)
	if i1.Equals(i2) {
		t.FailNow()
	}

	// different deletions
	i2.Delete(s2)
	i2.Add(s1)
	i2.DeleteParent(s1)
	i2.AddParent(s2)
	if i1.Equals(i2) {
		t.FailNow()
	}

	// equal
	i2.DeleteDeletion(s1)
	i2.AddDeletion(s3)
	if !i1.Equals(i2) {
		t.FailNow()
	}

}

// TestIndex_Update creates and updates an Index,
// checking the operations: change, move, delete, keep, create
func TestIndex_Update(t *testing.T) {
	i1, _ := MakeIndex()
	i1.Add(&Summary{ID: "f1.0", Path: "/f1", Blocks: []uint64{1}},
		&Summary{ID: "f2.0", Path: "/f2", Blocks: []uint64{2}},
		&Summary{ID: "f3.0", Path: "/f3", Blocks: []uint64{3}},
		&Summary{ID: "f4.0", Path: "/f4", Blocks: []uint64{4}})
	i2, _ := MakeIndex()
	i2.Files = make(map[string]*Summary)
	i2.Add(&Summary{ID: "f1.1", Path: "/f1", Blocks: []uint64{11}}, // change
		&Summary{ID: "f2.2", Path: "/n2", Blocks: []uint64{2}}, // move
		&Summary{ID: "f4.0", Path: "/f4", Blocks: []uint64{4}}, // keep
		&Summary{ID: "f5.0", Path: "/f5", Blocks: []uint64{4}}) // create

	i3 := Update(i1, i2)

	// change
	if i3.Files["/f1"].Parent != i1.Files["/f1"].ID {
		t.FailNow()
	}
	// move
	if i3.Files["/n2"].Parent != i1.Files["/f2"].ID ||
		i3.Files["/n2"].Path == i1.Files["/f2"].Path {
		t.FailNow()
	}
	// delete
	if _, found := i3.Files["/f3"]; found {
		t.FailNow()
	}
	if _, found := i3.Deletions["f3.0"]; !found {
		t.FailNow()
	}
	// keep
	if !i3.Files["/f4"].Equals(i1.Files["/f4"]) {
		t.FailNow()
	}
	// create
	if _, found := i3.Files["/f5"]; !found {
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
	testMerge7(t) // compare merges
}

// testMerge1 merges the same file
func testMerge1(t *testing.T) {
	s := &Summary{ID: "id", Path: "/path", Blocks: []uint64{0}}

	i1, _ := MakeIndex(s)
	i2, _ := MakeIndex(s)

	// equal summaries
	i3, err := Merge(i1, i2)
	if err != nil {
		t.Fatal(err)
	}
	if len(i3.Files) != 1 || len(i3.Parents) != 0 || len(i3.Deletions) != 0 {
		t.FailNow()
	}

	if !i3.Files[s.Path].Equals(s) {
		t.FailNow()
	}
}

// testMerge2 merges different branches of a same file
func testMerge2(t *testing.T) {
	s1_0ab := &Summary{ID: "id1", Path: "/path_1", Blocks: []uint64{0}}
	s1_1a := &Summary{ID: "id2", Parent: s1_0ab.ID, Path: "/path_1", Blocks: []uint64{1}}
	s1_2a := &Summary{ID: "id3", Parent: s1_1a.ID, Path: "/path_1", Blocks: []uint64{2}}
	s1_1b := &Summary{ID: "id4", Parent: s1_0ab.ID, Path: "/path_1", Blocks: []uint64{3}}

	i1, _ := MakeIndex(s1_2a)
	i2, _ := MakeIndex(s1_1b)
	i1.AddParent(s1_0ab, s1_1a)
	i2.AddParent(s1_0ab)

	i3, err := Merge(i1, i2)
	if err != nil {
		t.Fatal(err)
	}
	if len(i3.Files) != 2 || len(i3.Parents) != 2 || len(i3.Deletions) != 0 {
		t.FailNow()
	}
}

// testMerge3 merges an edited and a moved branch of a file
func testMerge3(t *testing.T) {
	s1_0ab := &Summary{ID: "id1", Path: "/path_1", Blocks: []uint64{0}}
	s1_1a := &Summary{ID: "id2", Parent: s1_0ab.ID, Path: "/path_1", Blocks: []uint64{1}}
	s1_2a := &Summary{ID: "id3", Parent: s1_1a.ID, Path: "/path_2", Blocks: []uint64{1}}
	s1_1b := &Summary{ID: "id4", Parent: s1_0ab.ID, Path: "/path_1", Blocks: []uint64{2}}

	i1, _ := MakeIndex(s1_2a)
	i2, _ := MakeIndex(s1_1b)
	i1.AddParent(s1_0ab, s1_1a)
	i2.AddParent(s1_0ab)

	i3, err := Merge(i1, i2)
	if err != nil {
		t.Fatal(err)
	}
	if len(i3.Files) != 2 || len(i3.Parents) != 2 || len(i3.Deletions) != 0 {
		t.FailNow()
	}
	if _, found := i3.Files[s1_2a.Path]; !found {
		t.FailNow()
	}
}

// testMerge4 creates a file in both branches and deletes it in one of them
func testMerge4(t *testing.T) {
	s := &Summary{ID: "id", Path: "/path", Blocks: []uint64{0}}

	i1, _ := MakeIndex(s)
	i2, _ := MakeIndex()
	i2.AddDeletion(s)

	i3, err := Merge(i1, i2)
	if err != nil {
		t.Fatal(err)
	}
	if len(i3.Files) != 0 || len(i3.Parents) != 0 || len(i3.Deletions) != 1 {
		t.FailNow()
	}
}

// testMerge5 deletes a file in one branch, edits it in another
func testMerge5(t *testing.T) {
	s1_0 := &Summary{ID: "id1", Path: "/path_1", Blocks: []uint64{0}}
	s1_1 := &Summary{ID: "id2", Parent: s1_0.ID, Path: "/path_1", Blocks: []uint64{1}}
	s1_2 := &Summary{ID: "id3", Parent: s1_1.ID, Path: "/path_1", Blocks: []uint64{2}}

	i1, _ := MakeIndex(s1_2)
	i2, _ := MakeIndex()

	i1.AddParent(s1_0, s1_1)
	i2.AddDeletion(s1_0)

	i3, err := Merge(i1, i2)
	if err != nil {
		t.Fatal(err)
	}
	if len(i3.Files) != 1 || len(i3.Parents) != 2 || len(i3.Deletions) != 0 {
		t.FailNow()
	}
}

// testMerge6 deletes a file in one branch, moves it in another
func testMerge6(t *testing.T) {
	s1_0 := &Summary{ID: "id1", Path: "/path_1", Blocks: []uint64{0}}
	s1_1 := &Summary{ID: "id2", Parent: s1_0.ID, Path: "/path_2", Blocks: []uint64{0}}
	s1_2 := &Summary{ID: "id3", Parent: s1_1.ID, Path: "/path_3", Blocks: []uint64{0}}

	i1, _ := MakeIndex(s1_2)
	i2, _ := MakeIndex()

	i1.AddParent(s1_0, s1_1)
	i2.AddDeletion(s1_0)

	i3, err := Merge(i1, i2)
	if err != nil {
		t.Fatal(err)
	}
	if len(i3.Files) != 1 || len(i3.Parents) != 2 || len(i3.Deletions) != 0 {
		t.FailNow()
	}
}

// testMerge7 merges two indices and compares two merges
func testMerge7(t *testing.T) {
	s1_0ab := &Summary{ID: "id1", Path: "/path_1", Blocks: []uint64{0}}
	s1_1a := &Summary{ID: "id2", Parent: s1_0ab.ID, Path: "/path_1", Blocks: []uint64{1}}
	s1_2a := &Summary{ID: "id3", Parent: s1_1a.ID, Path: "/path_1", Blocks: []uint64{2}}

	i1, _ := MakeIndex(s1_2a)
	i1.AddParent(s1_0ab, s1_1a)

	i2, _ := MakeIndex()
	i2.AddParent(s1_0ab)
	i2.AddDeletion(s1_1a)

	m1, err := Merge(i1, i2)
	if err != nil {
		t.FailNow()
	}

	m2, err := Merge(i2, i1)
	if err != nil {
		t.FailNow()
	}

	if !m1.Equals(m2) {
		t.FailNow()
	}

	if len(m1.Files) != 1 || len(m1.Parents) != 2 || len(m1.Deletions) != 0 {
		t.FailNow()
	}

	if file, found := m1.Files[s1_2a.Path]; !found || !file.Equals(s1_2a) {
		t.FailNow()
	}
}
