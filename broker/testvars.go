package broker

// this file contains all the global variables used *only* for testing

import (
	"fmt"
	"time"
)

const (
	// max amount of connection attempts to broker
	maxConnTries = 1e2
	// peer id
	peerID = "p1"
	// public RSA keys, faster than generating individual keys for each test
	//                                                                        ↓ changed char
	pubKey    = "MIGJAoGBAJJYXgBem1scLKPEjwKrW8+ci3B/YNN3aY2DJ3lc5e2wNc0SmFikDaow1TdYcKl2wdrXX7sMRsyjTk15IECMezyHzaJGQ9TinnkQixJ+YnlNdLC04TNWOg13plyahIXBforYAjYl2wVIA8Yma2bEQFhmAFkEX1A/Q1dIKy6EfQ+xAgMBAAE="
	pubKeyAlt = "MIGJAoGBAJJYXgBem1scLKPEjwKrW8+ci3B/YNN3aY2DJ3lc5e2wNc0SmFikDbow1TdYcKl2wdrXX7sMRsyjTk15IECMezyHzaJGQ9TinnkQixJ+YnlNdLC04TNWOg13plyahIXBforYAjYl2wVIA8Yma2bEQFhmAFkEX1A/Q1dIKy6EfQ+xAgMBAAE="
	// multiaddr used for testing
	testmultiaddr = "/ip4/127.0.0.1/tcp/3081"
)

var connSleepTime = time.Duration(5e5) * time.Nanosecond                 // time between atempts
var bURL = fmt.Sprintf("http://%s:%d", HTTPDefaultAddr, HTTPDefaultPort) // url of broker
var pURL = fmt.Sprintf("%s/peer", bURL)                                  // bURL + path to peer api
