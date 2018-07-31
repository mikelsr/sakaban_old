package auth

import (
	"testing"
)

func TestExtractPubKey(t *testing.T) {
	// invalid keys
	_, err := ExtractPubKey(" ")
	if err == nil {
		t.FailNow()
	}
	_, err = ExtractPubKey("")
	if err == nil {
		t.FailNow()
	}
	// valid key
	_, err = ExtractPubKey(pubKey)
	if err != nil {
		t.FailNow()
	}
}

func TestPrintPubKey(t *testing.T) {
	key, _ := ExtractPubKey(pubKey)
	keyStr := PrintPubKey(key)
	if keyStr != pubKey {
		t.FailNow()
	}
}
