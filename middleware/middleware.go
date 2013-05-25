//中间件,用于存储全局的信息
package middleware

import (
	"sync"
)

var m *SafeMap

func init() {
	m = NewSafeMap()
}
func Set(key interface{}, value interface{}) {
	m.Set(key, value)
}

//检查
func Check(key interface{}) (ok bool) {
	return m.Check(key)
}

//获取
func Get(key interface{}) (v interface{}) {
	return m.Get(key)
}

//删除
func Del(key interface{}) {
	m.lock.Lock()
	m.Del(key)
}

var mutSafe *SafeMap = NewSafeMap()

//一个可以通过用户id来实现同步锁,防止单次快速的重复提交
func GetUserMutex(key string) *sync.Mutex {
	if v := mutSafe.Get(key); v != nil {
		return v.(*sync.Mutex)
	} else {
		m := &sync.Mutex{}
		mutSafe.Set(key, m)
		return m
	}
}
