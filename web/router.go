package web

import (
	"fmt"
	"iyf.cc/gospeed/browser"
	. "iyf.cc/gospeed/log"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

type controllerInfo struct {
	pattern        string
	regex          *regexp.Regexp
	params         map[int]string
	controllerType reflect.Type
}

type userHandler struct {
	pattern string
	regex   *regexp.Regexp
	params  map[int]string
	h       http.Handler
}

type FilterRegistor func(w http.ResponseWriter, r *http.Request, ctx *Context)

type ControllerRegistor struct {
	routers      []*controllerInfo
	fixrouters   []*controllerInfo
	filters      []FilterRegistor
	userHandlers map[string]*userHandler
}

func NewControllerRegistor() *ControllerRegistor {
	return &ControllerRegistor{routers: make([]*controllerInfo, 0), userHandlers: make(map[string]*userHandler)}
}

func (p *ControllerRegistor) Add(pattern string, c ControllerInterface) {
	parts := strings.Split(pattern, "/")

	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "(.+)"
			//a user may choose to override the defult expression
			// similar to expressjs: ‘/user/:id([0-9]+)’
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[1:index]
				//match /user/:id:int ([0-9]+)
				//match /post/:username:string	([\w]+)
			} else if lindex := strings.LastIndex(part, ":"); lindex != 0 {
				switch part[lindex:] {
				case ":int":
					expr = "([0-9]+)"
					part = part[1:lindex]
				case ":string":
					expr = `([\w]+)`
					part = part[1:lindex]
				}
			}
			params[j] = part
			parts[i] = expr
			j++
		}
		if strings.HasPrefix(part, "*") {
			expr := "(.+)"
			if part == "*.*" {
				params[j] = "path"
				parts[i] = "([^.]+).([^.]+)"
				j++
				params[j] = "ext"
				j++
			} else {
				parts[i] = expr
				params[j] = "splat"
				j++
			}
		}
	}
	if j == 0 {
		//now create the Route
		t := reflect.Indirect(reflect.ValueOf(c)).Type()
		route := &controllerInfo{}
		route.pattern = pattern
		route.controllerType = t

		p.fixrouters = append(p.fixrouters, route)
	} else { // add regexp routers
		//recreate the url pattern, with parameters replaced
		//by regular expressions. then compile the regex
		pattern = strings.Join(parts, "/")
		regex, regexErr := regexp.Compile(pattern)
		if regexErr != nil {
			//TODO add error handling here to avoid panic
			panic(regexErr)
			return
		}

		//now create the Route
		t := reflect.Indirect(reflect.ValueOf(c)).Type()
		route := &controllerInfo{}
		route.regex = regex
		route.params = params
		route.pattern = pattern
		route.controllerType = t
		p.routers = append(p.routers, route)
	}
}

func (p *ControllerRegistor) AddHandler(pattern string, c http.Handler) {
	parts := strings.Split(pattern, "/")

	j := 0
	params := make(map[int]string)
	for i, part := range parts {
		if strings.HasPrefix(part, ":") {
			expr := "([^/]+)"
			//a user may choose to override the defult expression
			// similar to expressjs: ‘/user/:id([0-9]+)’
			if index := strings.Index(part, "("); index != -1 {
				expr = part[index:]
				part = part[:index]
			}
			params[j] = part
			parts[i] = expr
			j++
		}
	}
	if j == 0 {
		//now create the Route
		uh := &userHandler{}
		uh.pattern = pattern
		uh.h = c
		p.userHandlers[pattern] = uh
	} else { // add regexp routers
		//recreate the url pattern, with parameters replaced
		//by regular expressions. then compile the regex
		pattern = strings.Join(parts, "/")
		regex, regexErr := regexp.Compile(pattern)
		if regexErr != nil {
			//TODO add error handling here to avoid panic
			panic(regexErr)
			return
		}

		//now create the Route
		uh := &userHandler{}
		uh.regex = regex
		uh.params = params
		uh.pattern = pattern
		uh.h = c
		p.userHandlers[pattern] = uh
	}
}

// Filter adds the middleware filter.
func (p *ControllerRegistor) Filter(filter FilterRegistor) {
	p.filters = append(p.filters, filter)
}

// FilterParam adds the middleware filter if the REST URL parameter exists.
func (p *ControllerRegistor) FilterParam(param string, filter FilterRegistor) {
	if !strings.HasPrefix(param, ":") {
		param = ":" + param
	}

	p.Filter(func(w http.ResponseWriter, r *http.Request, ctx *Context) {
		p := r.URL.Query().Get(param)
		if len(p) > 0 {
			filter(w, r, ctx)
		}
	})
}

// FilterPrefixPath adds the middleware filter if the prefix path exists.
func (p *ControllerRegistor) FilterPrefixPath(path string, filter FilterRegistor) {
	p.Filter(func(w http.ResponseWriter, r *http.Request, ctx *Context) {
		if strings.HasPrefix(r.URL.Path, path) {
			filter(w, r, ctx)
		}
	})
}

// AutoRoute
func (p *ControllerRegistor) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			if !AppConfig.RecoverPanic {
				// go back to panic
				panic(err)
			} else {
				var stack string
				Critical("Handler crashed with error", err)
				for i := 1; ; i++ {
					_, file, line, ok := runtime.Caller(i)
					if !ok {
						break
					}
					Critical(file, line)
					if AppConfig.RunMode == "dev" {
						stack = stack + fmt.Sprintln(file, line)
					}
				}
				if AppConfig.RunMode == "dev" {
					ShowErr(err, rw, r, stack)
				}
			}
		}
	}()
	w := &responseWriter{writer: rw}

	var runrouter *controllerInfo
	var findrouter bool

	params := make(map[string]string)

	//static file server
	for prefix, staticDir := range StaticDir {
		if strings.HasPrefix(r.URL.Path, prefix) {
			file := staticDir + r.URL.Path[len(prefix):]
			http.ServeFile(w, r, file)
			w.started = true
			return
		}
	}

	requestPath := r.URL.Path

	//user defined Handler
	for pattern, c := range p.userHandlers {
		if c.regex == nil && pattern == requestPath {
			c.h.ServeHTTP(rw, r)
			return
		} else if c.regex == nil {
			continue
		}

		//check if Route pattern matches url
		if !c.regex.MatchString(requestPath) {
			continue
		}

		//get submatches (params)
		matches := c.regex.FindStringSubmatch(requestPath)

		//double check that the Route matches the URL pattern.
		if len(matches[0]) != len(requestPath) {
			continue
		}

		if len(c.params) > 0 {
			//add url parameters to the query param map
			values := r.URL.Query()
			for i, match := range matches[1:] {
				values.Add(c.params[i], match)
				params[c.params[i]] = match
			}
			//reassemble query params and add to RawQuery
			r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
			//r.URL.RawQuery = url.Values(values).Encode()
		}
		c.h.ServeHTTP(rw, r)
		return
	}

	//first find path from the fixrouters to Improve Performance
	//通过路径匹配,采用显示注册
	for _, route := range p.fixrouters {
		n := len(requestPath)
		//route like "/"
		if n == 1 {
			if requestPath == route.pattern {
				runrouter = route
				findrouter = true
				break
			} else {
				continue
			}
		}

		if (requestPath[n-1] != '/' && route.pattern == requestPath) ||
			(requestPath[n-1] == '/' && len(route.pattern) >= n-1 && requestPath[0:n-1] == route.pattern) {
			runrouter = route
			findrouter = true
			break
		}
	}

	if !findrouter {
		//find a matching Route
		for _, route := range p.routers {

			//check if Route pattern matches url
			if !route.regex.MatchString(requestPath) {
				continue
			}

			//get submatches (params)
			matches := route.regex.FindStringSubmatch(requestPath)

			//double check that the Route matches the URL pattern.
			if len(matches[0]) != len(requestPath) {
				continue
			}

			if len(route.params) > 0 {
				//add url parameters to the query param map
				values := r.URL.Query()
				for i, match := range matches[1:] {
					values.Add(route.params[i], match)
					params[route.params[i]] = match
				}
				//reassemble query params and add to RawQuery
				r.URL.RawQuery = url.Values(values).Encode() + "&" + r.URL.RawQuery
				//r.URL.RawQuery = url.Values(values).Encode()
			}
			runrouter = route
			findrouter = true
			break
		}
	}

	if runrouter != nil {
		b := browser.NewCheck()
		b.Parser(r)
		ct := &Context{ResponseWriter: w, Request: r, UrlParams: params, Browser: b} //创建上下文对象,在filter中也可以处理上下文,相较于来说,filter只是比controller具有优先执行权限,如果在filter中执行后结束了,那么controller不会被执行,这种需求一般用于全局校验/
		//如果没有满足条件的router,也就是找不到要处理的controller,filter也不会执行.
		//execute middleware filters
		for _, filter := range p.filters {
			filter(w, r, ct)
			if w.started {
				return
			}
		}

		//Invoke the request handler
		vc := reflect.New(runrouter.controllerType)

		//call the controller init function
		init := vc.MethodByName("Init")
		in := make([]reflect.Value, 3)
		in[0] = reflect.ValueOf(ct)
		in[1] = reflect.ValueOf(runrouter.controllerType.Name())
		in[2] = reflect.ValueOf(params)
		init.Call(in)
		//call prepare function
		in = make([]reflect.Value, 0)
		method := vc.MethodByName("Prepare")
		method.Call(in)

		//if response has written,yes don't run next
		if !w.started {
			//通过switch来筛选,防止出现恶意的Method
			switch r.Method {
			case "GET":
				method = vc.MethodByName("Get")
			case "POST":
				method = vc.MethodByName("Post")
			case "HEAD":
				method = vc.MethodByName("Head")
			case "DELETE":
				method = vc.MethodByName("Delete")
			case "PUT":
				method = vc.MethodByName("Put")
			case "PATCH":
				method = vc.MethodByName("Patch")
			case "OPTIONS":
				method = vc.MethodByName("Options")
			default:
				{
					ct.WriteString("Method Is Illegal")
					return
				}
			}
			method.Call(in)
			if !w.started {
				if AppConfig.AutoRender {
					method = vc.MethodByName("Render")
					method.Call(in)
				}
				if !w.started {
					method = vc.MethodByName("Finish")
					method.Call(in)
				}
			}
			ct.SessionRelease()
		}
	}

	//if no matches to url, throw a not found exception
	if w.started == false {
		http.NotFound(w, r)
	}
}

//responseWriter is a wrapper for the http.ResponseWriter
//started set to true if response was written to then don't execute other handler
type responseWriter struct {
	writer  http.ResponseWriter
	started bool
	status  int
}

// Header returns the header map that will be sent by WriteHeader.
func (w *responseWriter) Header() http.Header {
	return w.writer.Header()
}

// Write writes the data to the connection as part of an HTTP reply,
// and sets `started` to true
func (w *responseWriter) Write(p []byte) (int, error) {
	w.started = true
	return w.writer.Write(p)
}

// WriteHeader sends an HTTP response header with status code,
// and sets `started` to true
func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.started = true
	w.writer.WriteHeader(code)
}
