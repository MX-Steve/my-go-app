package main

import (
	"fmt"

	"github.com/MX-Steve/my-go-app/vip"
	"github.com/polevpn/webview"
)

func run() {
	w := webview.New(800, 600, false, true)
	defer w.Destroy()
	w.SetTitle("KTZ Platform")
	w.SetSize(800, 600, webview.HintNone)
	port := vip.GetIniData("version.port")
	w.Navigate(fmt.Sprintf("http://localhost:%s/login", port))
	w.Run()
}
