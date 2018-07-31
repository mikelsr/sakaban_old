package broker

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"bitbucket.org/mikelsr/sakaban/broker/auth"
)

var testBroker *Broker

func TestMain(m *testing.M) {
	testBroker = NewBroker()
	go testBroker.ListenAndServeDefault()

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

func TestBroker_genToken(t *testing.T) {
	b := NewBroker()
	t1 := newToken()
	b.Auths[t1] = *new(session)
	t2 := b.genToken()
	if t1 == t2 {
		t.FailNow()
	}
}

func TestBroker_handleAuth(t *testing.T) {
	_prv, _ := rsa.GenerateKey(rand.Reader, keySize)
	_pub := &_prv.PublicKey

	// no token
	r, err := http.Post(fmt.Sprintf("%s/auth", bURL), "text/plain", nil)
	if err != nil || r.StatusCode != http.StatusUnauthorized {
		t.FailNow()
	}

	// create auth request
	r, err = postTestPeer(auth.PrintPubKey(_pub), "pid1", "/ip4/0.0.0.0/tcp/3001")
	if err != nil || r.StatusCode != http.StatusOK {
		t.FailNow()
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		t.FailNow()
	}
	eR := new(auth.EncryptedRequest)
	err = json.Unmarshal(body, eR)
	if err != nil {
		t.FailNow()
	}
	req, err := auth.DecryptRequest(*eR, _prv)
	if err != nil {
		t.FailNow()
	}

	// incorrectly encrypted solution
	r, err = http.Post(fmt.Sprintf("%s/auth?token=%s", bURL, req.Token),
		"text/plain", nil)
	if err != nil || r.StatusCode != http.StatusUnauthorized {
		t.FailNow()
	}

	// correctly encrypted wrong solution
	brokerRSAKey, err := auth.ExtractPubKey(req.BrokerRSAKey)
	if err != nil {
		t.FailNow()
	}
	r, err = http.Post(fmt.Sprintf("%s/auth?token=%s", bURL, req.Token),
		"text/plain",
		bytes.NewReader(auth.RSAEncrypt(brokerRSAKey, []byte(fmt.Sprint(42)))))
	if err != nil || r.StatusCode != http.StatusUnauthorized {
		t.FailNow()
	}

	// correct encrypted correct solution
	problem, err := auth.MakeProblemFromString(req.Problem)
	if err != nil {
		t.FailNow()
	}
	r, err = http.Post(fmt.Sprintf("%s/auth?token=%s", bURL, req.Token),
		"text/plain",
		bytes.NewReader(auth.RSAEncrypt(brokerRSAKey, []byte(fmt.Sprint(problem.Solution())))))
	if err != nil || r.StatusCode != http.StatusOK {
		t.FailNow()
	}

	r, err = getTestPeer(auth.PrintPubKey(_pub))
	if err != nil {
		t.FailNow()
	}
	body, err = ioutil.ReadAll(r.Body)
	if err != nil {
		t.FailNow()
	}
	// check that the peer was correctly registered
	if !strings.Contains(string(body), "pid1") {
		t.FailNow()
	}
}

func TestBroker_handlePeer(t *testing.T) {
	// no public key param
	if r, err := http.Get(pURL); err != nil || r.StatusCode != http.StatusBadRequest {
		t.FailNow()
	}

	// invalid public key
	if r, _ := getTestPeer("_"); r.StatusCode != http.StatusNotAcceptable {
		t.FailNow()
	}

	// valid public key, invalid method (PUT)
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s?%s=%s", pURL, restPublicKey, url.QueryEscape(pubKey)), bytes.NewReader([]byte{}))
	if err != nil {
		t.FailNow()
	}
	if r, err := http.DefaultClient.Do(req); err != nil || r.StatusCode != http.StatusMethodNotAllowed {
		t.FailNow()
	}
}

func TestBroker_handlePeerGET(t *testing.T) {
	// non-existing peer
	if r, _ := getTestPeer(pubKeyAlt); r.StatusCode != http.StatusNotFound {
		t.FailNow()
	}
	// register test peer
	testBroker.Peers[pubKeyAlt] = Client{
		PeerID:    "g1",
		MultiAddr: "/ip4/0.0.0.0/tcp/3001",
	}
	// retrieve test peer
	if r, _ := getTestPeer(pubKeyAlt); r.StatusCode != http.StatusOK {
		t.FailNow()
	}
}

func TestBroker_handlePeerPOST(t *testing.T) {
	// fail cases
	if r, _ := postTestPeer(pubKey, "", ""); r.StatusCode != http.StatusBadRequest {
		t.FailNow()
	}
	if r, _ := postTestPeer(pubKey, "p1", ""); r.StatusCode != http.StatusBadRequest {
		t.FailNow()
	}
	if r, _ := postTestPeer(pubKey, "p1", "_"); r.StatusCode != http.StatusBadRequest {
		t.FailNow()
	}
	// POST invalid JSON
	cli := http.DefaultClient
	r, _ := cli.Post(fmt.Sprintf("%s?%s=%s", pURL, restPublicKey, pubKey),
		"application/json", bytes.NewReader([]byte("{")))
	if r.StatusCode != http.StatusBadRequest {
		t.FailNow()
	}
	// peer should not be registered
	if _, found := testBroker.Peers[pubKey]; found {
		t.FailNow()
	}

	// success cases
	if r, _ = postTestPeer(pubKey, "p1", "/ip4/0.0.0.0/tcp/3001"); r.StatusCode != http.StatusOK {
		t.FailNow()
	}
	// ommit auth, manually register peer
	testBroker.Peers[pubKey] = Client{
		PeerID:    "p1",
		MultiAddr: "/ip4/0.0.0.0/tcp/3001",
	}
	r, _ = getTestPeer(pubKey)
	if r.StatusCode != http.StatusOK {
		t.FailNow()
	}

	decoder := json.NewDecoder(r.Body)
	c := new(Client)
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
	r, err := http.Get(fmt.Sprintf("%s?%s=%s", pURL, restPublicKey, url.QueryEscape(pk)))
	if err != nil {
		return nil, err
	}
	return r, nil
}

// getTestPeer posts a peer to a broker running in the default address
func postTestPeer(pk string, pid string, ma string) (*http.Response, error) {
	peer := new(Client)
	peer.PeerID = pid
	peer.MultiAddr = ma

	body, err := json.Marshal(peer)
	if err != nil {
		return nil, err
	}
	response, err := http.Post(fmt.Sprintf("%s?%s=%s",
		pURL, restPublicKey, url.QueryEscape(pk)),
		"application/json",
		bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	return response, nil
}
