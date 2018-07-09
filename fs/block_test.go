package fs

import (
	"testing"
)

// TestBlock_Equals makes a shallow comparison of a block with itself
// and with a different block
func TestBlock_Equals(t *testing.T) {
	b1 := Block{Content: []byte{1}}
	b2 := Block{Content: []byte{1, 2}}
	if !b1.Equals(&b1) {
		t.FailNow()
	}
	if b1.Equals(&b2) {
		t.FailNow()
	}
}

// TestBlock_DeepEquals compares the content of a Block with itself,
// an equal block and a different block
func TestBlock_DeepEquals(t *testing.T) {
	b1 := Block{Content: []byte{1}}
	b2 := Block{Content: []byte{1}}
	b3 := Block{Content: []byte{1, 2}}
	if !b1.DeepEquals(&b1) {
		t.FailNow()
	}
	if !b1.DeepEquals(&b2) {
		t.FailNow()
	}
	if b1.DeepEquals(&b3) {
		t.FailNow()
	}
}

// TestBlock_Hash checks for hash collisions in different blocks
func TestBlock_Hash(t *testing.T) {
	// Contents of different lenght
	b1 := Block{Content: []byte{1}}
	b2 := Block{Content: []byte{1, 2}}
	if b1.Hash() == b2.Hash() {
		t.Fatalf("Hash collision")
	}
	// Contents of same lenght
	b1.Content = []byte{1, 2, 3}
	b2.Content = []byte{3, 2, 1}
	if b1.Hash() == b2.Hash() {
		t.Fatalf("Hash collision")
	}
}

// TestBlock_Size compares the output of Block.Size to the lenght
// of Block.Content
func TestBlock_Size(t *testing.T) {
	content := []byte{1, 2, 3}
	b := Block{Content: content}
	if b.Size() != len(content) {
		t.Fatalf("Failed to calculate block size, expected %d got %d",
			len(content), b.Size())
	}
}
