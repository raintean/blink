package main

import (
	"blink"
	"log"
)

func main() {
	//设置调试模式
	blink.SetDebugMode(true)

	//初始化blink模块
	err := blink.InitBlink()
	if err != nil {
		log.Fatal(err)
	}

	//新建view,加载URL
	view := blink.NewWebView(false, 1366, 920)
	view.LoadURL("https://github.com/raintean/blink.git")
	view.SetWindowTitle("Golang GUI Application")
	view.MoveToCenter()
	view.ShowWindow()
	view.ShowDevTools()

	<-make(chan bool)
}
