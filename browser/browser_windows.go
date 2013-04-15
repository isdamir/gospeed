// +build windows
//提供浏览器操作方法
package browser
import (
	"os/exec"
	"log"
)
//打开一个浏览器直至关闭
func OpenBrowserSync(url string){
	openBrowser(url)
}
//异步打开浏览器
func OpenBrowserAsync(url string){
	go openBrowser(url)
}
func openBrowser(url string){
	cmd := exec.Command("cmd", "/c","start",url);
	 err := cmd.Start()
	 if err!=nil{
		log.Fatal(err)
	 }
}
