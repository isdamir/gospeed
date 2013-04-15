package middleware

import (
	"sync"
)

type SafeMap struct {
	lock *sync.RWMutex
	bm   map[interface{}]interface{}
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		lock: new(sync.RWMutex),
		bm:   make(map[interface{}]interface{}),
	}
}

//设置
func (m *SafeMap) Set(key interface{}, value interface{}) {
	m.lock.Lock()
	m.bm[key] = value
	m.lock.Unlock()
}

//检查
func (m *SafeMap) Check(key interface{}) (ok bool) {
	_, ok = m.bm[key]
	return
}

//获取
func (m *SafeMap) Get(key interface{}) (v interface{}) {
	v, _ = m.bm[key]
	return
}

//删除
func (m *SafeMap) Del(key interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	delete(m.bm, key)
}
