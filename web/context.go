package web

import (
	"fmt"
	"iyf.cc/gospeed/browser"
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
	Browser        *browser.BrowserCheck
	sessionStore   session.SessionStore
	params         *url.Values
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
func EnUrl(sess map[interface{}]interface{}, ul string) string {
	if sess != nil {
		if v, ok := sess["SessionID"]; ok && AppConfig.SessionToUrl {
			vs := v.(string)
			if strings.Index(vs, "?") == -1 {
				return fmt.Sprintf("%s?%s=%s", ul, AppConfig.SessionName, url.QueryEscape(vs))
			}
			if strings.HasSuffix(vs, "&") {
				return fmt.Sprintf("%s=%s", ul, AppConfig.SessionName, url.QueryEscape(vs))
			} else {
				return fmt.Sprintf("&%s=%s", ul, AppConfig.SessionName, url.QueryEscape(vs))
			}
		}
	}
	return ul
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
func (ctx *Context) Params() *url.Values {
	if ctx.params == nil {
		ct := ctx.Request.Header.Get("Content-Type")
		if strings.Contains(ct, "multipart/form-data") {
			ctx.Request.ParseMultipartForm(AppConfig.MaxMemory) //64MB
		} else {
			ctx.Request.ParseForm()
		}
		ctx.params = &ctx.Request.Form
	}
	return ctx.params
}

//获取传递的参数并转化为string
func (ctx *Context) ParamString(key string) string {
	return ctx.Params().Get(key)
}

//获取传递的参数并转化为Int64
func (ctx *Context) ParamInt64(key string) (int64, error) {
	return strconv.ParseInt(ctx.Params().Get(key), 10, 64)
}

//获取传递的参数并转化为Int64
func (ctx *Context) ParamBool(key string) (bool, error) {
	return strconv.ParseBool(ctx.Params().Get(key))
}

//获取传递的文件
func (ctx *Context) ParamFile(key string) (multipart.File, *multipart.FileHeader, error) {
	ctx.Params()
	return ctx.Request.FormFile(key)
}

//获取传递的参数并转化为int
func (ctx *Context) ParamInt(key string) (int, error) {
	i, err := strconv.ParseInt(ctx.Params().Get(key), 10, 64)
	return int(i), err
}

//获取传递的参数并转化为float64
func (ctx *Context) ParamFloat64(key string) (float64, error) {
	return strconv.ParseFloat(ctx.Params().Get(key), 64)
}

//返回一个session.SessionStore
func (ctx *Context) Session() (sess session.SessionStore) {
	if !ctx.sessionStart {
		ctx.sessionStore = GlobalSessions.SessionStart(ctx.ResponseWriter, ctx.Request)
		ctx.sessionStart = true
	}
	return ctx.sessionStore
}

//用于释放session,主要是指文件保存时关闭文件
//方法一般不用手动调用
func (ctx *Context) SessionRelease() {
	if ctx.sessionStart {
		ctx.Session().SessionRelease()
	}
}

func (ctx *Context) SessionSet(name string, value interface{}) {
	ctx.Session().Set(name, value)
}

func (ctx *Context) SessionGet(name string) interface{} {
	return ctx.Session().Get(name)
}

func (ctx *Context) SessionDel(name string) {
	ctx.Session().Delete(name)
}
func (ctx *Context) SessionDestroy() {
	GlobalSessions.SessionDestroy(ctx.ResponseWriter, ctx.Request)
}
