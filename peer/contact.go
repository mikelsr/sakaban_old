package peer

import (
	p2peer "github.com/libp2p/go-libp2p-peer"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// Contact stores information about external peers
type Contact struct {
	Addr      string `json:"multiaddr"`
	PeerID    string `json:"peer_id"`
	RSAPubKEy string `json:"rsa_public_key"`
}

// ID returns a libp2p-peer.ID from Contact.ID
func (c Contact) ID() p2peer.ID {
	id, _ := p2peer.IDB58Decode(c.PeerID)
	return id
}

// MultiAddr returns a MultiAddr struct from the Contact.Addr string
func (c Contact) MultiAddr() multiaddr.Multiaddr {
	ma, _ := multiaddr.NewMultiaddr(c.Addr)
	return ma
}
