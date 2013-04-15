package middleware

import (
	"testing"
)

type test struct {
	I int
	S string
}

func TestSet(t *testing.T) {
	tem := test{1000, "test two"}
	M.Set("test", tem)
	ok := M.Check("test")
	if !ok {
		t.FailNow()
	}

}
func TestDel(t *testing.T) {
	te := test{100, "test it"}
	M.Set("test", te)
	M.Del("test")
	ok := M.Check("test")
	if ok {
		t.FailNow()
	}

}
