package hex

import (
	"bytes"
	"testing"
	"time"
	/*"time"*/)

func TestTimeToByte(t *testing.T) {
	da := time.Date(2014, 8, 29, 10, 22, 0, 0, time.UTC)
	data := TimeToByte(da)
	if !(data[0] == 0 && data[1] == 22 && data[2] == 10 && data[3] == 29 && data[4] == 8 && data[5] == 5 && data[6] == 14) {
		t.Error(data)
		t.FailNow()
	}
}
func TestInt16(t *testing.T) {
	data := make([]byte, 2)
	data[0] = 0x01
	data[1] = 0xf4
	if 500 != BytesToInt16(data) {
		t.Error("BytesToInt16")
		t.Error(BytesToInt16(data))
		t.FailNow()
	}
	if !bytes.Equal(data, Int16ToBytes(500)) {
		t.Error(Int16ToBytes(500), data)
		t.Error("Int16ToBytes")
		t.FailNow()
	}
}
func TestInt32(t *testing.T) {
	var i int32 = 2323
	if i != BytesToInt32(Int32ToBytes(i)) {
		t.Error(BytesToInt32(Int32ToBytes(i)))
		t.FailNow()
	}
}
func TestInt64(t *testing.T) {
	var i int64 = 2323
	if i != BytesToInt64(Int64ToBytes(i)) {
		t.Error(BytesToInt64(Int64ToBytes(i)))
		t.FailNow()
	}
}
