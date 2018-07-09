package fs

import (
	"testing"
)

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
