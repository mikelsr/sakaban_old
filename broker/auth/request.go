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
	Token   string `json:"token"`
}

// MakeRequest is the Request constructor
func MakeRequest(pubKey *rsa.PublicKey, problem Problem, token string) Request {
	return Request{
		Problem: problem.Formulate(),
		PubKey:  PrintPubKey(pubKey),
		Token:   token,
	}
}

// String returns a JSON string of the request
func (r *Request) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
