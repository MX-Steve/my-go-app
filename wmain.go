package main

import (
	"github.com/polevpn/webview"
)

func main() {

	w := webview.New(800, 600, false, true)
	defer w.Destroy()
	w.SetTitle("KTZ Platform")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://localhost:7715/login")
	w.Run()
}
