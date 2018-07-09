package fs

import (
	"fmt"
	"os"
	"testing"
)

var projectPath = fmt.Sprintf("%s/src/bitbucket.org/mikelsr/sakaban", os.Getenv("GOPATH"))
var filePath = fmt.Sprintf("%s/res/muffin.jpg", projectPath)

// TestMakeFile builds File from a valid and an invalid path
func TestMakeFile(t *testing.T) {
	// Correct file
	_, err := MakeFile(filePath)
	if err != nil {
		panic(err)
	}

	// Incorrect file
	if _, err = MakeFile(""); err == nil {
		t.Fatalf("Built File from empty path")
	}
}

// TestFile_Slice divides a valid and an invalid file into blocks,
// checks number of blocks in valid file
func TestFile_Slice(t *testing.T) {
	// Correct file
	f, _ := MakeFile(filePath)
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
