package broker

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	mrand "math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/mikelsr/sakaban/broker/auth"
)

// Broker stores information about peers indexed by their public key
type Broker struct {
	auths map[token]session
	peers map[string]client
}

type token string

// NewBroker should be used as the only constructor for Broker
func NewBroker() *Broker {
	b := new(Broker)
	b.auths = make(map[token]session)
	b.peers = make(map[string]client)
	return b
}

// newToken is the token constructor
func newToken() token {
	mrand.Seed(time.Now().UnixNano())
	str := make([]byte, tokenSize)
	for i := 0; i < tokenSize; i++ {
		str[i] = tokenChars[mrand.Intn(len(tokenChars))]
	}
	return token(str)
}

// genToken creates a unique token for the 'b' broker
func (b *Broker) genToken() token {
	t := newToken()
	for {
		if _, found := b.auths[t]; !found {
			return t
		}
		t = newToken()
	}
}

func (b *Broker) handleAuth(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	t := r.Form.Get("token")
	log.Printf("[Broker]\tReceived auth request, token '%s'\n", t)
	if t == "" {
		sendStatus(w, http.StatusUnauthorized)
	}

	s, found := b.auths[token(t)]
	if !found {
		sendStatus(w, http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendStatus(w, http.StatusInternalServerError)
		return
	}
	solBytes, err := auth.RSADecrypt(s.prv, body)
	// incorrect solution
	if err != nil {
		sendStatus(w, http.StatusUnauthorized)
		return
	}

	sol, err := strconv.ParseInt(string(solBytes), 10, 64)
	// incorrect solution
	if err != nil || sol != s.prob.Solution() {
		sendStatus(w, http.StatusUnauthorized)
		return
	}

	// update peer
	b.peers[s.peerKey] = s.changes
	sendStatus(w, http.StatusOK)
}

// handlePeer makes basic comprobation and delegates the request to the
// corresponding method
func (b *Broker) handlePeer(w http.ResponseWriter, r *http.Request) {
	log.Printf("[Broker]\t(%s)\t%s\n", r.Method, r.URL)
	r.ParseForm()
	// verify that a public key is provided
	keyStr := r.Form.Get(restPublicKey)
	if keyStr == "" {
		sendStatus(w, http.StatusBadRequest)
		return
	}

	keyStr, err := url.PathUnescape(keyStr)
	// FIXME why are '+' sings being treated literally and escaped as spaces
	keyStr = strings.Replace(keyStr, " ", "+", len(keyStr))
	if err != nil {
		sendStatus(w, http.StatusNotAcceptable)
		return
	}
	key, err := auth.ExtractPubKey(keyStr)
	if err != nil {
		sendStatus(w, http.StatusNotAcceptable)
		return
	}

	// delegate request
	if r.Method == http.MethodGet {
		b.handlePeerGET(w, r, keyStr)
	} else if r.Method == http.MethodPost {
		b.handlePeerPOST(w, r, keyStr, key)
	} else {
		sendStatus(w, http.StatusMethodNotAllowed)
	}
}

// handlePeerGET looks for the public key in Broker.peers and retuns it
// marshalled if found
func (b *Broker) handlePeerGET(w http.ResponseWriter, r *http.Request, pubKeyStr string) {
	if c, found := b.peers[pubKeyStr]; found {
		response, _ := json.Marshal(c)
		sendStatus(w, http.StatusOK)
		w.Write(response)
	} else {
		sendStatus(w, http.StatusNotFound)
	}
}

// handlePeerPOST registers/updates a client in Broker.peers
func (b *Broker) handlePeerPOST(w http.ResponseWriter, r *http.Request, pubKeyStr string, pubKey *rsa.PublicKey) {
	// unmarshall client (peer) from body
	var c client
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&c)
	if err != nil {
		sendStatus(w, http.StatusBadRequest)
		return
	}
	// verify integrity of client struct
	if ok, er := c.ok(); !ok {
		sendStatus(w, http.StatusBadRequest)
		w.Write([]byte(fmt.Sprint(er)))
		return
	}

	// create session for peer
	t := b.genToken()
	s := makeSession(pubKeyStr, c)
	b.auths[t] = s
	// creathe authentication request and encrypt it
	authReq := auth.MakeRequest(s.aesKey, auth.PrintPubKey(s.pub), s.prob, string(t))
	encryptedAuthReq := auth.MakeEncryptedRequest(pubKey, &authReq)
	// send encrypted authentication request
	sendStatus(w, http.StatusOK)
	w.Write([]byte(encryptedAuthReq.String()))
}

// ListenAndServe runs the http listener on the specified addr:port
func (b *Broker) ListenAndServe(addr string, port int) error {
	http.HandleFunc("/auth", b.handleAuth)
	http.HandleFunc("/peer", b.handlePeer)
	log.Printf("[Broker]\tListening at %s:%d", addr, port)
	err := http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), nil)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// ListenAndServeDefault calls ListenAndServe at HTTPDefaultAddr:HTTPDefaultPort
func (b *Broker) ListenAndServeDefault() error {
	return b.ListenAndServe(HTTPDefaultAddr, HTTPDefaultPort)
}

// sendStatus writes the status header to the ResponseWriter and logs it
func sendStatus(w http.ResponseWriter, s int) {
	w.WriteHeader(s)
	log.Printf("[Broker]\t%d\t%s", s, http.StatusText(s))
}
