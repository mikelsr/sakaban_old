package peer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"bitbucket.org/mikelsr/sakaban-broker/auth"
	"bitbucket.org/mikelsr/sakaban-broker/broker"
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
	os.Mkdir(filepath.Join(testDir, "peer"), 0755)
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

func TestExport(t *testing.T) {
	err := testPeer.Export(filepath.Join(testFailDir, "export"))
	if err == nil {
		t.FailNow()
	}

	// export to non-writeable directory
	err = testPeer.Export(testFailDir)
	if err == nil {
		t.FailNow()
	}

	// correct export
	dir := filepath.Join(testDir, "peer", "export")
	err = testPeer.Export(dir)
	if err != nil {
		t.FailNow()
	}
	files, _ := ioutil.ReadDir(dir)
	if len(files) != 3 {
		t.FailNow()
	}
}

func TestImport(t *testing.T) {
	// import from empty folder
	_, err := Import(testFailDir)
	if err == nil {
		t.FailNow()
	}

	// correct export
	dir := filepath.Join(testDir, "peer", "import")
	testPeer.Export(dir)
	// correct import
	_, err = Import(dir)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}

func TestPeer_RequestPeer(t *testing.T) {
	p, _ := NewPeer()
	port := p.BrokerPort
	p.BrokerPort = 0
	// invalid port
	_, err := p.RequestPeer("")
	if err == nil {
		t.FailNow()
	}
	p.BrokerPort = port
	// invalid key
	_, err = p.RequestPeer("")
	if err == nil {
		t.FailNow()
	}
	// post self
	p.Register()
	// successfully request self
	_, err = p.RequestPeer(auth.PrintPubKey(p.PubKey))
	if err != nil {
		t.FailNow()
	}
}

func TestPeer_Register(t *testing.T) {
	err := testPeer.Register()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}
}
