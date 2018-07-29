package broker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	b := NewBroker()
	go b.ListenAndServeDefault()

	// connection attempts to test broker
	connTries := 0

	_, err := http.Get(bURL)
	// wait for broker to start listening
	for err != nil && connTries < maxConnTries {
		time.Sleep(connSleepTime)
		connTries++
	}

	m.Run()
	//http.Get(fmt.Sprintf("%s/stop", bUrl))
}

func TestBroker_handlePeer(t *testing.T) {
	// no public key param
	if r, err := http.Get(pURL); err != nil || r.StatusCode != http.StatusBadRequest {
		t.FailNow()
	}
	// valid public key, invalid method (PUT)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s?%s=pk1", pURL, restPublicKey), bytes.NewReader([]byte{}))
	if err != nil {
		t.FailNow()
	}
	if r, err := http.DefaultClient.Do(req); err != nil || r.StatusCode != http.StatusMethodNotAllowed {
		t.FailNow()
	}
}

func TestBroker_handlePeerGET(t *testing.T) {
	// non-existing peer
	if r, _ := getTestPeer("_"); r.StatusCode != http.StatusNotFound {
		t.FailNow()
	}
	// register test peer
	postTestPeer("g1", "g1", "/ip4/0.0.0.0/tcp/3001")
	// retrieve test peer
	if r, _ := getTestPeer("g1"); r.StatusCode != http.StatusOK {
		t.FailNow()
	}
}

func TestBroker_handlePeerPOST(t *testing.T) {
	// fail cases
	if r, _ := postTestPeer("p1", "", ""); r.StatusCode != http.StatusBadRequest {
		t.FailNow()
	}
	if r, _ := postTestPeer("p1", "p1", ""); r.StatusCode != http.StatusBadRequest {
		t.FailNow()
	}
	if r, _ := postTestPeer("p1", "p1", "_"); r.StatusCode != http.StatusBadRequest {
		t.FailNow()
	}
	// POST invalid JSON
	cli := http.DefaultClient
	r, _ := cli.Post(fmt.Sprintf("%s?%s=%s", pURL, restPublicKey, "p1"),
		"application/json", bytes.NewReader([]byte("{")))
	if r.StatusCode != http.StatusBadRequest {
		t.FailNow()
	}
	// peer should not be registered
	if r, _ = getTestPeer("p1"); r.StatusCode != http.StatusNotFound {
		t.FailNow()
	}

	// success cases
	if r, _ = postTestPeer("p1", "p1", "/ip4/0.0.0.0/tcp/3001"); r.StatusCode != http.StatusOK {
		t.FailNow()
	}
	r, _ = getTestPeer("p1")
	if r.StatusCode != http.StatusOK {
		t.FailNow()
	}
	decoder := json.NewDecoder(r.Body)
	c := new(client)
	err := decoder.Decode(c)
	if err != nil {
		t.FailNow()
	}
	if c.PeerID != "p1" || c.MultiAddr != "/ip4/0.0.0.0/tcp/3001" {
		t.FailNow()
	}
}

// getTestPeer retrieves 'pk' peer from a broker running in the default address
func getTestPeer(pk string) (*http.Response, error) {
	r, err := http.Get(fmt.Sprintf("%s?%s=%s", pURL, restPublicKey, pk))
	if err != nil {
		return nil, err
	}
	return r, nil
}

// getTestPeer posts a peer to a broker running in the default address
func postTestPeer(pk string, pid string, ma string) (*http.Response, error) {
	peer := new(client)
	peer.PeerID = pid
	peer.MultiAddr = ma

	body, err := json.Marshal(peer)
	if err != nil {
		return nil, err
	}
	response, err := http.Post(fmt.Sprintf("%s?%s=%s",
		pURL, restPublicKey, pk),
		"application/json",
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return response, nil
}
