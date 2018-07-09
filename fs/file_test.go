package fs

import (
	"fmt"
	"os"
	"testing"

	"github.com/satori/go.uuid"
)

var muffinPath = fmt.Sprintf("%s/res/muffin.jpg", ProjectPath())

// TestMakeFile builds File from a valid and an invalid path
func TestMakeFile(t *testing.T) {
	// Correct file
	_, err := MakeFile(muffinPath)
	if err != nil {
		panic(err)
	}

	// Incorrect file
	if _, err = MakeFile(""); err == nil {
		t.Fatalf("Built File from empty path")
	}
}

// TestMakeFileSummary checks the replicability of the
// File->Summary->File->Summary construction cycle
func TestMakeFileFromSummary(t *testing.T) {
	f, _ := MakeFile(muffinPath)
	parent, _ := uuid.NewV1()
	f.Parent = &parent
	fSum := MakeFileSummary(f)
	f2, err := MakeFileFromSummary(fSum)
	if err != nil {
		t.FailNow()
	}

	// same FileSummary
	fSum2 := MakeFileSummary(f2)
	if !fSum.Is(fSum2) {
		t.FailNow()
	}

	// different ID and amount of Blocks
	f2.ID, _ = uuid.NewV1()
	fSum2 = MakeFileSummary(f2)
	fSum2.Blocks = []uint64{0, 1}
	if fSum.Equals(fSum2) || fSum.Is(fSum2) {
		t.FailNow()
	}

	// invalid ID
	fSum2.ID = "invalid uuid"
	_, err = MakeFileFromSummary(fSum2)
	if err == nil {
		t.FailNow()
	}

	// invalid parent ID
	fSum2.ID = fSum2.Parent
	fSum2.Parent = "invalid uuid"
	_, err = MakeFileFromSummary(fSum2)
	if err == nil {
		t.FailNow()
	}

	// invalid path
	fSum2.Parent = fSum2.ID
	fSum2.Path = ""
	_, err = MakeFileFromSummary(fSum2)
	if err == nil {
		t.FailNow()
	}

	// different blocks, same amount
	fSum2 = MakeFileSummary(f)
	fSum2.Blocks = make([]uint64, len(fSum.Blocks))
	if fSum.Equals(fSum2) {
		t.FailNow()
	}
}

// TestMakeFileSummary checks that a FileSummary is built properly from a File
func TestMakeFileSummary(t *testing.T) {
	f, _ := MakeFile(muffinPath)
	fSum := MakeFileSummary(f)
	if fSum.ID != f.ID.String() || fSum.Parent != "" {
		t.FailNow()
	}

	parent, _ := uuid.NewV4()
	f.Parent = &parent
	fSum = MakeFileSummary(f)
	if fSum.Parent != parent.String() {
		t.FailNow()
	}
}

// TestFile_DeepEquals makes a shallow comparison between a File with itself,
// and another file with: the same blocks, different blocks, different number
// of blocks
func TestFile_DeepEquals(t *testing.T) {
	f1, _ := MakeFile(muffinPath)
	f2 := *f1
	if !f1.DeepEquals(f1) {
		t.FailNow()
	}
	if !f1.DeepEquals(&f2) {
		t.FailNow()
	}

	// check that first and last block are not equal
	if !f2.Blocks[0].DeepEquals(f2.Blocks[len(f2.Blocks)-1]) {
		// rearrange f2 blocks
		f2.Blocks = append(f2.Blocks[1:], f2.Blocks[0])
		if f1.DeepEquals(&f2) {
			t.FailNow()
		}
	}

	// change amount of blocks
	f2.Blocks = f2.Blocks[1:]
	if len(f1.Blocks) > 1 && f1.DeepEquals(&f2) {
		t.FailNow()
	}
}

// TestFile_Equals makes a shallow comparison between a File with itself,
// and another file with: the same blocks, different blocks, different number
// of blocks
func TestFile_Equals(t *testing.T) {
	f1, _ := MakeFile(muffinPath)
	f2 := *f1
	if !f1.Equals(f1) {
		t.FailNow()
	}
	if !f1.Equals(&f2) {
		t.FailNow()
	}

	// check that first and last block are not equal
	if !f2.Blocks[0].Equals(f2.Blocks[len(f2.Blocks)-1]) {
		// rearrange f2 blocks
		f2.Blocks = append(f2.Blocks[1:], f2.Blocks[0])
		if f1.Equals(&f2) {
			t.FailNow()
		}
	}

	// change amount of blocks
	f2.Blocks = f2.Blocks[1:]
	if len(f1.Blocks) > 1 && f1.Equals(&f2) {
		t.FailNow()
	}
}

// TestFile_Slice divides a valid and an invalid file into blocks,
// checks number of blocks in valid file
func TestFile_Slice(t *testing.T) {
	// Correct file
	f, _ := MakeFile(muffinPath)
	blocks, err := f.Slice()
	if err != nil {
		panic(err)
	}
	// Check that the file is sliced into the correct
	// amount of blocks
	file, _ := os.Stat(f.Path)
	blockN := CalcBlockN(file)

	if len(blocks) != blockN {
		t.Fatalf("Incorrect block number after slicing: got %d expected %d",
			len(blocks), blockN)
	}

	// Incorrect file
	f = &File{Path: ""}
	if _, err := f.Slice(); err == nil {
		t.Fatalf("Sliced non-existing file")
	}
}

// TestFile_String only ensures that the generated string is not empty
func TestFile_String(t *testing.T) {
	f, _ := MakeFile(muffinPath)
	if f.String() == "" {
		t.FailNow()
	}
}
