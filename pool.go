package blink

//#include "wke.h"
import "C"
import (
	"github.com/lxn/win"
	"sync"
)

var windowMap = make(map[C.wkeWebView]*WebView)
var handleMap = make(map[win.HWND]*WebView)
var locker sync.RWMutex

func addViewToPool(view *WebView) {
	locker.Lock()
	windowMap[view.window] = view
	handleMap[view.handle] = view
	locker.Unlock()
	//如果webview destroy的话,从池中清除
	view.On("destroy", func(v *WebView) {
		locker.Lock()
		delete(windowMap, v.window)
		delete(handleMap, v.handle)
		locker.Unlock()
	})
}

func getWebViewByWindow(window C.wkeWebView) *WebView {
	locker.RLock()
	if view, exist := windowMap[window]; exist {
		locker.RUnlock()
		return view
	} else {
		locker.RUnlock()
		return nil
	}
}

func getWebViewByHandle(handle win.HWND) *WebView {
	locker.RLock()
	if view, exist := handleMap[handle]; exist {
		locker.RUnlock()
		return view
	} else {
		locker.RUnlock()
		return nil
	}
}
