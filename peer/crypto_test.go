package peer

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"bitbucket.org/mikelsr/sakaban-broker/auth"
)

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

func TestPeer_RSADecrypt(t *testing.T) {
	p, _ := NewPeer()
	data := []byte{42}
	encrypted := auth.RSAEncrypt(p.PubKey, data)

	if !bytes.Equal(p.RSADecrypt(encrypted), data) {
		t.FailNow()
	}
}
