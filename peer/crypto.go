package peer

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// MarshalRSAPrvKey marshals a RSA private key into a PEM file
func MarshalRSAPrvKey(prv *rsa.PrivateKey) []byte {
	return pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(prv),
		},
	)
}

// MarshalRSAPubKey marshals a RSA public key into a PEM file
func MarshalRSAPubKey(pub *rsa.PublicKey) ([]byte, error) {
	// pubBytes, err := x509.MarshalPKIXPublicKey(pub)
	// if err != nil {
	// 	return nil, err
	// }
	return pem.EncodeToMemory(
		&pem.Block{
			Type: "RSA PUBLIC KEY",
			// Bytes: pubBytes,
			Bytes: x509.MarshalPKCS1PublicKey(pub),
		},
	), nil
}

// UnmarshalRSAPrvKey unmarshals a private key from a PEM file
func UnmarshalRSAPrvKey(prvPEM []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(prvPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}
	prv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return prv, nil
}

// UnmarshalRSAPubKey unmarshals a public key from a PEM file
func UnmarshalRSAPubKey(pubPEM []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block")
	}
	pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pubKey, nil
}
