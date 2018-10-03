package peer

import (
	"strings"
	"testing"

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

	br := comm.BlockRequest{
		BlockN:   1,
		FileID:   id,
		FilePath: relPath,
	}
	payload := append(br.Dump(), 0xFF)
	n, err := s.Write(payload)
	if err != nil || n != len(payload) {
		t.FailNow()
	}
	response := make([]byte, 1024*1024*2)
	_, err = s.Read(response)
	if err != nil {
		t.FailNow()
	}
}
