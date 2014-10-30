package session

import (
	"bytes"
	"encoding/gob"
	"github.com/isdamir/gospeed/log"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"
)

var (
	filepder      = &FileProvider{}
	gcmaxlifetime int64
)

type FileSessionStore struct {
	f      string
	sid    string
	lock   sync.RWMutex
	values map[interface{}]interface{}
}

func (st *FileSessionStore) Map() map[interface{}]interface{} {
	return st.values
}
func (fs *FileSessionStore) Set(key, value interface{}) error {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	fs.values[key] = value
	return nil
}

func (fs *FileSessionStore) Get(key interface{}) interface{} {
	fs.lock.RLock()
	defer fs.lock.RUnlock()
	if v, ok := fs.values[key]; ok {
		return v
	} else {
		return nil
	}
	return nil
}

func (fs *FileSessionStore) Delete(key interface{}) error {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	delete(fs.values, key)
	return nil
}

func (fs *FileSessionStore) SessionID() string {
	return fs.sid
}

func (fs *FileSessionStore) SessionRelease() {
	fs.updatecontent()
}

func (fs *FileSessionStore) updatecontent() {
	fs.lock.Lock()
	defer fs.lock.Unlock()
	if len(fs.values) > 0 {
		b, err := encodeGob(fs.values)
		if err == nil {
			ioutil.WriteFile(fs.f, b, os.ModePerm)
		}
	}
}

type FileProvider struct {
	maxlifetime int64
	savePath    string
}

func (fp *FileProvider) SessionInit(maxlifetime int64, savePath string) error {
	fp.maxlifetime = maxlifetime
	fp.savePath = savePath
	return nil
}

func (fp *FileProvider) SessionRead(sid string) (SessionStore, error) {
	pf := path.Join(fp.savePath, string(sid[0]), string(sid[1]), sid)
	err := os.MkdirAll(path.Join(fp.savePath, string(sid[0]), string(sid[1])), 0777)
	if err != nil {
		log.Debug(err)
	}
	ss := &FileSessionStore{}
	ss.sid = sid
	ss.f = pf
	ss.lock.RLock()
	defer ss.lock.RUnlock()
	var kv map[interface{}]interface{}
	b, err := ioutil.ReadFile(pf)
	if err != nil {
		ss.values = make(map[interface{}]interface{})
		return ss, err
	}
	if len(b) == 0 {
		kv = make(map[interface{}]interface{})
	} else {
		kv, err = decodeGob(b)
		if err != nil {
			return nil, err
		}
	}
	ss.values = kv
	return ss, nil
}

func (fp *FileProvider) SessionDestroy(sid string) error {
	os.Remove(path.Join(fp.savePath))
	return nil
}

func (fp *FileProvider) SessionGC() {
	gcmaxlifetime = fp.maxlifetime
	filepath.Walk(fp.savePath, gcpath)
}

func gcpath(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}
	if info.IsDir() {
		return nil
	}
	if (info.ModTime().Unix() + gcmaxlifetime) < time.Now().Unix() {
		os.Remove(path)
	}
	return nil
}

func init() {
	Register("file", filepder)
	gob.Register([]interface{}{})
	gob.Register(map[int]interface{}{})
	gob.Register(map[string]interface{}{})
	gob.Register(map[interface{}]interface{}{})
	gob.Register(map[string]string{})
	gob.Register(map[int]string{})
	gob.Register(map[int]int{})
	gob.Register(map[int]int64{})
}

func encodeGob(obj map[interface{}]interface{}) ([]byte, error) {
	buf := bytes.NewBuffer(nil)
	enc := gob.NewEncoder(buf)
	err := enc.Encode(obj)
	if err != nil {
		return []byte(""), err
	}
	return buf.Bytes(), nil
}

func decodeGob(encoded []byte) (map[interface{}]interface{}, error) {
	buf := bytes.NewBuffer(encoded)
	dec := gob.NewDecoder(buf)
	var out map[interface{}]interface{}
	err := dec.Decode(&out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
