package auth

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"testing"
)

func TestDecryptRequest(t *testing.T) {
	problem := NewProblem()
	aesKey := AESNewKey()
	r1 := MakeRequest(aesKey, pubKey, problem, "testtoken")
	eR := MakeEncryptedRequest(pub, &r1)

	// invalid rsa private key
	wrongKey, _ := rsa.GenerateKey(rand.Reader, keySize)
	if _, err := DecryptRequest(eR, wrongKey); err == nil {
		t.FailNow()
	}

	// wrong AES key in EncryptedRequest
	eR1 := MakeEncryptedRequest(pub, &r1)
	eR1.Key = string(RSAEncrypt(pub, AESNewKey()))
	if _, err := DecryptRequest(eR1, prv); err == nil {
		t.FailNow()
	}

	// successful attempt
	r3, err := DecryptRequest(eR, prv)
	if err != nil {
		t.FailNow()
	}
	if !bytes.Equal(r1.AESKey, r3.AESKey) || r1.Problem != r3.Problem || r1.Token != r3.Token {
		t.FailNow()
	}
}

func TestEncryptedRequest_String(t *testing.T) {
	problem := NewProblem()
	aesKey := AESNewKey()
	r := MakeRequest(aesKey, pubKey, problem, "testtoken")
	eR1 := MakeEncryptedRequest(pub, &r)
	str := eR1.String()
	eR2 := new(EncryptedRequest)
	err := json.Unmarshal([]byte(str), eR2)
	if err != nil {
		t.FailNow()
	}
}

func TestRequest_String(t *testing.T) {
	problem := NewProblem()
	aesKey := AESNewKey()
	r1 := MakeRequest(aesKey, pubKey, problem, "testtoken")
	str := r1.String()
	r2 := new(Request)
	err := json.Unmarshal([]byte(str), r2)
	if err != nil {
		t.FailNow()
	}
}
