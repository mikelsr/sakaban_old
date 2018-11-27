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

	"bitbucket.org/mikelsr/sakaban/fs"
	"bitbucket.org/mikelsr/sakaban/peer/comm"
	"github.com/satori/go.uuid"
)

func TestPeer_HandleRequestMTBlockContent(t *testing.T) {
	fmt.Println("-- Starting BlockContent test")
	fileName := filepath.Join(testDir, "testfile")
	fileID, _ := uuid.NewV4()
	testIntPeer1.stack = *newFileStack()
	testIntPeer1.stack.push(&fs.Summary{
		ID:     fileID.String(),
		Parent: "",
		Path:   fileName,
		Blocks: []uint64{1, 1},
	})
	// push and iter nil file to generate tmpFile for first summary
	testIntPeer1.stack.push(nil)
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
	s, err := testIntPeer2.ConnectTo(testIntPeer2.Contacts[0 /* testIntPeer1 */])
	if err != nil {
		t.FailNow()
	}
	s.Write(bc1.Dump())
	s.Close()
	log.Println("[Test]\tWaiting for Peer 1 to receive first block...")
	for testIntPeer1.stack.tmpFile.Blocks[0] == nil {
		// test will timeout if content isn't stored by testIntPeer1
	}
	// ensure that blocks are equal
	if !bytes.Equal(testIntPeer1.stack.tmpFile.Blocks[0].Content, bc1.Content) {
		t.FailNow()
	}

	// connection to send second block
	s, err = testIntPeer2.ConnectTo(testIntPeer2.Contacts[0 /* testIntPeer1 */])
	if err != nil {
		t.FailNow()
	}
	s.Write(bc2.Dump())
	s.Close()
	log.Println("[Test]\tWaiting for Peer 1 to receive second block...")
	for testIntPeer1.stack.tmpFile.Blocks[1] == nil {
		// test will timeout if content isn't stored by testIntPeer1
	}
	log.Println("[Test]\tWaiting for Peer 1 to write file...")
	for {
		if _, err = os.Stat(fileName); !os.IsNotExist(err) {
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
	fmt.Println("-- Ending BlockContent test")
}

func TestPeer_HandleRequestMTBlockRequest(t *testing.T) {
	fmt.Println("-- Starting BlockRequest test")
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
	fmt.Println("-- Ending BlockRequest test")
}

func TestPeer_HandleRequestMTIndexContent(t *testing.T) {
	fmt.Println("-- Starting IndexContent test")
	testIntPeer3.waiting = true
	ic := comm.IndexContent{Index: testIntPeer1.RootIndex}
	s, err := testIntPeer2.ConnectTo(testIntPeer2.Contacts[1 /* testIntPeer3 */])
	if err != nil {
		t.FailNow()
	}
	s.Write(ic.Dump())
	log.Println("[Test]\tWaiting for testIntPeer3 to receive index...")
	for len(testIntPeer3.stack.files) != len(ic.Index.Files) {
		// timeout if index isn't loaded correctly by Peer 3
	}
	fmt.Println("-- Ending IndexContent test")
}

func TestPeer_HandleRequestMTIndexRequest(t *testing.T) {
	fmt.Println("-- Starting IndexRequest test")
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
	fmt.Println("-- Ending IndexRequest test")
}
