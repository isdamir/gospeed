package nettool

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	nt := NetTool{}
	str, err := nt.Get("http://360.cn")
	fmt.Println(str, err)
}
