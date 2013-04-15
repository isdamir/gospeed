package web

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"html/template"
	. "iyf.cc/gospeed/log"
	"net/http"
	"path"
	"strconv"
)

type Controller struct {
	Ctx       *Context
	Data      map[string]interface{}
	ChildName string
	TplName   string
	//将这里面的模板进行解析,并将结果存到Data[key]中用于输出
	TplIn  map[string]string
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
	Options()
	Finish()
	Render() error
}

func (c *Controller) Init(ctx *Context, cn string, up map[string]string) {
	c.Data = make(map[string]interface{})
	c.TplIn = make(map[string]string)
	c.TplName = ""
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

func (c *Controller) Render() error {
	rb, err := c.RenderBytes()

	if err != nil {
		return err
	} else {
		c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(rb)), true)
		c.Ctx.ContentType("text/html")
		c.Ctx.ResponseWriter.Write(rb)
		return nil
	}
	return nil
}

func (c *Controller) RenderString() (string, error) {
	b, e := c.RenderBytes()
	return string(b), e
}

func (c *Controller) RenderBytes() ([]byte, error) {
	if c.TplName == "" {
		c.TplName = fmt.Sprint(path.Join(c.ChildName, c.Ctx.Request.Method), ".", c.TplExt)
	}
	c.Data["Custom"] = AppConfig.Custom
	c.Data["Session"] = c.Ctx.SessionStore.Map()
	if len(c.TplIn) > 0 {
		for k, v := range c.TplIn {
			buf, err := RenderTemplate(v, c.Data)
			if err != nil {
				Debug(err)
				continue
			}
			c.Data[k] = template.HTML(buf.String())
		}
	}
	buf, err := RenderTemplate(c.TplName, c.Data)
	if err != nil {
		Trace("template Execute err:", err)
	}
	return buf.Bytes(), nil
}

func (c *Controller) Redirect(url string, code int) {
	c.Ctx.Redirect(url, code)
}

func (c *Controller) ServeJson() {
	content, err := json.MarshalIndent(c.Data, "", "  ")
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	c.Ctx.ContentType("json")
	c.Ctx.ResponseWriter.Write(content)
}

func (c *Controller) ServeXml() {
	content, err := xml.Marshal(c.Data)
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ctx.SetHeader("Content-Length", strconv.Itoa(len(content)), true)
	c.Ctx.ContentType("xml")
	c.Ctx.ResponseWriter.Write(content)
}
