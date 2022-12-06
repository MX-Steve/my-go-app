package main

import (
	"github.com/polevpn/webview"
)

// .\make_version.bat ./ ./version.h
// 		可以生成最后一位版本号和短哈希值
// rsrc.exe -manifest main.manifest -o app.syso -ico three.ico
// go build -ldflags="-H windowsgui" -o app.exe

func main() {

	w := webview.New(800, 600, false, true)
	defer w.Destroy()
	w.SetTitle("KTZ Platform")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://localhost:7715/login")
	// w.Navigate("https://www.baidu.com")
	w.Run()
}
