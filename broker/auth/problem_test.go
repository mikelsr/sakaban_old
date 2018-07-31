package auth

import (
	"testing"
)

func TestMakeProblemFromString(t *testing.T) {
	// invalid strings
	_, err := MakeProblemFromString("")
	if err == nil {
		t.FailNow()
	}
	_, err = MakeProblemFromString("_*_")
	if err == nil {
		t.FailNow()
	}
	_, err = MakeProblemFromString("15*_")
	if err == nil {
		t.FailNow()
	}

	// valid string
	p1 := MakeProblem(46, 2)
	p2, err := MakeProblemFromString(p1.Formulate())
	if err != nil || p1.c1 != p2.c1 || p1.c2 != p2.c2 {
		t.FailNow()
	}
}

func TestProblem_Formulate(t *testing.T) {
	p1 := NewProblem()
	p2, _ := MakeProblemFromString(p1.Formulate())
	if p1.c1 != p2.c1 || p1.c2 != p2.c2 {
		t.FailNow()
	}
}

func TestProblem_Solution(t *testing.T) {
	p := MakeProblem(46, 2)
	if p.Solution() != 46*2 {
		t.FailNow()
	}
}
