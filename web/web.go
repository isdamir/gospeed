package web

import (
	"fmt"
	"iyf.cc/gospeed/config"
	"iyf.cc/gospeed/log"
	"iyf.cc/gospeed/web/session"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"path"
	"runtime"
)

const VERSION = "0.2.0"

type AppConfigData struct {
	AppName      string
	HttpAddr     string //set "" is bind all adapter
	HttpPort     int
	RecoverPanic bool
	AutoRender   bool
	PprofOn      bool
	ViewsPath    string
	LogLevel     int
	RunMode      string //"dev" or "pro"
	//related to session
	SessionOn            bool   // wheather auto start session,default is false
	SessionProvider      string // default session provider  memory mysql redis
	SessionName          string // sessionName cookie's name
	SessionGCMaxLifetime int64  // session's gc maxlifetime
	SessionSavePath      string // session savepath if use mysql/redis/file this set to the connectinfo
	SessionToUrl         bool   //session save to url,in template ,must use enurl "http://example.com"
	UseFcgi              bool
	MaxMemory            int64
	AutoDevice           bool              //自动设备处理
	EnableGzip           bool              //是否启用gzip压缩
	Custom               map[string]string //自定义信息
}

var (
	AppPath        string
	SpeedApp       *App
	GlobalSessions *session.Manager //GlobalSessions
	StaticDir      map[string]string
	AppConfig      *AppConfigData
)

func init() {
	os.Chdir(path.Dir(os.Args[0]))
	SpeedApp = NewApp()
	AppPath, _ = os.Getwd()
	StaticDir = make(map[string]string)
	AppConfig = &AppConfigData{}
	//set default
	AppConfig.AppName = "Test"
	AppConfig.HttpPort = 80
	AppConfig.RecoverPanic = true
	AppConfig.AutoRender = true
	AppConfig.ViewsPath = "views/"
	AppConfig.RunMode = "pro"
	AppConfig.SessionOn = true
	AppConfig.SessionProvider = "memory"
	AppConfig.SessionName = "gospeed"
	AppConfig.SessionGCMaxLifetime = 1440000
	AppConfig.SessionToUrl = true
	AppConfig.MaxMemory = 64
	AppConfig.AutoDevice = true
	AppConfig.EnableGzip = true
	err := config.GetConfig().Register("conf/app.json", "WebApp", AppConfig)
	if err != nil {
		log.Error(err)
	}
	log.SetLevelBind(&AppConfig.LogLevel)
	StaticDir["/static"] = "static"
	log.Trace(*AppConfig, AppPath)
}

type App struct {
	Handlers *ControllerRegistor
}

// New returns a new PatternServeMux.
func NewApp() *App {
	cr := NewControllerRegistor()
	app := &App{Handlers: cr}
	return app
}

func (app *App) Start() {
	addr := fmt.Sprintf("%s:%d", AppConfig.HttpAddr, AppConfig.HttpPort)
	var err error
	if AppConfig.UseFcgi {
		l, e := net.Listen("tcp", addr)
		if e != nil {
			log.SpeedLogger.Fatal("Listen: ", e)
		}
		err = fcgi.Serve(l, app.Handlers)
	} else {
		log.Debug(addr)
		err = http.ListenAndServe(addr, app.Handlers)
	}
	if err != nil {
		log.SpeedLogger.Fatal("ListenAndServe: ", err)
	}
}

func (app *App) Router(path string, c ControllerInterface) *App {
	app.Handlers.Add(path, c)
	return app
}

func (app *App) Filter(filter FilterRegistor) *App {
	app.Handlers.Filter(filter)
	return app
}

func (app *App) FilterParam(param string, filter FilterRegistor) *App {
	app.Handlers.FilterParam(param, filter)
	return app
}

func (app *App) FilterPrefixPath(path string, filter FilterRegistor) *App {
	app.Handlers.FilterPrefixPath(path, filter)
	return app
}

func (app *App) SetViewsPath(path string) *App {
	AppConfig.ViewsPath = path
	return app
}

func (app *App) SetStaticPath(url string, path string) *App {
	StaticDir[url] = path
	return app
}

func (app *App) ErrorLog(ctx *Context) {
	log.SpeedLogger.Printf("[ERR] host: '%s', request: '%s %s', proto: '%s', ua: '%s', remote: '%s'\n", ctx.Request.Host, ctx.Request.Method, ctx.Request.URL.Path, ctx.Request.Proto, ctx.Request.UserAgent(), ctx.Request.RemoteAddr)
}

func (app *App) AccessLog(ctx *Context) {
	log.SpeedLogger.Printf("[ACC] host: '%s', request: '%s %s', proto: '%s', ua: %s'', remote: '%s'\n", ctx.Request.Host, ctx.Request.Method, ctx.Request.URL.Path, ctx.Request.Proto, ctx.Request.UserAgent(), ctx.Request.RemoteAddr)
}

func RegisterRouter(path string, c ControllerInterface) *App {
	SpeedApp.Router(path, c)
	return SpeedApp
}

func RouterHandler(path string, c http.Handler) *App {
	SpeedApp.Handlers.AddHandler(path, c)
	return SpeedApp
}

func Filter(filter FilterRegistor) *App {
	SpeedApp.Filter(filter)
	return SpeedApp
}

func FilterParam(param string, filter FilterRegistor) *App {
	SpeedApp.FilterParam(param, filter)
	return SpeedApp
}

func FilterPrefixPath(path string, filter FilterRegistor) *App {
	SpeedApp.FilterPrefixPath(path, filter)
	return SpeedApp
}

func Start() {
	if AppConfig.PprofOn {
		SpeedApp.Router(`/debug/pprof`, &ProfController{})
		SpeedApp.Router(`/debug/pprof/:pp([\w]+)`, &ProfController{})
	}
	if AppConfig.SessionOn {
		GlobalSessions, _ = session.NewManager(AppConfig.SessionProvider, AppConfig.SessionName, AppConfig.SessionGCMaxLifetime, AppConfig.SessionSavePath, &AppConfig.SessionToUrl)
		go GlobalSessions.GC()
	}
	err := WatchTemplate()
	if err != nil {
		log.Warn(err)
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	registerErrorHander()
	SpeedApp.Start()
}
