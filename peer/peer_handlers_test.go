package peer

import (
	"fmt"
	"strings"
	"testing"

	"bitbucket.org/mikelsr/sakaban/peer/comm"
	"github.com/satori/go.uuid"
)

func TestPeer_HandleRequestMTBlockRequest(t *testing.T) {
	s, err := testIntPeer1.ConnectTo(testIntPeer1.Contacts[0 /* testIntPeer2 */])
	if err != nil {
		t.FailNow()
	}
	summary, found := testIntPeer1.RootIndex.Files[muffinPath]
	if !found {
		t.FailNow()
	}
	absPath := summary.Path
	relPath := strings.Replace(absPath, testIntPeer1.RootDir, "", 1)
	id, _ := uuid.FromBytes([]byte(summary.ID))

	br := comm.BlockRequest{
		BlockN:   0,
		FileID:   id,
		FilePath: relPath,
	}
	payload := br.Dump()
	if n, err := s.Write(payload); err != nil || n != len(payload) {
		t.FailNow()
	}
	response := make([]byte, 1024*1024*2)
	n, err := s.Read(response)
	fmt.Println(n)
	if err != nil {
		t.FailNow()
	}
}
