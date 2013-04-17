package web

import (
	"fmt"
	"iyf.cc/gospeed/utils"
	"iyf.cc/gospeed/web/session"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Context struct {
	ResponseWriter http.ResponseWriter
	Request        *http.Request
	UrlParams      map[string]string
	Params         *url.Values
	SessionStore   session.SessionStore
	sessionStart   bool
}

func (ctx *Context) WriteString(content string) {
	ctx.ResponseWriter.Write([]byte(content))
}

func (ctx *Context) Abort(body string, status int) {
	ctx.ResponseWriter.WriteHeader(status)
	ctx.ResponseWriter.Write([]byte(body))
}

func (ctx *Context) Redirect(url_ string, status int) {
	ctx.ResponseWriter.Header().Set("Location", url_)
	ctx.ResponseWriter.WriteHeader(status)
}

func (ctx *Context) NotModified() {
	ctx.ResponseWriter.WriteHeader(304)
}

func (ctx *Context) NotFound(message string) {
	ctx.ResponseWriter.WriteHeader(404)
	ctx.ResponseWriter.Write([]byte(message))
}

//Sets the content type by extension, as defined in the mime package.
//For example, ctx.ContentType("json") sets the content-type to "application/json"
func (ctx *Context) ContentType(ext string) {
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}
	ctype := mime.TypeByExtension(ext)
	if ctype != "" {
		ctx.ResponseWriter.Header().Set("Content-Type", ctype)
	}
}

func (ctx *Context) SetHeader(hdr string, val string, unique bool) {
	if unique {
		ctx.ResponseWriter.Header().Set(hdr, val)
	} else {
		ctx.ResponseWriter.Header().Add(hdr, val)
	}
}

//Sets a cookie -- duration is the amount of time in seconds. 0 = forever
func (ctx *Context) SetCookie(name string, value string, age int64) {
	var utctime time.Time
	if age == 0 {
		// 2^31 - 1 seconds (roughly 2038)
		utctime = time.Unix(2147483647, 0)
	} else {
		utctime = time.Unix(time.Now().Unix()+age, 0)
	}
	cookie := fmt.Sprintf("%s=%s; expires=%s", name, value, utils.WebTime(utctime))
	ctx.SetHeader("Set-Cookie", cookie, false)
}
func (ctx *Context) InitInput() *url.Values {
	if ctx.Params == nil {
		ct := ctx.Request.Header.Get("Content-Type")
		if strings.Contains(ct, "multipart/form-data") {
			ctx.Request.ParseMultipartForm(AppConfig.MaxMemory) //64MB
		} else {
			ctx.Request.ParseForm()
		}
		ctx.Params = &ctx.Request.Form
	}
	return ctx.Params
}

//获取传递的参数并转化为string
func (ctx *Context) ParamString(key string) string {
	return ctx.InitInput().Get(key)
}

//获取传递的参数并转化为Int64
func (ctx *Context) ParamInt64(key string) (int64, error) {
	return strconv.ParseInt(ctx.InitInput().Get(key), 10, 64)
}

//获取传递的参数并转化为Int64
func (ctx *Context) ParamBool(key string) (bool, error) {
	return strconv.ParseBool(ctx.InitInput().Get(key))
}

//获取传递的文件
func (ctx *Context) ParamFile(key string) (multipart.File, *multipart.FileHeader, error) {
	ctx.InitInput()
	return ctx.Request.FormFile(key)
}

//获取传递的参数并转化为int
func (ctx *Context) ParamInt(key string) (int, error) {
	i, err := strconv.ParseInt(ctx.InitInput().Get(key), 10, 64)
	return int(i), err
}

//获取传递的参数并转化为float64
func (ctx *Context) ParamFloat64(key string) (float64, error) {
	return strconv.ParseFloat(ctx.InitInput().Get(key), 64)
}

//返回一个session.SessionStore
//如果需要直接对这个对象操作才需要调用
func (ctx *Context) InitSession() (sess session.SessionStore) {
	if !ctx.sessionStart {
		ctx.SessionStore = GlobalSessions.SessionStart(ctx.ResponseWriter, ctx.Request)
	}
	return ctx.SessionStore
}

//用于释放session,主要是指文件保存时关闭文件
//方法一般不用手动调用
func (ctx *Context) EndSession() {
	if ctx.SessionStore != nil {
		ctx.SessionStore.SessionRelease()
	}
}

func (ctx *Context) SetSession(name string, value interface{}) {
	ctx.InitSession().Set(name, value)
}

func (ctx *Context) GetSession(name string) interface{} {
	return ctx.InitSession().Get(name)
}

func (ctx *Context) DelSession(name string) {
	ctx.InitSession().Delete(name)
}
