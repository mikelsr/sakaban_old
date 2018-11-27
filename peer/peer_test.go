package peer

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"bitbucket.org/mikelsr/sakaban-broker/auth"
	"bitbucket.org/mikelsr/sakaban-broker/broker"
	"bitbucket.org/mikelsr/sakaban/fs"
	"bitbucket.org/mikelsr/sakaban/peer/comm"

	libp2p "github.com/libp2p/go-libp2p"
	net "github.com/libp2p/go-libp2p-net"
	p2peer "github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	multiaddr "github.com/multiformats/go-multiaddr"
)

func createTestPeers() (*Peer, *Peer, *Peer, *Peer) {
	p1, _ := NewPeer()
	p2, _ := NewPeer()
	p3, _ := NewPeer()
	p4, _ := NewPeer()

	// p1
	addr, _ := multiaddr.NewMultiaddr(testListenMultiAddr1)
	options := []libp2p.Option{
		libp2p.ListenAddrs(addr), // listeing multiaddr
	}
	h, _ := libp2p.New(context.Background(), options...)
	p1.Host = h
	p1.Host.SetStreamHandler(protocolID, p1.HandleStream)

	// p2
	addr, _ = multiaddr.NewMultiaddr(testListenMultiAddr2)
	options = []libp2p.Option{
		libp2p.ListenAddrs(addr), // listeing multiaddr
	}
	h, _ = libp2p.New(context.Background(), options...)
	p2.Host = h
	p2.Host.SetStreamHandler(protocolID, p2.HandleStream)

	// p3
	addr, _ = multiaddr.NewMultiaddr(testListenMultiAddr3)
	options = []libp2p.Option{
		libp2p.ListenAddrs(addr), // listeing multiaddr
	}
	h, _ = libp2p.New(context.Background(), options...)
	p3.Host = h
	p3.Host.SetStreamHandler(protocolID, p3.HandleStream)

	// p4
	addr, _ = multiaddr.NewMultiaddr(testListenMultiAddr4)
	options = []libp2p.Option{
		libp2p.ListenAddrs(addr), // listeing multiaddr
	}
	h, _ = libp2p.New(context.Background(), options...)
	p4.Host = h
	p4.Host.SetStreamHandler(protocolID, p4.HandleStream)

	c1 := Contact{
		Addr:      testListenMultiAddr1,
		PeerID:    p1.Host.ID().Pretty(),
		RSAPubKEy: auth.PrintPubKey(p1.PubKey),
	}
	c2 := Contact{
		Addr:      testListenMultiAddr2,
		PeerID:    p2.Host.ID().Pretty(),
		RSAPubKEy: auth.PrintPubKey(p2.PubKey),
	}
	c4 := Contact{
		Addr:      testListenMultiAddr4,
		PeerID:    p4.Host.ID().Pretty(),
		RSAPubKEy: auth.PrintPubKey(p4.PubKey),
	}

	p1.Contacts = []Contact{c2}
	p2.Contacts = []Contact{c1}
	p3.Contacts = []Contact{c4}

	p1.Host.Peerstore().AddAddr(c2.ID(), c2.MultiAddr(), pstore.PermanentAddrTTL)
	p2.Host.Peerstore().AddAddr(c1.ID(), c1.MultiAddr(), pstore.PermanentAddrTTL)
	p3.Host.Peerstore().AddAddr(c4.ID(), c4.MultiAddr(), pstore.PermanentAddrTTL)

	p1.RootDir = testPeerRootDir
	p2.RootDir = testPeerRootDir
	p3.RootDir = testPeerRootDir
	p4.RootDir = testPeerRootDir

	p1.ReloadIndex()
	p3.ReloadIndex()
	i4, _ := fs.MakeIndex()
	p4.RootIndex = *i4

	p1.Register()
	p2.Register()
	p3.Register()
	p4.Register()
	return p1, p2, p3, p4
}

// testListenAndServe stops the tests whenever broker.ListenAndServe fails
func testListenAndServe(b *broker.Broker, ip string, port int) {
	if err := testBroker.ListenAndServe(testBrokerIP, testBrokerPort); err != nil {
		panic(err)
	}
}

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
	go testListenAndServe(&testBroker, testBrokerIP, testBrokerPort)

	testIntPeer1, testIntPeer2, testIntPeer3, testIntPeer4 = createTestPeers()

	// cleanup
	defer os.RemoveAll(testDir)

	// run tests
	m.Run()
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
	addr, _ := multiaddr.NewMultiaddr(testListenMultiAddrAux)
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
	err := errors.New("")
	var s net.Stream
	for err != nil {
		// create valid peer listening in unused addr
		p, _ := NewPeer()
		addr, _ := multiaddr.NewMultiaddr(testListenMultiAddrAux)
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
		s, err = testPeer.ConnectTo(*c)
		// close host to avoid conflict in the next iteration
		if err != nil {
			p.Host.Close()
		}
	}

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

func TestPeer_ReloadPeer(t *testing.T) {
	testPeer.RootDir = testPeerRootDir
	testPeer.ReloadIndex()
	if reflect.DeepEqual(testPeer.RootIndex, fs.Index{}) {
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
