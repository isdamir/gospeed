package web

//@todo add template funcs

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/howeyc/fsnotify"
	"html/template"
	"io/ioutil"
	"iyf.cc/gospeed/log"
	"iyf.cc/gospeed/utils"
	"os"
	"path/filepath"
	"strings"
)

var (
	speedTplFuncMap template.FuncMap
	Templates       map[string]*template.Template
	TemplateExt     []string
	globalTemplate  *template.Template
	templateEven    *fsnotify.Watcher
)

func init() {
	Templates = make(map[string]*template.Template)
	speedTplFuncMap = make(template.FuncMap)
	TemplateExt = make([]string, 0)
	TemplateExt = append(TemplateExt, "tpl", "html")
	speedTplFuncMap["markdown"] = utils.MarkDown
	speedTplFuncMap["dateformat"] = utils.DateFormat
	speedTplFuncMap["date"] = utils.Date
	speedTplFuncMap["compare"] = utils.Compare
	speedTplFuncMap["substr"] = utils.Substr
	speedTplFuncMap["html2str"] = utils.Html2str
	speedTplFuncMap["str2html"] = utils.Str2html
	speedTplFuncMap["htmlquote"] = utils.Htmlquote
	speedTplFuncMap["htmlunquote"] = utils.Htmlunquote
	speedTplFuncMap["op"] = utils.Operator

	speedTplFuncMap["enurl"] = EnUrl
	var err error
	templateEven, err = fsnotify.NewWatcher()
	if err != nil {
		log.Trace(err)
	}
}

// AddFuncMap let user to register a func in the template
func AddFuncMap(key string, funname interface{}) error {
	if _, ok := speedTplFuncMap[key]; ok {
		return errors.New("funcmap already has the key")
	}
	speedTplFuncMap[key] = funname
	return nil
}

// type templatefile struct {
// 	root  string
// 	files map[string][]string
// }

func IsTemplate(path string) (b bool) {
	for _, v := range TemplateExt {
		if strings.HasSuffix(path, v) {
			return true
		}
	}
	return false
}

func AddTemplateExt(ext string) {
	for _, v := range TemplateExt {
		if v == ext {
			return
		}
	}
	TemplateExt = append(TemplateExt, ext)
}

//监听某个目录的模板
func WatchTemplate() (err error) {
	//取出目录中的gobal部分
	buildGlobalTemplate(AppConfig.ViewsPath + "_global")
	err = filepath.Walk(AppConfig.ViewsPath, func(pa string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if f.IsDir() {
			templateEven.Watch(pa)
		}
		return nil
	})
	if err != nil {
		log.Debug("filepath.Walk() returned %v\n", err)
		return
	}
	templateEven.Watch(AppConfig.ViewsPath)
	go startWatch()
	return
}

func startWatch() {
	for {
		select {
		case v := <-templateEven.Event:
			{
				log.Debug(v.String())
				log.Trace("File:", v.String())
				fi, err := os.Stat(v.Name)
				if err == nil {
					log.Debug(v.String())
					if v.IsCreate() || v.IsModify() || v.IsRename() {
						if strings.Index(v.Name, "_global") != -1 {
							buildGlobalTemplate(AppConfig.ViewsPath + "_global")
						} else if !fi.IsDir() && IsTemplate(v.Name) {
							buildTemplate(v.Name)
						}
						if fi.IsDir() {
							templateEven.Watch(v.Name)
						}
					}
					if v.IsDelete() {
						templateEven.RemoveWatch(v.Name)
					}
				}
			}
		case err := <-templateEven.Error:
			log.Debug("template error:", err)
		}
	}
}
func buildAllTemplate(dir string) (err error) {
	err = filepath.Walk(dir, func(pa string, f os.FileInfo, err error) error {
		if f == nil {
			return err
		}
		if !f.IsDir() {
			if strings.Index(pa, "_global") == -1 && IsTemplate(pa) {
				buildTemplate(pa)
			}
		}
		return nil
	})
	return
}
func buildGlobalTemplate(dir string) {
	if _, err := os.Stat(dir); err != nil {
		log.Debug(err)
	} else {
		var file []string
		err := filepath.Walk(dir, func(pa string, f os.FileInfo, err error) error {
			if f == nil {
				return err
			}
			if f.IsDir() || (f.Mode()&os.ModeSymlink) > 0 {
				return nil
			} else {
				if IsTemplate(pa) {
					file = append(file, pa)
				}
			}
			return nil
		})
		if err != nil {
			log.Debug("filepath.Walk() returned %v\n", err)
		}
		t, err := template.New("speedTemplateGlobal").Funcs(speedTplFuncMap).ParseFiles(file...)
		if err == nil && t != nil {
			globalTemplate = t
		} else {
			log.Debug(err)
		}
	}
	buildAllTemplate(AppConfig.ViewsPath)
}
func buildTemplate(file string) {
	if strings.HasPrefix(file, "./") {
		file = file[2:]
	}
	if _, err := os.Stat(file); err != nil {
		log.Debug(err)
		return
	}
	file = strings.Replace(file, "\\", "/", -1)
	file = strings.Replace(file, "//", "/", -1)
	log.Trace("build template", file)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	s := string(b)
	g := GlobalTemplate()
	if g != nil {
		g, _ = g.Clone()
		t, err := g.New(file).Funcs(speedTplFuncMap).Parse(s)
		if err == nil && t != nil {
			Templates[file] = t
		} else {
			log.Debug(err)
		}
	} else {
		t, err := template.New(file).Funcs(speedTplFuncMap).Parse(s)
		if err == nil && t != nil {
			Templates[file] = t
		} else {
			log.Debug(err)
		}
	}
}

//取得全局模板
func GlobalTemplate() (t *template.Template) {
	return globalTemplate
}

//获取到一个模板
//模板会解析全局的模板
func GetTemplate(file string) (t *template.Template) {
	return Templates[file]
}
func ExsitTemplate(file string) bool {
	_, b := Templates[fmt.Sprint(AppConfig.ViewsPath, file)]
	return b
}

//解析模板
func RenderTemplate(file string, data map[string]interface{}, tm template.FuncMap) (wr *bytes.Buffer, err error) {
	file = fmt.Sprint(AppConfig.ViewsPath, file)
	wr = &bytes.Buffer{}
	t := GetTemplate(file)
	if t != nil {
		log.Trace("Render:", file)
		if len(tm) == 0 {
			err = t.ExecuteTemplate(wr, file, data)
		} else {
			err = t.Funcs(tm).ExecuteTemplate(wr, file, data)
		}
		if err != nil {
			log.Debug("Render:", file, err)
		}
	} else {
		err = errors.New("no template")
		log.Debug("Render:", file, "no template")
	}
	return
}
