package web

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"
)

var (
	DefaultEvery int = 60 // 1 minute
)

type SpeedItem struct {
	val        interface{}
	LastAccess time.Time
	expired    int
}

func (itm *SpeedItem) Access() interface{} {
	itm.LastAccess = time.Now()
	return itm.val
}

type SpeedCache struct {
	lock  sync.RWMutex
	dur   time.Duration
	items map[string]*SpeedItem
	Every int // Run an expiration check Every seconds
}

// NewDefaultCache returns a new FileCache with sane defaults.
func NewSpeedCache() *SpeedCache {
	cache := SpeedCache{dur: time.Since(time.Now()),
		Every: DefaultEvery}
	return &cache
}

func (bc *SpeedCache) Get(name string) interface{} {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	itm, ok := bc.items[name]
	if !ok {
		return nil
	}
	return itm.Access()
}

func (bc *SpeedCache) Put(name string, value interface{}, expired int) error {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	t := SpeedItem{val: value, LastAccess: time.Now(), expired: expired}
	if _, ok := bc.items[name]; ok {
		return errors.New("the key is exist")
	} else {
		bc.items[name] = &t
	}
	return nil
}

func (bc *SpeedCache) Delete(name string) (ok bool, err error) {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	if _, ok = bc.items[name]; !ok {
		return
	}
	delete(bc.items, name)
	_, valid := bc.items[name]
	if valid {
		ok = false
	}
	return
}

func (bc *SpeedCache) IsExist(name string) bool {
	bc.lock.RLock()
	defer bc.lock.RUnlock()
	_, ok := bc.items[name]
	return ok
}

// Start activates the file cache; it will 
func (bc *SpeedCache) Start() error {
	dur, err := time.ParseDuration(fmt.Sprintf("%ds", bc.Every))
	if err != nil {
		return err
	}
	bc.dur = dur
	bc.items = make(map[string]*SpeedItem, 0)
	go bc.vaccuum()
	return nil
}

func (bc *SpeedCache) vaccuum() {
	if bc.Every < 1 {
		return
	}
	for {
		<-time.After(time.Duration(bc.dur))
		if bc.items == nil {
			return
		}
		for name, _ := range bc.items {
			bc.item_expired(name)
		}
	}
}

// item_expired returns true if an item is expired.
func (bc *SpeedCache) item_expired(name string) bool {
	bc.lock.Lock()
	defer bc.lock.Unlock()
	itm, ok := bc.items[name]
	if !ok {
		return true
	}
	dur := time.Now().Sub(itm.LastAccess)
	sec, err := strconv.Atoi(fmt.Sprintf("%0.0f", dur.Seconds()))
	if err != nil {
		delete(bc.items, name)
		return true
	} else if sec >= itm.expired {
		delete(bc.items, name)
		return true
	}
	return false
}
