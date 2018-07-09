package fs

import (
	"fmt"
	"os"
	"testing"
)

var projectPath = fmt.Sprintf("%s/src/bitbucket.org/mikelsr/sakaban", os.Getenv("GOPATH"))
var muffinPath = fmt.Sprintf("%s/res/muffin.jpg", projectPath)

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
