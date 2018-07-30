package peer

import (
	"crypto/rand"
	"crypto/rsa"
)

// Peer represents an individual device
type Peer struct {
	// PrvKey and PubKey are used to verify the identity of the Peer
	PrvKey *rsa.PrivateKey
	PubKey *rsa.PublicKey
}

// NewPeer creates a peer with a NEW PAIR OF KEYS
// for creating a peer with an existing pair of keys, use MakePeer
func NewPeer() (*Peer, error) {
	prv, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}
	prv.Precompute()
	err = prv.Validate()
	if err != nil {
		return nil, err
	}
	p := &Peer{PrvKey: prv, PubKey: &prv.PublicKey}
	return p, nil
}
