package middleware

import (
	"testing"
)

func Test_safemap(t *testing.T) {
	bm := NewSafeMap()
	bm.Set("test", 1)
	if !bm.Check("test") {
		t.Error("check err")
	}

	if v := bm.Get("test"); v.(int) != 1 {
		t.Error("get err")
	}

	bm.Del("test")
	if bm.Check("test") {
		t.Error("delete err")
	}
}
