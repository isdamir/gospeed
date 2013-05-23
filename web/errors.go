package web

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"runtime"
)

var tpl = `
<!DOCTYPE html> 
<html> 
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=UTF-8" />
    <title>gospeed application error</title>
    <style>
        html, body, body * {padding: 0; margin: 0;}
        #header {background:#ffd; border-bottom:solid 2px #A31515; padding: 20px 10px;}
        #header h2{ }
        #footer {border-top:solid 1px #aaa; padding: 5px 10px; font-size: 12px; color:green;}
        #content {padding: 5px;}
        #content .stack b{ font-size: 13px; color: red;}
        #content .stack pre{padding-left: 10px;}
        table {}
        td.t {text-align: right; padding-right: 5px; color: #888;}
    </style> 
    <script type="text/javascript">
    </script>
</head> 
<body>
    <div id="header">
        <h2>{{.AppError}}</h2>
    </div>
    <div id="content">
        <table>
            <tr>
                <td class="t">Request Method: </td><td>{{.RequestMethod}}</td>
            </tr>
            <tr>
                <td class="t">Request URL: </td><td>{{.RequestURL}}</td>
            </tr>
            <tr>
                <td class="t">RemoteAddr: </td><td>{{.RemoteAddr }}</td>
            </tr>
        </table>
        <div class="stack">
            <b>Stack</b>
            <pre>{{.Stack}}</pre>
        </div>
    </div>
    <div id="footer">
        <p>gospeed {{ .SpeedgoVersion }} (gospeed framework)</p>
        <p>golang version: {{.GoVersion}}</p>
    </div>
</body>
</html>        
`

func ShowErr(err interface{}, rw http.ResponseWriter, r *http.Request, Stack string) {
	t, err := template.New("gospeederrortemp").Parse(tpl)
	data := make(map[string]string)
	data["AppError"] = AppConfig.AppName + ":" + fmt.Sprint(err)
	data["RequestMethod"] = r.Method
	data["RequestURL"] = r.RequestURI
	data["RemoteAddr"] = r.RemoteAddr
	data["Stack"] = Stack
	data["SpeedgoVersion"] = VERSION
	data["GoVersion"] = runtime.Version()
	t.Execute(rw, data)
}

var errtpl = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
        <title>Page Not Found</title>
        <style type="text/css">
            * {
                margin:0;
                padding:0;
            }

            body {
                background-color:#EFEFEF;
                font: .9em "Lucida Sans Unicode", "Lucida Grande", sans-serif;
            }

            #wrapper{
                width:600px;
                margin:40px auto 0;
                text-align:center;
                -moz-box-shadow: 5px 5px 10px rgba(0,0,0,0.3);
                -webkit-box-shadow: 5px 5px 10px rgba(0,0,0,0.3);
                box-shadow: 5px 5px 10px rgba(0,0,0,0.3);
            }

            #wrapper h1{
                color:#FFF;
                text-align:center;
                margin-bottom:20px;
            }

            #wrapper a{
                display:block;
                font-size:.9em;
                padding-top:20px;
                color:#FFF;
                text-decoration:none;
                text-align:center;
            }

            #container {
                width:600px;
                padding-bottom:15px;
                background-color:#FFFFFF;
            }

            .navtop{
                height:40px;
                background-color:#24B2EB;
                padding:13px;
            }

            .content {
                padding:10px 10px 25px;
                background: #FFFFFF;
                margin:;
                color:#333;
            }

            a.button{
                color:white;
                padding:15px 20px;
                text-shadow:1px 1px 0 #00A5FF;
                font-weight:bold;
                text-align:center;
                border:1px solid #24B2EB;
                margin:0px 200px;
                clear:both;
                background-color: #24B2EB;
                border-radius:100px;
                -moz-border-radius:100px;
                -webkit-border-radius:100px;
            }

            a.button:hover{
                text-decoration:none;
                background-color: #24B2EB;
            }

        </style>
    </head>
    <body>
        <div id="wrapper">
            <div id="container">
                <div class="navtop">
                    <h1>{{.Title}}</h1>
                </div>
                <div id="content">
                    {{.Content}}
                    <a href="/" title="Home" class="button">Go Home</a><br />

                    <br>power by gospeed {{.SpeedVersion}}
                </div>
            </div>
        </div>
    </body>
</html>
`

var errorMaps map[int]func() []byte

func init() {
	errorMaps = make(map[int]func() []byte)
}
func ErrorCode(rw http.ResponseWriter, code int) {
	if v, ok := errorMaps[code]; ok {
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		rw.WriteHeader(code)
		rw.Write(v())
	}
}

func registerErrorHander() {
	if _, ok := errorMaps[404]; !ok {
		errorMaps[404] = notFound
	}

	if _, ok := errorMaps[401]; !ok {
		errorMaps[401] = unauthorized
	}

	if _, ok := errorMaps[403]; !ok {
		errorMaps[403] = forbidden
	}

	if _, ok := errorMaps[503]; !ok {
		errorMaps[503] = serviceUnavailable
	}

	if _, ok := errorMaps[500]; !ok {
		errorMaps[500] = internalServerError
	}
}

func ErrorCodeReg(code int, f func() []byte) {
	errorMaps[code] = f
}

//404
func notFound() []byte {
	t, _ := template.New("errortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Page Not Found"
	data["Content"] = template.HTML("<br>The Page You have requested flown the coop." +
		"<br>Perhaps you are here because:" +
		"<br><br><ul>" +
		"<br>The page has moved" +
		"<br>The page no longer exists" +
		"<br>You were looking for your puppy and got lost" +
		"<br>You like 404 pages" +
		"</ul>")
	data["SpeedVersion"] = VERSION
	wr := &bytes.Buffer{}
	t.Execute(wr, data)
	return wr.Bytes()
}

//401
func unauthorized() []byte {
	t, _ := template.New("errortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Unauthorized"
	data["Content"] = template.HTML("<br>The Page You have requested can't authorized." +
		"<br>Perhaps you are here because:" +
		"<br><br><ul>" +
		"<br>Check the credentials that you supplied" +
		"<br>Check the address for errors" +
		"</ul>")
	data["SpeedVersion"] = VERSION
	wr := &bytes.Buffer{}
	t.Execute(wr, data)
	return wr.Bytes()
}

//403
func forbidden() []byte {
	t, _ := template.New("errortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Forbidden"
	data["Content"] = template.HTML("<br>The Page You have requested forbidden." +
		"<br>Perhaps you are here because:" +
		"<br><br><ul>" +
		"<br>Your address may be blocked" +
		"<br>The site may be disabled" +
		"<br>You need to log in" +
		"</ul>")
	data["SpeedVersion"] = VERSION
	wr := &bytes.Buffer{}
	t.Execute(wr, data)
	return wr.Bytes()
}

//503
func serviceUnavailable() []byte {
	t, _ := template.New("errortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Service Unavailable"
	data["Content"] = template.HTML("<br>The Page You have requested unavailable." +
		"<br>Perhaps you are here because:" +
		"<br><br><ul>" +
		"<br><br>The page is overloaded" +
		"<br>Please try again later." +
		"</ul>")
	data["SpeedVersion"] = VERSION
	wr := &bytes.Buffer{}
	t.Execute(wr, data)
	return wr.Bytes()
}

//500
func internalServerError() []byte {
	t, _ := template.New("errortemp").Parse(errtpl)
	data := make(map[string]interface{})
	data["Title"] = "Internal Server Error"
	data["Content"] = template.HTML("<br>The Page You have requested has down now." +
		"<br><br><ul>" +
		"<br>simply try again later" +
		"<br>you should report the fault to the website administrator" +
		"</ul>")
	data["SpeedVersion"] = VERSION
	wr := &bytes.Buffer{}
	t.Execute(wr, data)
	return wr.Bytes()
}
