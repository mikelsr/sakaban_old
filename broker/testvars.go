package broker

// this file contains all the global variables used *only* for testing

import (
	"fmt"
	"time"
)

const maxConnTries = 1e2                                                 // max amount of connection attempts to broker
var connSleepTime = time.Duration(5e5) * time.Nanosecond                 // time between atempts
var pubKey = "k1"                                                        // public key
var peerID = "p1"                                                        // peer id
var testmultiaddr = "/ip4/127.0.0.1/tcp/3081"                            // multiaddr used for testing
var bURL = fmt.Sprintf("http://%s:%d", httpDefaultAddr, httpDefaultPort) // url of broker
var pURL = fmt.Sprintf("%s/peer", bURL)                                  // bURL + path to peer api
