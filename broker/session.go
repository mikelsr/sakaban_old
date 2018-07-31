package broker

import (
	"crypto/rand"
	"crypto/rsa"

	"bitbucket.org/mikelsr/sakaban/broker/auth"
)

// session stores details for an authorization session
type session struct {
	aesKey  []byte
	changes client
	peerKey string
	prob    auth.Problem
	prv     *rsa.PrivateKey
	pub     *rsa.PublicKey
}

func makeSession(rsaKey string, changes client) session {
	prv, _ := rsa.GenerateKey(rand.Reader, keySize)
	aesKey := make([]byte, 32)
	rand.Read(aesKey)
	return session{
		aesKey:  aesKey,
		changes: changes,
		peerKey: rsaKey,
		prob:    auth.NewProblem(),
		prv:     prv,
		pub:     &prv.PublicKey,
	}
}
