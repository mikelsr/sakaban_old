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

func TestAESDecrypt(t *testing.T) {
	aesKey := AESNewKey()
	encrypted := []byte{42} // smaller than aes.BlockSize
	if _, err := AESDecrypt(aesKey, encrypted); err == nil {
		t.FailNow()
	}
}

func TestAESEncrypt(t *testing.T) {
	aesKey := AESNewKey()
	data := []byte{42}
	decrypted, err := AESDecrypt(aesKey, AESEncrypt(aesKey, data))
	if err != nil || !bytes.Equal(decrypted, data) {
		t.FailNow()
	}
}

func TestRSADecrypt(t *testing.T) {
	wrongKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	encrypted := RSAEncrypt(pub, testData)

	// decrypt with wrong key
	_, err := RSADecrypt(wrongKey, encrypted)
	if err == nil {
		t.FailNow()
	}

	// decrypt with correct key
	decrypted, err := RSADecrypt(prv, encrypted)
	if err != nil || !bytes.Equal(testData, decrypted) {
		t.FailNow()
	}
}

func TestRSAEncrypt(t *testing.T) {
	encrypted := RSAEncrypt(pub, testData)
	decrypted, _ := RSADecrypt(prv, encrypted)
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
