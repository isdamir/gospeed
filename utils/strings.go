package utils

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/russross/blackfriday"
	"html"
	"html/template"
	"io"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type Strings string

func NewString(i interface{}) Strings {
	s := Strings(fmt.Sprintf("%v", i))
	return s
}

// match regexp with string, and return a named group map
// Example:
//   regexp: "(?P<name>[A-Za-z]+)-(?P<age>\\d+)"
//   string: "CGC-30"
//   return: map[string]string{ "name":"CGC", "age":"30" }
func NamedRegexpGroup(str string, reg *regexp.Regexp) (ng map[string]string, matched bool) {
	rst := reg.FindStringSubmatch(str)
	//fmt.Printf("%s => %s => %s\n\n", reg, str, rst)
	if len(rst) < 1 {
		return
	}
	ng = make(map[string]string)
	lenRst := len(rst)
	sn := reg.SubexpNames()
	for k, v := range sn {
		// SubexpNames contain the none named group,
		// so must filter v == ""
		if k == 0 || v == "" {
			continue
		}
		if k+1 > lenRst {
			break
		}
		ng[v] = rst[k]
	}
	matched = true
	return
}

func (s Strings) String() string {
	return string(s)
}

func (s Strings) Md5() string {
	m := md5.New()
	io.WriteString(m, s.String())

	return fmt.Sprintf("%x", m.Sum(nil))
}
func Substr(s string, start, length int) string {
	bt := []rune(s)
	if start < 0 {
		start = 0
	}
	var end int
	if (start + length) > (len(bt) - 1) {
		end = len(bt) - 1
	} else {
		end = start + length
	}
	return string(bt[start:end])
}
func RandomString(len int) (str string) {
	ch := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	sb := bytes.Buffer{}
	rand.Seed(time.Now().Unix())
	for i := 0; i < len; i++ {
		sb.WriteRune(ch[rand.Intn(len)])
	}
	return sb.String()
}

// MarkDown parses a string in MarkDown format and returns HTML. Used by the template parser as "markdown"
func MarkDown(raw string) (output template.HTML) {
	input := []byte(raw)
	bOutput := blackfriday.MarkdownBasic(input)
	output = template.HTML(string(bOutput))
	return
}

// Html2str() returns escaping text convert from html
func Html2str(html string) string {
	src := string(html)

	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)

	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")

	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")

	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")

	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")

	return strings.TrimSpace(src)
}

func Str2html(raw string) template.HTML {
	return template.HTML(raw)
}

func Htmlquote(src string) string {
	return html.EscapeString(src)
}

func Htmlunquote(src string) string {
	//实体符号解释为HTML
	return html.UnescapeString(src)
}

// convert like this: "HelloWorld" to "hello_world"
func (s Strings) SnakeCasedName() string {
	newstr := make([]rune, 0)
	firstTime := true

	for _, chr := range string(s) {
		if isUpper := 'A' <= chr && chr <= 'Z'; isUpper {
			if firstTime == true {
				firstTime = false
			} else {
				newstr = append(newstr, '_')
			}
			chr -= ('A' - 'a')
		}
		newstr = append(newstr, chr)
	}

	return string(newstr)
}

// convert like this: "hello_world" to "HelloWorld"
func (s Strings) TitleCasedName() string {
	newstr := make([]rune, 0)
	upNextChar := true

	for _, chr := range string(s) {
		switch {
		case upNextChar:
			upNextChar = false
			chr -= ('a' - 'A')
		case chr == '_':
			upNextChar = true
			continue
		}

		newstr = append(newstr, chr)
	}

	return string(newstr)
}

func (s Strings) PluralizeString() string {
	str := string(s)
	if strings.HasSuffix(str, "y") {
		str = str[:len(str)-1] + "ie"
	}
	return str + "s"
}
