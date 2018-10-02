package peer

import (
	"fmt"
	"math/rand"

	"bitbucket.org/mikelsr/sakaban-broker/broker"
	"bitbucket.org/mikelsr/sakaban/fs"
)

const (
	testBrokerIP       = "127.0.0.1"
	testBrokerPort     = 3080
	testErrDialBackOff = "dial backoff" // hopefully fixed
)

var (
	muffinPath           = fmt.Sprintf("%s/res/muffin.jpg", fs.ProjectPath())
	testBrokerAddr       = fmt.Sprintf("http://%s:%d", testBrokerIP, testBrokerPort)
	testBroker           broker.Broker
	testDir              = fmt.Sprintf("/tmp/sakaban-test-%d", rand.Intn(1e8))
	testFailDir          = testDir + "/fail"
	testIntPeer1         Peer // used for integration testing
	testIntPeer2         Peer // used for integration testing
	testPeer             Peer
	testPeerRootDir      = fmt.Sprintf("%s/res", fs.ProjectPath())
	testListenMultiAddr1 = "/ip4/0.0.0.0/tcp/3011"
	testListenMultiAddr2 = "/ip4/0.0.0.0/tcp/3012"
	testListenMultiAddr3 = "/ip4/0.0.0.0/tcp/3013"
)
