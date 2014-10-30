package hex

import (
	"encoding/binary"
	"errors"
	"time"
)

//全部使用大端
func Int16ToBytes(i int16) []byte {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(i))
	return buf
}

func BytesToInt16(buf []byte) int16 {
	return int16(binary.BigEndian.Uint16(buf))
}
func Int32ToBytes(i int32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(i))
	return buf
}

func BytesToInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

//秒 分 时 日 月 星期 年(表示20范围,去掉20)
func TimeToByte(t time.Time) []byte {
	data := make([]byte, 7)
	data[0] = byte(t.Second())
	data[1] = byte(t.Minute())
	data[2] = byte(t.Hour())
	data[3] = byte(t.Day())
	data[4] = byte(t.Month())
	data[5] = byte(t.Weekday())           //星期一是0
	data[6] = byte(t.Year() - ReduceYeer) //2000年开始
	return data
}

const (
	ReduceYeer int = 2000
)

//秒 分 时 日 月 星期 年(表示20范围,去掉20)
func ByteToTime(data []byte) (t time.Time, err error) {
	if len(data) != 7 {
		err = errors.New("data is faild")
		return
	}
	t = time.Date(int(data[6])+ReduceYeer, time.Month(int(data[4])), int(data[3]), int(data[2]), int(data[1]), int(data[0]), time.Now().Nanosecond(), time.UTC)
	return
}
