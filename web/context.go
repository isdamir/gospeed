package web

import (
	"fmt"
	"iyf.cc/gospeed/browser"
	"iyf.cc/gospeed/web/session"
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

func (ctx *Context) SetHeader(hdr string, val string, unique bool) {
	if unique {
		ctx.ResponseWriter.Header().Set(hdr, val)
	} else {
		ctx.ResponseWriter.Header().Add(hdr, val)
	}
}
func (ctx *Context) EnUrl(ul string) string {
	if AppConfig.SessionToUrl {
		if _, ok := ctx.Session().Map()["__ToUrl"]; ok {
			if strings.Index(ul, fmt.Sprint(AppConfig.SessionName, "=")) == -1 {
				vs := ctx.Session().SessionID()
				if strings.Index(ul, "?") == -1 {
					return fmt.Sprintf("%s?%s=%s", ul, AppConfig.SessionName, url.QueryEscape(vs))
				}
				if strings.HasSuffix(ul, "&") {
					return fmt.Sprintf("%s=%s", ul, AppConfig.SessionName, url.QueryEscape(vs))
				} else {
					return fmt.Sprintf("&%s=%s", ul, AppConfig.SessionName, url.QueryEscape(vs))
				}
			}
		}
	}
	return ul
}
func EnUrl(ul string) string {
	return ul
}

/*
cookie
cookie[0] => name string
cookie[1] => value string
cookie[2] => expires string
cookie[3] => path string
cookie[4] => domain string
cookie[5] => httpOnly bool
cookie[6] => secure bool
*/
func (ctx *Context) SetCookie(args ...interface{}) *http.Cookie {
	if len(args) < 2 {
		return nil
	}

	const LEN = 7
	var cookie = [LEN]interface{}{}

	for k, v := range args {
		if k >= LEN {
			break
		}

		cookie[k] = v
	}

	var (
		name     string
		value    string
		expires  int
		path     string
		domain   string
		httpOnly bool
		secure   bool
	)

	if v, ok := cookie[0].(string); ok {
		name = v
	} else {
		return nil
	}

	if v, ok := cookie[1].(string); ok {
		value = v
	} else {
		return nil
	}

	if v, ok := cookie[2].(int); ok {
		expires = v
	}

	if v, ok := cookie[3].(string); ok {
		path = v
	}

	if v, ok := cookie[4].(string); ok {
		domain = v
	}

	if v, ok := cookie[5].(bool); ok {
		httpOnly = v
	}

	if v, ok := cookie[6].(bool); ok {
		secure = v
	}

	pCookie := &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     path,
		Domain:   domain,
		HttpOnly: httpOnly,
		Secure:   secure,
	}

	if expires != 0 {
		d, _ := time.ParseDuration(strconv.Itoa(expires) + "s")
		pCookie.Expires = time.Now().Add(d)
	}

	http.SetCookie(ctx.ResponseWriter, pCookie)

	return pCookie
}
func (ctx *Context) Params() *url.Values {
	return &ctx.Request.Form
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
