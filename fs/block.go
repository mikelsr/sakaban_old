package fs

import (
	"bytes"
	"hash/fnv" // non-cryptographic hash functions
)

// Block contains a portion of a file and the hash corresponding
// to that portion
//	Content: bytes forming the block of the file
type Block struct {
	Content []byte
}

// DeepEquals compares the contents of the body of a block,
// avoiding false positives caused by hash collisions
func (b *Block) DeepEquals(b2 *Block) bool {
	if b == b2 {
		return true
	}
	if b.Hash() != b2.Hash() {
		return false
	}
	return bytes.Equal(b.Content, b2.Content)
}

// Equals compares two blocks by comparing their hashes, vulnerable
// to hash collisions
func (b *Block) Equals(b2 *Block) bool {
	if b == b2 {
		return true
	}
	return b.Hash() == b2.Hash()
}

// Hash uses the hash/fnv package to generate the hash of the Block Content
func (b *Block) Hash() uint64 {
	hfn := fnv.New64a()
	hfn.Write(b.Content)
	return hfn.Sum64()
}

// Size returns the number of bytes in the content of a Block
func (b *Block) Size() int {
	return len(b.Content)
}
