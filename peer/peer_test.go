package peer

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"bitbucket.org/mikelsr/sakaban/broker"
)

func TestMain(m *testing.M) {
	// create test peer with key pair
	tp, err := NewPeer()
	if err != nil {
		os.Exit(1)
	}
	testPeer = *tp

	// create test directories
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		os.Exit(1)
	}
	os.MkdirAll(testFailDir, 0000)

	// export necessary keys for testing
	ExportRSAKeys(testDir, testPeer.PrvKey, testPeer.PubKey)
	ExportRSAKeys(filepath.Join(testDir, "import"),
		testPeer.PrvKey, testPeer.PubKey)
	// create and run test broker
	testBroker = *broker.NewBroker()
	go testBroker.ListenAndServe(testBrokerIP, testBrokerPort)
	// run tests
	m.Run()
	// cleanup
	os.RemoveAll(testDir)
}

func TestPeer_UpdateInfo(t *testing.T) {
	p, _ := NewPeer()
	err := p.UpdateInfo()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}
