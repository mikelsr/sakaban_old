package peer

import (
	"fmt"

	"bitbucket.org/mikelsr/sakaban/broker"
)

const (
	testBrokerIP   = "127.0.0.1"
	testBrokerPort = 3080
	testDir        = "/tmp/peertest"
	testFailDir    = "/tmp/peertest/fail"
)

var (
	testBrokerAddr = fmt.Sprintf("http://%s:%d", testBrokerIP, testBrokerPort)
	testBroker     broker.Broker
	testPeer       Peer
)
