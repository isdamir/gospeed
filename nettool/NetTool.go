package nettool

import (
	"code.google.com/p/mahonia"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"time"
)

type NetTool struct {
	ProxyIP    string
	Sleep      time.Duration
	Getheader  *http.Header
	Postheader *http.Header
}
type Response struct {
	resp *http.Response
	err  error
}

var timeout = 10 * time.Second

func (this *NetTool) Get(url string) (str string, err error) {
	str, _, err = this.Do(url)
	return
}
func (this *NetTool) Do(st ...string) (str string, u string, err error) {
	time.Sleep(this.Sleep)
	client := &http.Client{}
	client.Jar = NewJar()
	var oUrl string
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		oUrl = req.URL.String()
		req.Header = this.GetGETHeader()
		if len(via) >= 10 {
			return errors.New("stopped after 10 redirects")
		}
		return nil
	}
	this.SetProxy(client)
	var req *http.Request
	if len(st) == 1 {
		req, err = http.NewRequest("GET", st[0], nil)

		if err != nil {
			return
		}
		req.Header = this.GetGETHeader()
	} else {
		req, err = http.NewRequest("POST", st[0], strings.NewReader(st[1]))

		if err != nil {
			return
		}
		req.Header = this.GetPOSTHeader()
	}
	done := make(chan Response)
	go func() {
		resp, err := client.Do(req)
		done <- Response{resp: resp, err: err}
	}()
	select {
	case resp := <-done:
		if resp.err == nil {
			data, err := ioutil.ReadAll(resp.resp.Body) //取出主体的内容
			if err == nil {
				str = string(data)
				enc := GetCharset(str, resp.resp)
				if enc != "utf-8" {
					dec := mahonia.NewDecoder(enc)
					str = dec.ConvertString(str)
				}
				if oUrl == "" {
					u = st[0]
				} else {
					u = oUrl
				}
			}
		}
	case <-time.After(timeout):
		//timeout
		if trans, ok := http.DefaultTransport.(*http.Transport); ok {
			trans.CancelRequest(req)
		}
		err = errors.New("request timeout")
	}
	return

}
func GetUrl(str string) string {
	ind := strings.LastIndex(str, "/")
	return str[:ind]
}
func GetHost(str string) string {
	ind := strings.Index(str[7:], "/")
	if ind == -1 {
		return str
	}
	return str[:ind+7]
}
func SetUrl(str, old string) string {
	str = strings.Replace(str, "&amp;", "&", -1)
	if !strings.HasPrefix(strings.ToLower(str), "http") {
		if strings.HasPrefix(str, "/") {
			return GetHost(old) + str
		} else {
			return fmt.Sprint(GetUrl(old), "/", str)
		}
	}
	return str
}

func (this *NetTool) Post(url, data string) (str string, err error) {
	str, _, err = this.Do(url, data)
	return
}

func (this *NetTool) GetPOSTHeader() http.Header {
	if this.Postheader == nil {
		this.Postheader = &http.Header{}
		this.Postheader.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.22 (KHTML, like Gecko) Chrome/25.0.1364.160 Safari/537.22")
		this.Postheader.Set("accept-charset", "utf-8;q=0.7,*;q=0.3")
	}
	return *this.Postheader
}
func (this *NetTool) GetGETHeader() http.Header {
	if this.Getheader == nil {
		this.Getheader = &http.Header{}
		this.Getheader.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.22 (KHTML, like Gecko) Chrome/25.0.1364.160 Safari/537.22")
		this.Getheader.Set("accept-charset", "utf-8;q=0.7,*;q=0.3")
	}
	return *this.Getheader
}
func (this *NetTool) SetProxy(ch *http.Client) {
	if this.ProxyIP != "" {
		fixedURL, err := url.Parse(fmt.Sprintf("http://%v", this.ProxyIP))
		if err == nil {
			transport := &http.Transport{Proxy: http.ProxyURL(fixedURL),
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
			transport.ResponseHeaderTimeout = time.Second * 5
			ch.Transport = transport
		}
	}
}

type Jar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	j := new(Jar)
	j.cookies = map[string][]*http.Cookie{}
	return j
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	jar.cookies[u.Host] = cookies
	jar.lk.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

var nChar *regexp.Regexp
var nEChar *regexp.Regexp

func init() {
	nChar, _ = regexp.Compile(`(?is)charset="?(.+?)[";]`)
	nEChar, _ = regexp.Compile(`(?is)charset=(.+?)$`)
}
func GetCharset(str string, r *http.Response) string {
	ms := ""
	if rt := r.Header.Get("content-type"); rt != "" {
		ms = compileChar(rt)
	}
	if ms == "" {
		ms = compileChar(str)
	}
	if ms == "" {
		return "utf-8"
	} else {
		return ms
	}
}
func compileChar(str string) string {
	rt := strings.ToLower(str)
	if strings.Index(rt, "charset") != -1 {
		st := nChar.FindStringSubmatch(rt)
		if len(st) > 1 {
			return st[1]
		}
		st = nEChar.FindStringSubmatch(rt)
		if len(st) > 1 {
			return st[1]
		}
	}
	return ""
}
