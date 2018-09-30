package peer

import (
	"fmt"
	"math/rand"

	"bitbucket.org/mikelsr/sakaban-broker/broker"
)

const (
	testBrokerIP       = "127.0.0.1"
	testBrokerPort     = 3080
	testErrDialBackOff = "dial backoff"
)

var (
	testBrokerAddr = fmt.Sprintf("http://%s:%d", testBrokerIP, testBrokerPort)
	testBroker     broker.Broker
	testDir        = fmt.Sprintf("sakaban-test-%d", rand.Intn(1e8))
	testFailDir    = testDir + "/fail"
	testPeer       Peer
	testIntPeer1   Peer // used for integration testing
	testIntPeer2   Peer // used for integration testing
)
