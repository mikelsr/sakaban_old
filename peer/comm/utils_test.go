package comm

import (
	"bytes"
	"testing"
)

func TestUint16FromBytes(t *testing.T) {
	b := []byte{}
	expected := uint16(0)
	if uint16FromBytes(b) != expected {
		t.FailNow()
	}

	b = []byte{0, 0}
	if uint16FromBytes(b) != expected {
		t.FailNow()
	}

	b = []byte{0, 1}
	expected = uint16(256)
	if uint16FromBytes(b) != expected {
		t.FailNow()
	}
}

func TestUint16ToBytes(t *testing.T) {
	n := uint16(0)
	expected := []byte{0, 0}
	if !bytes.Equal(uint16ToBytes(n), expected) {
		t.FailNow()
	}
	n = uint16(256)
	expected = []byte{0, 1}
	if !bytes.Equal(uint16ToBytes(n), expected) {
		t.FailNow()
	}
}

func TestUint64FromBytes(t *testing.T) {
	b := []byte{}
	expected := uint64(0)
	if uint64FromBytes(b) != expected {
		t.FailNow()
	}

	b = []byte{0, 0, 0, 0, 0, 0, 0, 0}
	if uint64FromBytes(b) != expected {
		t.FailNow()
	}

	b = []byte{0, 1, 0, 0, 0, 0, 0, 0}
	expected = uint64(256)
	if uint64FromBytes(b) != expected {
		t.FailNow()
	}
}

func TestUint64ToBytes(t *testing.T) {
	n := uint64(0)
	expected := []byte{0, 0, 0, 0, 0, 0, 0, 0}
	if !bytes.Equal(uint64ToBytes(n), expected) {
		t.FailNow()
	}
	n = uint64(256)
	expected = []byte{0, 1, 0, 0, 0, 0, 0, 0}
	if !bytes.Equal(uint64ToBytes(n), expected) {
		t.FailNow()
	}
}
