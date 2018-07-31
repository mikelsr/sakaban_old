package auth

import (
	"crypto/rsa"
	"encoding/json"
)

// Request is used to send the problem statement and public key of the
// broker session
type Request struct {
	PubKey  string `json:"public_key"`
	Problem string `json:"problem"`
}

// MakeRequest is the Request constructor
func MakeRequest(pubKey *rsa.PublicKey, problem Problem) Request {
	return Request{
		Problem: problem.Formulate(),
		PubKey:  PrintPubKey(pubKey),
	}
}

// String returns a JSON string of the request
func (r *Request) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
