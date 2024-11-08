package main

import webview "github.com/webview/webview_go"

func main() {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Todo List")
	w.SetSize(400, 600, webview.HintNone)
	w.Navigate("data:text/html,<html><body><h1>Hello, WebView!</h1></body></html>")
	w.Run()
}
