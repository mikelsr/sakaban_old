package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	prvKey, err := rsa.GenerateKey(rand.Reader, keySize)
	if err != nil {
		os.Exit(1)
	}
	prv = prvKey
	pub = &prv.PublicKey
	m.Run()
}

func TestDecrypt(t *testing.T) {
	wrongKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	encrypted := Encrypt(pub, testData)

	// decrypt with wrong key
	_, err := Decrypt(wrongKey, encrypted)
	if err == nil {
		t.FailNow()
	}

	// decrypt with correct key
	decrypted, err := Decrypt(prv, encrypted)
	if err != nil || !bytes.Equal(testData, decrypted) {
		t.FailNow()
	}
}

func TestEncrypt(t *testing.T) {
	encrypted := Encrypt(pub, testData)
	decrypted, _ := Decrypt(prv, encrypted)
	if !bytes.Equal(testData, decrypted) {
		t.FailNow()
	}
}

func TestExtractPubKey(t *testing.T) {
	// invalid keys
	_, err := ExtractPubKey(" ")
	if err == nil {
		t.FailNow()
	}
	_, err = ExtractPubKey("")
	if err == nil {
		t.FailNow()
	}
	// valid key
	_, err = ExtractPubKey(pubKey)
	if err != nil {
		t.FailNow()
	}
}

func TestPrintPubKey(t *testing.T) {
	key, _ := ExtractPubKey(pubKey)
	keyStr := PrintPubKey(key)
	if keyStr != pubKey {
		t.FailNow()
	}
}
