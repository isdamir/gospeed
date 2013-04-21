package browser

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type BrowserCheck struct {
	IsWml     bool
	mobile    int
	IsMobile  bool
	IsIe      bool
	IeVersion int
	com       *regexp.Regexp
}

func NewCheck() *BrowserCheck {
	c := BrowserCheck{}
	c.com, _ = regexp.Compile(`MSIE (\d+)\.`)
	return &c
}
func (this *BrowserCheck) Parser(h *http.Request) {
	ct := h.Header.Get("Content-Type")
	if strings.Index(ct, "vnd.wap.xhtml+xml") != -1 {
		this.mobile++
	}
	if strings.Index(ct, "vnd.wap.wml") != -1 {
		if this.mobile == 0 {
			this.IsWml = true
		}
		this.mobile++
	}
	ua := h.UserAgent()
	if b, err := regexp.MatchString("(up.browser|up.link|mmp|symbian|smartphone|midp|wap|phone)", ua); err == nil {
		if b {
			this.mobile++
		}
	}
	m := []string{"w3c ", "acs-", "alav", "alca", "amoi", "audi", "avan", "benq", "bird", "blac",
		"blaz", "brew", "cell", "cldc", "cmd-", "dang", "doco", "eric", "hipt", "inno",
		"ipaq", "java", "jigs", "kddi", "keji", "leno", "lg-c", "lg-d", "lg-g", "lge-",
		"maui", "maxo", "midp", "mits", "mmef", "mobi", "mot-", "moto", "mwbp", "nec-",
		"newt", "noki", "oper", "palm", "pana", "pant", "phil", "play", "port", "prox",
		"qwap", "sage", "sams", "sany", "sch-", "sec-", "send", "seri", "sgh-", "shar",
		"sie-", "siem", "smal", "smar", "sony", "sph-", "symb", "t-mo", "teli", "tim-",
		"tosh", "tsm-", "upg1", "upsi", "vk-v", "voda", "wap-", "wapa", "wapi", "wapp",
		"wapr", "webc", "winw", "winw", "xda", "xda-"}
	for _, v := range m {
		if strings.Index(ua, v) != -1 {
			this.mobile++
			break
		}
	}
	if strings.Index(ua, "MSIE") != -1 {
		this.mobile = 0
		this.IsIe = true
	}
	t := this.com.FindStringSubmatch(ua)
	if len(t) > 1 {
		this.IeVersion, _ = strconv.Atoi(t[1])
	}
	if this.mobile > 0 {
		this.IsMobile = true
	} else {
		this.IsMobile = false
	}
}
