package peer

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"bitbucket.org/mikelsr/sakaban-broker/auth"
	"bitbucket.org/mikelsr/sakaban-broker/broker"
	libp2p "github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	multiaddr "github.com/multiformats/go-multiaddr"
)

// Peer represents an individual device
type Peer struct {
	BrokerIP   string    // IPv4 address of the host
	BrokerPort int       // TCP port of the host
	Contacts   []Contact `json:"contacts"` // List of trusted contacts
	Host       host.Host `json:"-"`        // Host is the libp2p host
	// PrvKey and PubKey are used to verify the identity of the Peer
	PrvKey *rsa.PrivateKey `json:"-"`
	PubKey *rsa.PublicKey  `json:"-"`
}

// BrokerAddr returns the formatted address of the broker assigned to the peer
func (p *Peer) BrokerAddr() string {
	return fmt.Sprintf("%s:%d", p.BrokerIP, p.BrokerPort)
}

// CleanUp removes all exportable elements of a peer from a 'dir' directory
func CleanUp(dir string) {
	os.Remove(filepath.Join(dir, filenamePeer))
	os.Remove(filepath.Join(dir, filenamePrv))
	os.Remove(filepath.Join(dir, filenamePub))
}

// Export marshals the Peer struct and its keys into files located in 'dir'
func (p *Peer) Export(dir string) error {
	err := os.MkdirAll(dir, permissionDir)
	if err != nil {
		CleanUp(dir)
		return err
	}

	b, _ := json.Marshal(p)
	err = ioutil.WriteFile(filepath.Join(dir, filenamePeer), b, permissionFile)
	if err != nil {
		CleanUp(dir)
		return err
	}

	err = ExportRSAKeys(dir, p.PrvKey, p.PubKey)
	if err != nil {
		CleanUp(dir)
		return err
	}

	return err
}

// Import unmarshals a Peer from a directory containing the struct and keys
func Import(dir string) (*Peer, error) {
	b, err := ioutil.ReadFile(filepath.Join(dir, filenamePeer))
	if err != nil {
		return nil, err
	}
	prv, pub, err := ImportRSAKeys(dir)
	if err != nil {
		return nil, err
	}
	p := new(Peer)
	err = json.Unmarshal(b, p)
	if err != nil {
		return nil, err
	}
	p.PrvKey, p.PubKey = prv, pub
	return p, nil
}

// NewPeer creates a peer with a NEW PAIR OF KEYS
// for creating a peer with an existing pair of keys, use MakePeer
func NewPeer() (*Peer, error) {
	// generate RSA private key to communicate with broker
	prv, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		return nil, err
	}
	prv.Precompute()
	err = prv.Validate()
	if err != nil {
		return nil, err
	}

	// generate p2p RSA private key
	p2prv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA,
		keySize, rand.Reader)
	if err != nil {
		return nil, err
	}

	// p2p multiaddr
	addr, _ := multiaddr.NewMultiaddr(listenMultiAddr)

	// libp2p options
	options := []libp2p.Option{
		libp2p.ListenAddrs(addr), // listeing multiaddr
		libp2p.Identity(p2prv),   // private key
	}

	// create libp2p host
	h, err := libp2p.New(context.Background(), options...)
	if err != nil {
		return nil, err
	}

	// create peer
	return &Peer{
		BrokerIP:   brokerIP,
		BrokerPort: brokerPort,
		Host:       h,
		PrvKey:     prv,
		PubKey:     &prv.PublicKey,
	}, nil
}

// Register updates info about peer 'p' at the Broker
func (p *Peer) Register() error {
	// create client
	c := broker.Client{
		PeerID:    p.Host.ID().String(),
		MultiAddr: p.Host.Addrs()[0].String(),
	}
	// marshal client
	body, err := json.Marshal(c)
	if err != nil {
		return err
	}
	// send post request to /peer, obtain problem
	r, err := http.Post(fmt.Sprintf("http://%s/peer?publicKey=%s",
		p.BrokerAddr(), auth.PrintPubKey(p.PubKey)),
		"application/json", bytes.NewReader(body))
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("peer: %s", err)
	}
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	// receive encrypted problem
	eR := new(auth.EncryptedRequest)
	err = json.Unmarshal(body, eR)
	if err != nil {
		return err
	}
	req, err := auth.DecryptRequest(*eR, p.PrvKey)
	if err != nil {
		return err
	}
	// extract compontents from response
	brokerKey, err := auth.ExtractPubKey(req.BrokerRSAKey)
	if err != nil {
		return err
	}
	problem, err := auth.MakeProblemFromString(req.Problem)
	if err != nil {
		return err
	}
	// solve problem, convert to bytes
	solution := []byte(fmt.Sprint(problem.Solution()))
	// send post request to /auth, receive confirmation
	r, err = http.Post(fmt.Sprintf("http://%s/auth?token=%s",
		p.BrokerAddr(), req.Token), "application/json",
		bytes.NewReader(auth.RSAEncrypt(brokerKey, solution)))
	if err != nil {
		return err
	}
	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("auth: %s", err)
	}
	// no errors
	return nil
}

// RequestPeer obtains info about a peer from a broker given the public key
// of the peer
func (p *Peer) RequestPeer(publicKey string) (*Contact, error) {
	r, err := http.Get(fmt.Sprintf("http://%s/peer?publicKey=%s",
		p.BrokerAddr(), url.QueryEscape(publicKey)))
	if err != nil {
		return nil, err
	}
	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("couldn't retrieve peer '%s': '%s'",
			publicKey, r.Status)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	c := new(Contact)
	err = json.Unmarshal(body, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
