package peer

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const testDir = "/tmp/peertest"
const testFailDir = "/tmp/peertest/fail"

var testPeer Peer

func TestMain(m *testing.M) {
	// create test peer with key pair
	tp, err := NewPeer()
	if err != nil {
		os.Exit(1)
	}
	testPeer = *tp

	// create test directories
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		os.Exit(1)
	}
	os.MkdirAll(testFailDir, 0000)

	// export necessary keys for testing
	ExportRSAKeys(testDir, testPeer.PrvKey, testPeer.PubKey)
	ExportRSAKeys(filepath.Join(testDir, "import"),
		testPeer.PrvKey, testPeer.PubKey)
	// run tests
	m.Run()
	// cleanup
	os.RemoveAll(testDir)
}

func TestPeer_Decrypt(t *testing.T) {
	// create and encrypt data
	data := []byte{42}
	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, testPeer.PubKey, data)
	if err != nil {
		t.FailNow()
	}

	// wrong seed
	_, err = testPeer.Decrypt(data, bytes.NewReader([]byte{0}))
	if err == nil {
		t.FailNow()
	}

	// correct decryption
	decryptedData, err := testPeer.Decrypt(ciphertext, rand.Reader)
	if err != nil || !bytes.Equal(decryptedData, data) {
		t.FailNow()
	}
}

func TestExportRSAKeys(t *testing.T) {
	err := ExportRSAKeys(testFailDir, testPeer.PrvKey, testPeer.PubKey)
	if err == nil {
		t.FailNow()
	}
	err = ExportRSAKeys(testDir, testPeer.PrvKey, testPeer.PubKey)
	if err != nil {
		t.FailNow()
	}
}

func TestExportRSAPrvKey(t *testing.T) {
	// export to non-writeable dir
	err := ExportRSAPrvKey(testFailDir, testPeer.PrvKey)
	if err == nil {
		t.FailNow()
	}
	// export to writeable dir
	err = ExportRSAPrvKey(filepath.Join(testDir, "export", "prv"), testPeer.PrvKey)
	if err != nil {
		t.FailNow()
	}
}

func TestExportRSAPubKey(t *testing.T) {
	// export to non-writeable dir
	err := ExportRSAPubKey(testFailDir, testPeer.PubKey)
	if err == nil {
		t.FailNow()
	}
	// export to writeable dir
	err = ExportRSAPubKey(filepath.Join(testDir, "export", "pub"), testPeer.PubKey)
	if err != nil {
		t.FailNow()
	}
}

func TestImportRSAKeys(t *testing.T) {
	_, _, err := ImportRSAKeys(testFailDir)
	if err == nil {
		t.FailNow()
	}
	_, _, err = ImportRSAKeys(testDir)
	if err != nil {
		t.FailNow()
	}
}

func TestImportRSAPrvKey(t *testing.T) {
	// import from non-existing file
	_, err := ImportRSAPrvKey(testFailDir)
	if err == nil {
		t.FailNow()
	}

	// import from invalid file
	// create invalid file
	failDir := filepath.Join(testDir, "import", "fail")
	os.MkdirAll(failDir, permissionDir)
	ioutil.WriteFile(filepath.Join(failDir, filenamePrv), []byte{0}, permissionFile)
	_, err = ImportRSAPrvKey(failDir)
	if err == nil {
		t.FailNow()
	}

	// import existing key
	_, err = ImportRSAPrvKey(filepath.Join(testDir, "import"))
	if err != nil {
		t.FailNow()
	}
}

func TestImportRSAPubKey(t *testing.T) {
	// import from non-existing file
	_, err := ImportRSAPubKey(testFailDir)
	if err == nil {
		t.FailNow()
	}

	// import from invalid file
	// create invalid file
	failDir := filepath.Join(testDir, "import", "fail")
	os.MkdirAll(failDir, permissionDir)
	ioutil.WriteFile(filepath.Join(failDir, filenamePub), []byte{0}, permissionFile)
	_, err = ImportRSAPubKey(failDir)
	if err == nil {
		t.FailNow()
	}

	// import existing key
	_, err = ImportRSAPubKey(filepath.Join(testDir, "import"))
	if err != nil {
		t.FailNow()
	}
}
