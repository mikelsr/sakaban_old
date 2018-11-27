package peer

import (
	"testing"

	"bitbucket.org/mikelsr/sakaban/fs"
)

func TestFileStack_iterFile(t *testing.T) {
	s1 := fs.Summary{Path: "1"}
	s2 := fs.Summary{Path: "2"}
	stack := fileStack{files: []*fs.Summary{&s1, &s2}}
	stack.iterFile()
	s := stack.peek()
	if s == nil || s.Path != s1.Path {
		t.FailNow()
	}
	stack.iterFile()
	s = stack.peek()
	if s != nil {
		t.FailNow()
	}
}

func TestFileStack_peek(t *testing.T) {
	s1 := fs.Summary{Path: "1"}
	stack := fileStack{files: []*fs.Summary{}}
	if stack.peek() != nil {
		t.FailNow()
	}
	stack.push(&s1)
	if stack.peek().Path != s1.Path {
		t.FailNow()
	}
}

func TestFileStack_pop(t *testing.T) {
	s1 := fs.Summary{Path: "1"}
	s2 := fs.Summary{Path: "2"}
	stack := fileStack{files: []*fs.Summary{&s1, &s2}}
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

func TestFileStack_push(t *testing.T) {
	s1 := fs.Summary{Path: "1"}
	s2 := fs.Summary{Path: "2"}

	stack := newFileStack()
	stack.push(&s1)
	stack.push(&s2)

	if stack.files[len(stack.files)-1].Path != s2.Path {
		t.FailNow()
	}
}
