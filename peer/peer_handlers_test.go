package peer

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"bitbucket.org/mikelsr/sakaban-broker/auth"
	"bitbucket.org/mikelsr/sakaban/fs"
	"bitbucket.org/mikelsr/sakaban/peer/comm"
	uuid "github.com/satori/go.uuid"
)

func TestPeer_HandleRequestMTBlockContent(t *testing.T) {
	fileName := filepath.Join(testDir, "testfile")
	fileID, _ := uuid.NewV4()
	testIntPeer1.stack = *newFileStack()
	testIntPeer1.stack.push(&fs.Summary{
		ID:     fileID.String(),
		Parent: "",
		Path:   fileName,
		Blocks: []uint64{1, 1},
	}, &testIntPeer1.Contacts[0])
	// push and iter nil file to generate tmpFile for first summary
	testIntPeer1.stack.push(nil, nil)
	testIntPeer1.stack.iterFile()

	content1 := make([]byte, fs.BlockSize)
	content2 := make([]byte, fs.BlockSize)
	rand.Read(content1)
	rand.Read(content2)

	bc1 := comm.BlockContent{
		BlockN:    0,
		BlockSize: uint16(fs.BlockSize / 1024),
		Content:   content1,
		FileID:    fileID,
	}
	bc2 := comm.BlockContent{
		BlockN:    1,
		BlockSize: uint16(fs.BlockSize / 1024),
		Content:   content2,
		FileID:    fileID,
	}

	// connection to send first block
	dump := bc1.Dump()

	log.Println("[Test]\tWaiting for Peer 1 to receive first block...")
	for testIntPeer1.stack.tmpFile.Blocks[0] == nil {
		// test will timeout if content isn't stored by testIntPeer1
		s, err := testIntPeer2.ConnectTo(testIntPeer2.Contacts[0 /* testIntPeer1 */])
		if err != nil {
			t.FailNow()
		}
		s.Write(dump)
		time.Sleep(time.Millisecond * 100)
		s.Close()
	}

	// ensure that blocks are equal
	if !bytes.Equal(testIntPeer1.stack.tmpFile.Blocks[0].Content, bc1.Content) {
		t.FailNow()
	}

	// connection to send second block
	dump = bc2.Dump()
	log.Println("[Test]\tWaiting for Peer 1 to receive second block...")
	for testIntPeer1.stack.tmpFile.Blocks[1] == nil {
		// test will timeout if content isn't stored by testIntPeer1
		s, err := testIntPeer2.ConnectTo(testIntPeer2.Contacts[0 /* testIntPeer1 */])
		if err != nil {
			t.FailNow()
		}
		s.Write(dump)
		time.Sleep(time.Millisecond * 100)
		s.Close()
	}
	log.Println("[Test]\tWaiting for Peer 1 to write file...")
	for {
		if _, err := os.Stat(fileName); !os.IsNotExist(err) {
			break
		}
	}
	// compare written file to sent file
	f, err := fs.MakeFile(fileName)
	if err != nil {
		t.FailNow()
	}
	if f.Blocks, err = f.Slice(); err != nil {
		fmt.Println(err)
		t.FailNow()
	}
	if !bytes.Equal(f.Blocks[0].Content, bc1.Content) ||
		!bytes.Equal(f.Blocks[1].Content, bc2.Content) {
		t.FailNow()
	}
}

func TestPeer_HandleRequestMTBlockRequest(t *testing.T) {
	s, err := testIntPeer2.ConnectTo(testIntPeer2.Contacts[0 /* testIntPeer1 */])
	if err != nil {
		t.FailNow()
	}
	summary, found := testIntPeer1.RootIndex.Files[muffinPath]
	if !found {
		t.FailNow()
	}
	absPath := summary.Path
	relPath := strings.Replace(absPath, testIntPeer2.RootDir+"/", "", 1)
	id, _ := uuid.FromString(summary.ID)

	blockN := uint8(1)
	br := comm.BlockRequest{
		BlockN:   blockN,
		FileID:   id,
		FilePath: relPath,
	}
	payload := br.Dump()
	n, err := s.Write(payload)
	if err != nil || n != len(payload) {
		t.FailNow()
	}

	buf := bufio.NewReader(s)
	bc := comm.BlockContent{}
	msg, err := bc.Recv(buf)
	if err != nil {
		t.FailNow()
	}
	if err = bc.Load(msg); err != nil {
		t.FailNow()
	}

	f, _ := fs.MakeFile(absPath)
	if !bytes.Equal(f.Blocks[blockN].Content, bc.Content) {
		t.FailNow()
	}
}

func TestPeer_HandleRequestMTIndexContent(t *testing.T) {
	c, _ := testIntPeer4.RequestPeer(auth.PrintPubKey(testIntPeer3.PubKey))
	testIntPeer4.Contacts = []Contact{*c}
	testIntPeer4.waiting = true
	ic := comm.IndexContent{Index: testIntPeer3.RootIndex}
	dump := ic.Dump()
	// connect to peer 4 and send index until success or timeout
	for testIntPeer4.stack.tmpFile == nil {
		s, err := testIntPeer3.ConnectTo(testIntPeer3.Contacts[0 /* testIntPeer4 */])
		if err != nil {
			t.FailNow()
		}
		log.Printf("[Test]\tSending '%s' to testIntPeer4...\n", string(dump))
		s.Write(dump)
		log.Println("[Test]\tWaiting for testIntPeer4 to receive index...")
		time.Sleep(time.Millisecond * 100)
		s.Close()
	}
}

func TestPeer_HandleRequestMTIndexRequest(t *testing.T) {
	s, err := testIntPeer2.ConnectTo(testIntPeer2.Contacts[0 /* testIntPeer1 */])
	if err != nil {
		t.FailNow()
	}
	ir := comm.IndexRequest{}
	payload := ir.Dump()
	n, err := s.Write(payload)
	if err != nil || n != len(payload) {
		t.FailNow()
	}
	buf := bufio.NewReader(s)
	ic := comm.IndexContent{}
	msg, err := ic.Recv(buf)
	if err != nil {
		t.FailNow()
	}
	if err = ic.Load(msg); err != nil {
		t.FailNow()
	}
	if !ic.Index.Equals(&testIntPeer1.RootIndex) {
		t.FailNow()
	}
}
