package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"encoding/asn1"
	"encoding/base64"
	"errors"
	"io"
)

// TODO: rand

// AESDecrypt is used to encrypt messages given an aes private key
func AESDecrypt(key []byte, data []byte) ([]byte, error) {
	ciphertext, _ := base64.StdEncoding.DecodeString(string(data))
	block, _ := aes.NewCipher(key)
	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("decoding data is too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}

// AESEncrypt is used to encrypt messages given an aes public key
func AESEncrypt(key []byte, data []byte) []byte {
	block, _ := aes.NewCipher(key)
	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	io.ReadFull(rand.Reader, iv)
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)
	return []byte(base64.StdEncoding.EncodeToString(ciphertext))
}

// AESNewKey generates a new AES key
func AESNewKey() []byte {
	key := make([]byte, 32)
	rand.Read(key)
	return key
}

// RSADecrypt is used to encrypt messages given an rsa private key
func RSADecrypt(key *rsa.PrivateKey, data []byte) ([]byte, error) {
	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, key, data)
	if err != nil {
		return nil, err
	}
	return decrypted, nil
}

// RSAEncrypt is used to encrypt messages given an rsa public key
func RSAEncrypt(key *rsa.PublicKey, data []byte) []byte {
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
