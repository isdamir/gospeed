package browser

import (
	"net/http"
	"strconv"
	"time"
)

/*
设置浏览器cookie
cookie[0] => name string
cookie[1] => value string
cookie[2] => expires string
cookie[3] => path string
cookie[4] => domain string
cookie[5] => httpOnly bool
cookie[6] => secure bool
*/
func SetCookie(w http.ResponseWriter, target map[string]string, args ...interface{}) *http.Cookie {
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

		if target != nil {
			if expires > 0 {
				target[pCookie.Name] = pCookie.Value
			} else {
				delete(target, pCookie.Name)
			}
		}
	}

	http.SetCookie(w, pCookie)

	return pCookie
}
