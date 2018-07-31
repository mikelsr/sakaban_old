package peer

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"bitbucket.org/mikelsr/sakaban/broker/auth"
)

// ExportRSAKeys calls ExportRSAPrvKey, ExportRSAPubKey to store 'prv' and 'pub' in 'dir'
func ExportRSAKeys(dir string, prv *rsa.PrivateKey, pub *rsa.PublicKey) error {
	err := ExportRSAPrvKey(dir, prv)
	if err != nil {
		return err
	}
	return ExportRSAPubKey(dir, pub)
}

// ExportRSAPrvKey writes a PEM containing the 'prv' RSA private key to a file in 'dir'
func ExportRSAPrvKey(dir string, prv *rsa.PrivateKey) error {
	// TODO: folder ownership/permissions
	os.MkdirAll(dir, permissionDir)
	prvBytes := MarshalRSAPrvKey(prv)
	err := ioutil.WriteFile(filepath.Join(dir, filenamePrv), prvBytes, permissionFile)
	return err
}

// ExportRSAPubKey writes a PEM containing the 'pub' RSA public key to a file in 'dir'
func ExportRSAPubKey(dir string, pub *rsa.PublicKey) error {
	// TODO: folder ownership/permissions
	os.MkdirAll(dir, permissionDir)
	pubBytes, err := MarshalRSAPubKey(pub)
	if err != nil {
		return nil
	}
	err = ioutil.WriteFile(filepath.Join(dir, filenamePub), pubBytes, permissionFile)
	return err
}

// ImportRSAKeys imports a private and a public key from 'dir'
func ImportRSAKeys(dir string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	prv, err := ImportRSAPrvKey(dir)
	if err != nil {
		return nil, nil, err
	}
	pub, err := ImportRSAPubKey(dir)
	if err != nil {
		return nil, nil, err
	}
	return prv, pub, nil
}

// ImportRSAPrvKey reads a private RSA key from 'dir'
func ImportRSAPrvKey(dir string) (*rsa.PrivateKey, error) {
	prvBytes, err := ioutil.ReadFile(filepath.Join(dir, filenamePrv))
	if err != nil {
		return nil, err
	}
	return UnmarshalRSAPrvKey(prvBytes)
}

// ImportRSAPubKey reads a public RSA key from 'dir'
func ImportRSAPubKey(dir string) (*rsa.PublicKey, error) {
	pubBytes, err := ioutil.ReadFile(filepath.Join(dir, filenamePub))
	if err != nil {
		return nil, err
	}
	return UnmarshalRSAPubKey(pubBytes)
}

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

// RSADecrypt calls auth.RSADecrypt with p.PrvKey
func (p *Peer) RSADecrypt(ciphertext []byte) []byte {
	data, _ := auth.RSADecrypt(p.PrvKey, ciphertext)
	return data
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
