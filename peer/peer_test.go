package peer

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"bitbucket.org/mikelsr/sakaban-broker/auth"
	"bitbucket.org/mikelsr/sakaban-broker/broker"
	"bitbucket.org/mikelsr/sakaban/peer/comm"

	libp2p "github.com/libp2p/go-libp2p"
	p2peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func TestMain(m *testing.M) {
	// create test peer with key pair
	tp, err := NewPeer()
	if err != nil {
		os.Exit(1)
	}
	testPeer = *tp

	tip1, _ := NewPeer()
	tip2, _ := NewPeer()
	testIntPeer1 = *tip1
	testIntPeer2 = *tip2

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

func TestPeer_ConnectTo(t *testing.T) {
	// create incorrect peer
	p, _ := NewPeer()
	wrongContact := Contact{
		Addr:   "/ip4/127.0.0.0/tcp/1",
		PeerID: p2peer.IDB58Encode(p.Host.ID()),
	}
	// add incorrect peer to peerstore
	testPeer.Host.Peerstore().AddAddr(wrongContact.ID(),
		wrongContact.MultiAddr(), pstore.PermanentAddrTTL)
	// connect to invalid peer
	_, err := testPeer.ConnectTo(wrongContact)
	if err == nil {
		t.FailNow()
	}

	// create valid peer listening in unused addr
	p, _ = NewPeer()
	addr, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/3002")
	options := []libp2p.Option{
		libp2p.ListenAddrs(addr), // listeing multiaddr
	}
	h, _ := libp2p.New(context.Background(), options...)
	p.Host = h
	// set stream handler of new peer
	p.Host.SetStreamHandler(protocolID, p.HandleStream)
	// register and retrieve peer at test broker
	p.Register()
	c, _ := p.RequestPeer(auth.PrintPubKey(p.PubKey))
	// add valid peer to peerstore
	testPeer.Host.Peerstore().AddAddr(c.ID(), c.MultiAddr(), pstore.PermanentAddrTTL)
	_, err = testPeer.ConnectTo(*c)
	if err != nil {
		t.FailNow()
	}
}

func TestPeer_Export(t *testing.T) {
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

func TestPeer_HandleStream(t *testing.T) {
	// create valid peer listening in unused addr
	p, _ := NewPeer()
	addr, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/3002")
	options := []libp2p.Option{
		libp2p.ListenAddrs(addr), // listeing multiaddr
	}
	h, _ := libp2p.New(context.Background(), options...)
	p.Host = h
	// set stream handler of new peer
	p.Host.SetStreamHandler(protocolID, p.HandleStream)
	// register and retrieve peer at test broker
	p.Register()
	c, _ := p.RequestPeer(auth.PrintPubKey(p.PubKey))
	// add valid peer to peerstore
	testPeer.Host.Peerstore().AddAddr(c.ID(), c.MultiAddr(), pstore.PermanentAddrTTL)

	// connect to peer
	// var s net.Stream
	// var err error
	// for {
	s, err := testPeer.ConnectTo(*c)
	if err != nil {
		// FIXME: why does this happen?
		if err.Error() == testErrDialBackOff {
			t.Log("Dial backoff error")
		}
		t.FailNow()
	}
	// 	continue
	// }

	// begin HandleStream test
	msg := append(comm.IndexRequest{}.Dump(), byte(0))
	if n, err := s.Write(msg); n != len(msg) || err != nil {
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
		t.FailNow()
	}
}

func TestPeer_SetRootDir(t *testing.T) {
	if err := testPeer.SetRootDir(""); err == nil {
		t.FailNow()
	}
	if err := testPeer.SetRootDir(testDir + filenamePub); err == nil {
		t.FailNow()
	}
	if err := testPeer.SetRootDir(testDir); err != nil {
		t.FailNow()
	}
}
