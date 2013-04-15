package config

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/howeyc/fsnotify"
	"io/ioutil"
	"iyf.cc/gospeed/log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	object map[string]*configData
	nick   map[string]*configData
	w      *fsnotify.Watcher
	watch  map[string]bool
}
type configData struct {
	filetype string
	filepath string
	nickname string
	data     interface{}
}

var cf *Config

//尽量都使用这个来获取单例的指针
func GetConfig() *Config {
	if cf == nil {
		cf = &Config{map[string]*configData{}, map[string]*configData{}, nil, map[string]bool{}}
		cf.w, _ = fsnotify.NewWatcher()
		go cf.startWatch()
	}
	return cf
}

//通过传入一个路径,昵称,一个interface来注册一个配置
//程序自动通过判断类型,建议使用相应的后缀,否则程序将依次尝试
//新建文件的话如果无后缀则使用json类型
func (c *Config) Register(path, nickname string, in interface{}) error {
	fi, err := os.Stat(path)
	if err != nil {
		stype, err := getFileType(path, in)
		if err != nil {
			stype = "json"
		}
		cf := &configData{stype, path, nickname, in}
		c.object[path] = cf
		c.nick[nickname] = cf
		c.addWatch(cf)
		return err
	}
	if fi.IsDir() {
		return errors.New("path must a file,not dir")
	}
	if _, ok := c.object[path]; ok {
		return errors.New("file is registerd,you can use Get(path)")
	}
	if _, ok := c.nick[path]; ok {
		return errors.New("nickname is use,you must change")
	}
	stype, err := getFileType(path, in)
	if err != nil {
		return err
	}
	cf := &configData{stype, path, nickname, in}
	err = readConfig(cf)
	if err != nil {
		return err
	}
	c.object[path] = cf
	c.nick[nickname] = cf
	c.addWatch(cf)
	return nil
}

//增加一个监视器,只监视修改
func (c *Config) addWatch(cf *configData) {
	dir := filepath.Dir(cf.filepath)
	if _, ok := c.watch[dir]; !ok {
		c.watch[dir] = true
		c.w.WatchFlags(dir, fsnotify.FSN_MODIFY)
	}
}

//通过nickname进行查询,无法找到则使用file查询
func (c *Config) Save(name string) error {
	var en *configData
	var ok bool
	if en, ok = c.nick[name]; !ok {
		if en, ok = c.object[name]; !ok {
			return errors.New("not found name with interface")
		}
	}
	switch en.filetype {
	case ".json":
		{
			b, err := json.Marshal(en.data)
			if err != nil {
				return err
			}
			return ioutil.WriteFile(en.filepath, b, 0644)
		}
	case ".gob":
		{
			var fout bytes.Buffer
			enc := gob.NewEncoder(&fout)
			err := enc.Encode(en.data)
			if err != nil {
				return err
			}
			return ioutil.WriteFile(en.filepath, fout.Bytes(), 0644)
		}
	case ".xml":
		{
			b, err := xml.Marshal(en.data)
			if err != nil {
				return err
			}
			return ioutil.WriteFile(en.filepath, b, 0644)
		}
	}
	return nil
}

//通过nickname进行查询,无法找到则使用file查询,查询后删除信息
func (c *Config) Close(name string) {
	if v, ok := c.nick[name]; ok {
		delete(c.nick, v.nickname)
		delete(c.object, v.filepath)
	}
	if v, ok := c.object[name]; ok {
		delete(c.nick, v.nickname)
		delete(c.object, v.filepath)
	}
}
func (c *Config) Get(name string) (interface{}, error) {
	if v, ok := c.nick[name]; ok {
		return v.data, nil
	}
	if v, ok := c.object[name]; ok {
		return v.data, nil
	}
	return nil, errors.New("not found name with interface")
}
func (c *Config) startWatch() {
	for {
		select {
		case v := <-c.w.Event:
			{
				file := ""
				if strings.HasPrefix(v.Name, "./") {
					file = v.Name[2:]
				} else {
					file = v.Name
				}
				if cf, ok := c.object[file]; ok {
					log.Trace("read:", file)
					readConfig(cf)
				}
			}
		case err := <-c.w.Error:
			log.Debug("error:", err)
		}
	}
}

func getFileType(path string, in interface{}) (string, error) {
	ext := filepath.Ext(path)
	if ext == ".dat" {
		ext = ".gob"
	}
	if ext == "" {
		bi, err := ioutil.ReadFile(path)
		if err != nil {
			return "", errors.New("not read file,can not judge type")
		}
		err = json.Unmarshal(bi, in)
		if err == nil {
			return ".json", nil
		}
		err = xml.Unmarshal(bi, in)
		if err == nil {
			return ".xml", nil
		}
		return ".gob", nil
	}
	return ext, nil
}
func readConfig(cf *configData) error {
	bi, err := ioutil.ReadFile(cf.filepath)
	if err != nil {
		return err
	}
	switch cf.filetype {
	case ".json":
		{
			return json.Unmarshal(bi, cf.data)
		}
	case ".gob":
		{
			fin, err := os.Open(cf.filepath)
			if err != nil {
				return err
			}
			dec := gob.NewDecoder(fin)
			err = dec.Decode(cf.data)
			return err
		}
	case ".xml":
		{
			return xml.Unmarshal(bi, cf.data)
		}
	}
	return nil
}
