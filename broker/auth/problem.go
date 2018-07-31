package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// Problem is used to authenticate peers
type Problem struct {
	c1  int64
	c2  int64
	sol int64
}

// NewProblem creates a simple multiplication and it's solution
func NewProblem() Problem {
	rand.Seed(time.Now().UnixNano())
	c1 := rand.Int63()
	c2 := rand.Int63()
	return Problem{c1: c1, c2: c2, sol: c1 * c2}
}

// MakeProblem creates a problem given it's components
func MakeProblem(c1 int64, c2 int64) Problem {
	return Problem{c1: c1, c2: c2}
}

// MakeProblemFromString creates a problem given it's formulation string
func MakeProblemFromString(str string) (*Problem, error) {
	strs := strings.Split(str, "*")
	if len(strs) != 2 {
		return nil, errors.New("Invalid string for a problem")
	}
	c1, err := strconv.ParseInt(strs[0], 10, 64)
	if err != nil {
		return nil, err
	}
	c2, err := strconv.ParseInt(strs[1], 10, 64)
	if err != nil {
		return nil, err
	}
	return &Problem{
		c1: c1,
		c2: c2,
	}, nil
}

// Formulate formulates the problem as a string
func (p *Problem) Formulate() string {
	return fmt.Sprintf("%d*%d", p.c1, p.c2)
}

// Solution returns the solution of a problem
func (p *Problem) Solution() int64 {
	return p.c1 * p.c2
}
