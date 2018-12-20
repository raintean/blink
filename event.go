package blink

//#include "event.h"
import "C"
import (
	"unsafe"
)

//export goOnWindowDestroyCallback
func goOnWindowDestroyCallback(window C.wkeWebView, param unsafe.Pointer) {
	go func() {
		view := getWebViewByWindow(window)
		view.Emit("destroy", view)
	}()
}

//export goOnDocumentReadyCallback
func goOnDocumentReadyCallback(window C.wkeWebView) {
	go func() {
		view := getWebViewByWindow(window)
		view.Emit("documentReady", view)
	}()
}
