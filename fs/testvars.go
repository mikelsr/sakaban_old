package fs

import (
	"fmt"
	"math/rand"
)

var (
	muffinPath = fmt.Sprintf("%s/res/muffin.jpg", ProjectPath())
	// testDir will contain the files generated for this tests
	testDir     = fmt.Sprintf("/tmp/sakaban-test-%d", rand.Intn(1e8))
	testFailDir = testDir + "/fail"
)
