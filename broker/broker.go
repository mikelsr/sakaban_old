package broker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	multiaddr "github.com/multiformats/go-multiaddr"
)

// Broker stores information about peers indexed by their public key
type Broker struct {
	peers map[string]client
}

// client stores the ID and MultiAddr of a peer
type client struct {
	PeerID    string `json:"peer_id"`
	MultiAddr string `json:"multiaddr"`
	// TODO: timeout?
}

// NewBroker should be used as the only constructor for Broker
func NewBroker() *Broker {
	b := new(Broker)
	b.peers = make(map[string]client)
	return b
}

// handlePeer makes basic comprobation and delegates the request to the
// corresponding method
func (b *Broker) handlePeer(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	// verify that a public key is provided
	publicKey := r.Form.Get(restPublicKey)
	if publicKey == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// delegate request
	if r.Method == http.MethodGet {
		b.handlePeerGET(w, r)
	} else if r.Method == http.MethodPost {
		b.handlePeerPOST(w, r)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handlePeerGET looks for the public key in Broker.peers and retuns it
// marshalled if found
func (b *Broker) handlePeerGET(w http.ResponseWriter, r *http.Request) {
	if c, found := b.peers[r.Form.Get(restPublicKey)]; found {
		response, _ := json.Marshal(c)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

// handlePeerPOST registers/updates a client in Broker.peers
// TODO: verify that public key belongs to contacting peer
func (b *Broker) handlePeerPOST(w http.ResponseWriter, r *http.Request) {

	// TODO: check Content-Type?

	// Unmarshall client (peer) from body
	var c client
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if ok, err := verifyClient(c); !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(err)))
		return
	}
	b.peers[r.Form.Get(restPublicKey)] = c
	w.WriteHeader(http.StatusOK)
}

// ListenAndServe runs the http listener on the specified addr:port
func (b *Broker) ListenAndServe(addr string, port int) error {
	http.HandleFunc("/peer", b.handlePeer)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// ListenAndServeDefault calls ListenAndServe at httpDefaultAddr:httpDefaultPort
func (b *Broker) ListenAndServeDefault() {
	b.ListenAndServe(httpDefaultAddr, httpDefaultPort)
}

// verifyClient checks that attributes are set and multiaddr is valid
func verifyClient(c client) (bool, error) {
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
