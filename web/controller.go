package web

import (
	"compress/gzip"
	"compress/zlib"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"iyf.cc/gospeed/log"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type Controller struct {
	Ctx       *Context
	Data      map[string]interface{}
	ChildName string
	tplName   string
	//将这里面的模板进行解析,并将结果存到Data[key]中用于输出
	tplIn  map[string]string
	TplExt string
}

type ControllerInterface interface {
	Init(ct *Context, cn string, up map[string]string)
	Prepare()
	Get()
	Post()
	Delete()
	Put()
	Head()
	Patch()
	Link()
	UnLink()
	Purge()
	Options()
	Finish()
	Render() error
}

func (c *Controller) Init(ctx *Context, cn string, up map[string]string) {
	c.Data = make(map[string]interface{})
	c.tplIn = make(map[string]string)
	c.tplName = ""
	c.ChildName = cn
	c.Ctx = ctx
	c.TplExt = "html"
}

func (c *Controller) Prepare() {

}

func (c *Controller) Finish() {
}

func (c *Controller) Get() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Post() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Delete() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Put() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Head() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Patch() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Options() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Link() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) UnLink() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Purge() {
	http.Error(c.Ctx.ResponseWriter, "Method Not Allowed", 405)
}

func (c *Controller) Render() error {
	rb, err := c.RenderBytes()

	if err != nil {
		return err
	} else {
		c.Ctx.ResponseWriter.Header().Set("Content-Type", "text/html; charset=utf-8")
		output_writer := c.Ctx.ResponseWriter.(io.Writer)
		if AppConfig.EnableGzip == true && c.Ctx.Request.Header.Get("Accept-Encoding") != "" {
			splitted := strings.SplitN(c.Ctx.Request.Header.Get("Accept-Encoding"), ",", -1)
			encodings := make([]string, len(splitted))

			for i, val := range splitted {
				encodings[i] = strings.TrimSpace(val)
			}
			for _, val := range encodings {
				if val == "gzip" {
					c.Ctx.ResponseWriter.Header().Set("Content-Encoding", "gzip")
					output_writer, _ = gzip.NewWriterLevel(c.Ctx.ResponseWriter, gzip.BestSpeed)

					break
				} else if val == "deflate" {
					c.Ctx.ResponseWriter.Header().Set("Content-Encoding", "deflate")
					output_writer, _ = zlib.NewWriterLevel(c.Ctx.ResponseWriter, zlib.BestSpeed)
					break
				}
			}
		} else {
			c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(rb)), true)
		}
		output_writer.Write(rb)
		switch output_writer.(type) {
		case *gzip.Writer:
			output_writer.(*gzip.Writer).Close()
		case *zlib.Writer:
			output_writer.(*zlib.Writer).Close()
		case io.WriteCloser:
			output_writer.(io.WriteCloser).Close()
		}
		return nil
	}
	return nil
}

func (c *Controller) RenderString() (string, error) {
	b, e := c.RenderBytes()
	return string(b), e
}

func (c *Controller) RenderBytes() ([]byte, error) {
	if c.tplName == "" {
		c.tplName = fmt.Sprint(path.Join(c.ChildName, c.Ctx.Request.Method), ".", c.TplExt)
	}
	c.Data["Custom"] = AppConfig.Custom
	c.Data["Browser"] = c.Ctx.Browser
	if c.Ctx.sessionStart {
		mp := c.Ctx.Session().Map()
		if _, ok := mp["__ToUrl"]; ok {
			mp["SessionID"] = c.Ctx.Session().SessionID()
		}
		c.Data["Session"] = mp
	}
	if len(c.tplIn) > 0 {
		for k, v := range c.tplIn {
			buf, err := RenderTemplate(v, c.Data)
			if err != nil {
				log.Debug(err)
				continue
			}
			c.Data[k] = template.HTML(buf.String())
		}
	}
	buf, err := RenderTemplate(c.tplName, c.Data)
	if err != nil {
		log.Trace("template Execute err:", err)
	}
	return buf.Bytes(), nil
}

func (c *Controller) Redirect(url string, code int) {
	c.Ctx.Redirect(url, code)
}

func (c *Controller) ServeJson(data interface{}) {
	content, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	c.Ctx.ResponseWriter.Write(content)
}

func (c *Controller) ServeXml(data interface{}) {
	content, err := xml.Marshal(data)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	c.Ctx.ResponseWriter.Header().Set("Content-Type", "application/xml")
	c.Ctx.ResponseWriter.Write(content)
}
func (c *Controller) ServeTpl(tplpath string) {
	if AppConfig.AutoDevice {
		ext := filepath.Ext(tplpath)
		file := tplpath[:len(tplpath)-len(ext)]
		c.tplName = c.templatePath(file, ext)
		log.Debug(file, ext)
	} else {
		c.tplName = tplpath
	}
}
func (c *Controller) ServetplIn(tplIn map[string]string) {
	c.tplIn = tplIn
}

//针对不同浏览器进行解析
func (c *Controller) templatePath(path, ext string) string {
	if c.Ctx.Browser.IsMobile {
		if c.Ctx.Browser.IsWml {
			t := fmt.Sprintf("%s_wml%s", path, ext)
			if ExsitTemplate(t) {
				return t
			}
		}
		t := fmt.Sprintf("%s_html5%s", path, ext)
		if ExsitTemplate(t) {
			return t
		}
	}
	return fmt.Sprintf("%s%s", path, ext)
}
