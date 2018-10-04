package peer

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"bitbucket.org/mikelsr/sakaban/fs"
	"bitbucket.org/mikelsr/sakaban/peer/comm"
	"github.com/satori/go.uuid"
)

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

	bc, err := recvBlockContent(s)
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	f, _ := fs.MakeFile(absPath)
	if !bytes.Equal(f.Blocks[blockN].Content, bc.Content) {
		t.FailNow()
	}
}
