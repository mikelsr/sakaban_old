package auth

import (
	"encoding/json"
	"testing"
)

func TestRequest_String(t *testing.T) {
	problem := NewProblem()
	r1 := MakeRequest(pub, problem, "testtoken")
	str := r1.String()
	r2 := new(Request)
	err := json.Unmarshal([]byte(str), r2)
	if err != nil {
		t.FailNow()
	}
}
