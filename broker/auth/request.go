package auth

import (
	"crypto/rsa"
	"encoding/json"
)

// EncryptedRequest contains the AES key the peer and the broker will share
// encrypted with the peer's RSA key and the data the peer will need for
// authentication encrypted with the AES key
type EncryptedRequest struct {
	Key  string `json:"string"`
	Data string `json:"string"`
}

// MakeEncryptedRequest encrypts the AES key with the RSA key and the data
// with the AES key
func MakeEncryptedRequest(RSAkey *rsa.PublicKey, r *Request) EncryptedRequest {
	eR := new(EncryptedRequest)
	eR.Key = string(RSAEncrypt(RSAkey, r.AESKey))
	eR.Data = string(AESEncrypt(r.AESKey, []byte(r.String())))
	return *eR
}

// Request is used to send the problem statement and public key of the
// broker session
type Request struct {
	AESKey  []byte
	Problem string `json:"problem"`
	Token   string `json:"token"`
}

// DecryptRequest decrypts an EncryptedRequest given the RSA private key
func DecryptRequest(eR EncryptedRequest, RSAkey *rsa.PrivateKey) (*Request, error) {
	aesKey, err := RSADecrypt(RSAkey, []byte(eR.Key))
	if err != nil {
		return nil, err
	}
	data, err := AESDecrypt(aesKey, []byte(eR.Data))
	if err != nil {
		return nil, err
	}
	req := new(Request)
	err = json.Unmarshal(data, req)
	if err != nil {
		return nil, err
	}
	req.AESKey = aesKey
	return req, nil
}

// MakeRequest is the Request constructor
func MakeRequest(RSAkey *rsa.PublicKey, AESkey []byte, problem Problem, token string) Request {
	return Request{
		AESKey:  AESkey,
		Problem: problem.Formulate(),
		Token:   token,
	}
}

// String returns a JSON string of the request
func (r *Request) String() string {
	b, _ := json.Marshal(r)
	return string(b)
}
