package peer

import (
	"testing"

	"bitbucket.org/mikelsr/sakaban/fs"
)

func TestFileStack_Pop(t *testing.T) {
	s1 := fs.Summary{Path: "1"}
	s2 := fs.Summary{Path: "2"}
	stack := fileStack{[]*fs.Summary{&s1, &s2}}
	stack.pop()
	s, n := stack.pop()
	if n != 0 || s.Path != s1.Path {
		t.FailNow()
	}
	s, n = stack.pop()
	if n != -1 || s != nil {
		t.FailNow()
	}
}

func TestFileStack_Push(t *testing.T) {
	s1 := fs.Summary{Path: "1"}
	s2 := fs.Summary{Path: "2"}

	stack := newFileStack()
	stack.push(&s1)
	stack.push(&s2)

	if stack.files[len(stack.files)-1].Path != s2.Path {
		t.FailNow()
	}
}
