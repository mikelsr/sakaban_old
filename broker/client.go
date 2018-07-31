package broker

import (
	"errors"

	multiaddr "github.com/multiformats/go-multiaddr"
)

// Client stores the ID and MultiAddr of a peer
type Client struct {
	PeerID    string `json:"peer_id"`
	MultiAddr string `json:"multiaddr"`
	// TODO: timeout?
}

// Ok checks that attributes are set and multiaddr is valid
func (c Client) Ok() (bool, error) {
	if c.PeerID == "" {
		return false, errors.New("PeerID can't be empty")
	}
	if c.MultiAddr == "" {
		return false, errors.New("MultiAddr can't be empty")
	}
	_, err := multiaddr.NewMultiaddr(c.MultiAddr)
	if err != nil {
		return false, err
	}
	return true, nil
}
