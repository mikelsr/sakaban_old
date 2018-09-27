package comm

import (
	"log"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
)

/* Block content */

// testBlockContent_Dump checks that the dumped slice has the expected length
func testBlockContentDump(t *testing.T, bc BlockContent) {
	if len(bc.Dump()) != 19+len(bc.content) {
		t.FailNow()
	}
}

// testBlockContent_Load loads a BlockContent from a bc.Dump() and compares it
// to the original (bc)
func testBlockContentLoad(t *testing.T, bc BlockContent) {
	bcLoaded := new(BlockContent)
	err := bcLoaded.Load(bc.Dump())
	if err != nil {
		log.Fatalln(err)
	}
	if !reflect.DeepEqual(*bcLoaded, bc) {
		t.FailNow()
	}

	/* error cases */
	if err = bc.Load([]byte{}); err == nil {
		t.FailNow()
	}
	if err = bc.Load([]byte{42}); err == nil {
		t.FailNow()
	}
	// Invalid block size
	bc.content = make([]byte, int(^uint8(0))*1024+1)
	if err = bc.Load(bc.Dump()); err == nil {
		t.FailNow()
	}
	// Mismatched block size
	bc.content = []byte{0}
	if err = bc.Load(bc.Dump()); err == nil {
		t.FailNow()
	}
}

func TestBlockContent(t *testing.T) {
	bc := *new(BlockContent)
	bc.blockN = 1
	bc.blockSize = 1
	bc.content = make([]byte, 1024)
	id, _ := uuid.NewV4()
	bc.fileID = id

	testBlockContentDump(t, bc)
	testBlockContentLoad(t, bc)
}

func TestBlockContent_Type(t *testing.T) {
	bc := *new(BlockContent)
	if bc.Type() != MTBlockContent {
		t.FailNow()
	}
}

/* Block request */

func testBlockRequestDump(t *testing.T, br BlockRequest) {
	if len(br.Dump()) != 18 {
		t.FailNow()
	}
}

func testBlockRequestLoad(t *testing.T, br BlockRequest) {
	b := br.Dump()
	if err := br.Load(b); err != nil {
		t.FailNow()
	}

	/* error cases */
	if err := br.Load([]byte{}); err == nil {
		t.FailNow()
	}
}

func TestBlockRequest(t *testing.T) {
	id, _ := uuid.NewV4()
	br := BlockRequest{blockN: 0, fileID: id}

	testBlockRequestDump(t, br)
	testBlockRequestLoad(t, br)
}

func TestBlockRequest_Type(t *testing.T) {
	bc := new(BlockRequest)
	if bc.Type() != MTBlockRequest {
		t.FailNow()
	}
}
