package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/asn1"
	"encoding/base64"
)

// TODO: rand

// Decrypt is used to encrypt messages given an rsa private key
func Decrypt(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, key, data)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

// Encrypt is used to encrypt messages given an rsa public key
func Encrypt(key *rsa.PublicKey, data []byte) []byte {
	encrypted, _ := rsa.EncryptPKCS1v15(rand.Reader, key, data)
	return encrypted
}

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
