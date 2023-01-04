package main

import (
	"fmt"

	"github.com/MX-Steve/my-go-app/vip"
	"github.com/polevpn/webview"
)

func run() {
	w := webview.New(1000, 800, false, true)
	defer w.Destroy()
	w.SetTitle("KTZ Platform")
	w.SetSize(1000, 800, webview.HintNone)
	port := vip.GetIniData("version.port")
	w.Navigate(fmt.Sprintf("http://localhost:%s/login", port))
	w.Run()
}
