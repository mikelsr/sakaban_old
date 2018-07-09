package fs

import (
	"hash/fnv" // non-cryptographic hash functions
)

// Block contains a portion of a file and the hash corresponding
// to that portion
//	Hash: hash of the block
//	Content: bytes forming the block of the file
type Block struct {
	Content []byte
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
