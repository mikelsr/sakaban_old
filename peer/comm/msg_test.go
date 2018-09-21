package comm

import (
	"log"
	"reflect"
	"testing"

	uuid "github.com/satori/go.uuid"
)

// testBlockContent_Dump checks that the dumped slice has the expected length
func testBlockContent_Dump(t *testing.T, bc BlockContent) {
	if len(bc.Dump()) != 19+len(bc.content) {
		t.FailNow()
	}
}

// testBlockContent_Load loads a BlockContent from a bc.Dump() and compares it
// to the original (bc)
func testBlockContent_Load(t *testing.T, bc BlockContent) {
	bcLoaded := new(BlockContent)
	err := bcLoaded.Load(bc.Dump())
	if err != nil {
		log.Fatalln(err)
	}
	if reflect.DeepEqual(*bcLoaded, bc) {
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

	testBlockContent_Dump(t, bc)
	testBlockContent_Load(t, bc)
}

func TestBlockContent_Type(t *testing.T) {
	bc := *new(BlockContent)
	if bc.Type() != MTBlockContent {
		t.FailNow()
	}
}
