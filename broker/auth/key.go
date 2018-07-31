package auth

import (
	"crypto/rsa"
	"encoding/asn1"
	"encoding/base64"
)

// ExtractPubKey extracts a RSA public key from a string base64-ecoded string
func ExtractPubKey(key string) (*rsa.PublicKey, error) {
	b, err := base64.StdEncoding.DecodeString(key)
	if err != nil {
		return nil, err
	}
	var pub rsa.PublicKey
	_, err = asn1.Unmarshal(b, &pub)
	if err != nil {
		return nil, err
	}
	return &pub, nil
}

// PrintPubKey marshals and base64-encodes a RSA public key
func PrintPubKey(key *rsa.PublicKey) string {
	b, _ := asn1.Marshal(*key)
	return base64.StdEncoding.EncodeToString(b)
}
