package auth

import "crypto/rsa"

const (
	// public RSA keys, faster than generating individual keys for each test
	//                                                                        ↓ changed char
	pubKey    = "MIGJAoGBAJJYXgBem1scLKPEjwKrW8+ci3B/YNN3aY2DJ3lc5e2wNc0SmFikDaow1TdYcKl2wdrXX7sMRsyjTk15IECMezyHzaJGQ9TinnkQixJ+YnlNdLC04TNWOg13plyahIXBforYAjYl2wVIA8Yma2bEQFhmAFkEX1A/Q1dIKy6EfQ+xAgMBAAE="
	pubKeyAlt = "MIGJAoGBAJJYXgBem1scLKPEjwKrW8+ci3B/YNN3aY2DJ3lc5e2wNc0SmFikDbow1TdYcKl2wdrXX7sMRsyjTk15IECMezyHzaJGQ9TinnkQixJ+YnlNdLC04TNWOg13plyahIXBforYAjYl2wVIA8Yma2bEQFhmAFkEX1A/Q1dIKy6EfQ+xAgMBAAE="
)

var prv *rsa.PrivateKey
var pub *rsa.PublicKey
var testData = []byte{42}
