package comm

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"testing"

	"bitbucket.org/mikelsr/sakaban/fs"
	uuid "github.com/satori/go.uuid"
)

func TestMessageTypeFromBytes(t *testing.T) {
	if _, err := MessageTypeFromBytes([]byte{}); err == nil {
		t.FailNow()
	}
	if _, err := MessageTypeFromBytes([]byte{byte(0xFF)}); err == nil {
		t.FailNow()
	}
	bc := BlockContent{}
	br := BlockRequest{}
	ic := IndexContent{}
	ir := IndexRequest{}
	if mt, err := MessageTypeFromBytes(bc.Dump()); err != nil || *mt != MTBlockContent {
		t.FailNow()
	}
	if mt, err := MessageTypeFromBytes(br.Dump()); err != nil || *mt != MTBlockRequest {
		t.FailNow()
	}
	if mt, err := MessageTypeFromBytes(ic.Dump()); err != nil || *mt != MTIndexContent {
		t.FailNow()
	}
	if mt, err := MessageTypeFromBytes(ir.Dump()); err != nil || *mt != MTIndexRequest {
		t.FailNow()
	}
}

/* Block content */

// testBlockContent_Dump checks that the dumped slice has the expected length
func testBlockContentDump(t *testing.T, bc BlockContent) {
	d := bc.Dump()
	if len(d) != 19+len(bc.content) || MessageType(d[0]) != MTBlockContent {
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
	d := br.Dump()
	if len(d) != 18 || MessageType(d[0]) != MTBlockRequest {
		t.FailNow()
	}
}

func testBlockRequestLoad(t *testing.T, br BlockRequest) {
	b := br.Dump()
	if err := br.Load(b); err != nil {
		t.FailNow()
	}

	/* error case */
	if err := br.Load([]byte{}); err == nil {
		t.FailNow()
	}
}

func TestBlockRequest(t *testing.T) {
	id, _ := uuid.NewV4()
	br := BlockRequest{BlockN: 0, FileID: id}

	testBlockRequestDump(t, br)
	testBlockRequestLoad(t, br)
}

func TestBlockRequest_Type(t *testing.T) {
	bc := new(BlockRequest)
	if bc.Type() != MTBlockRequest {
		t.FailNow()
	}
}

/* Index content*/

func testIndexContentDump(t *testing.T, ic IndexContent) {
	marshalledIndex, _ := json.Marshal(ic.index)
	if len(ic.Dump()) != len(marshalledIndex)+1 {
		t.FailNow()
	}
}

func testIndexContentLoad(t *testing.T, ic IndexContent) {
	dump := ic.Dump()
	if err := ic.Load(dump); err != nil {
		t.FailNow()
	}

	/* error cases */
	if err := ic.Load(dump[1:]); err == nil {
		t.FailNow()
	}
	if err := ic.Load(dump[:len(dump)-1]); err == nil {
		t.FailNow()
	}
}

func TestIndexContent(t *testing.T) {
	muffinPath := fmt.Sprintf("%s/res/muffin.jpg", fs.ProjectPath())
	f, err := fs.MakeFile(muffinPath)
	if err != nil {
		t.FailNow()
	}
	s := fs.MakeSummary(f)
	index, err := fs.MakeIndex(s)
	if err != nil {
		t.FailNow()
	}
	ic := IndexContent{index: *index}

	testIndexContentDump(t, ic)
	testIndexContentLoad(t, ic)
}

func TestIndexContent_Type(t *testing.T) {
	ic := new(IndexContent)
	if ic.Type() != MTIndexContent {
		t.FailNow()
	}
}

/* Index request */

func testIndexRequestDump(t *testing.T, ir IndexRequest) {
	d := ir.Dump()
	if len(d) != 1 || MessageType(d[0]) != MTIndexRequest {
		t.FailNow()
	}
}

func testIndexRequestLoad(t *testing.T, ir IndexRequest) {
	d := ir.Dump()
	if err := ir.Load(d); err != nil {
		t.FailNow()
	}

	/* error cases */
	if err := ir.Load([]byte{}); err == nil {
		t.FailNow()
	}
	if err := ir.Load(make([]byte, 2)); err == nil {
		t.FailNow()
	}
}

func TestIndexRequest(t *testing.T) {
	ir := *new(IndexRequest)

	testIndexRequestDump(t, ir)
	testIndexRequestLoad(t, ir)
}

func TestIndexRequest_Type(t *testing.T) {
	ir := new(IndexRequest)
	if ir.Type() != MTIndexRequest {
		t.FailNow()
	}
}
