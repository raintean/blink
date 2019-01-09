//webview接口头文件
#ifndef WEBVIEW_DEFINE_H
#define WEBVIEW_DEFINE_H

#include "stdio.h"
#include "wke.h"
#include "netfs.h"
#include "interop.h"
#include "event.h"

wchar_t *wlocalstorage;
wchar_t *wcookiejar;

wkeWebView createWebWindow(bool isTransparent, int x, int y, int width, int height);
HWND getWindowHandle(wkeWebView window);
void loadURL(wkeWebView window, char *url);
void reloadURL(wkeWebView window);
void setWindowTitle(wkeWebView window, char *title);
const char *getWebTitle(wkeWebView window);
void destroyWindow(wkeWebView window);
void showDevTools(wkeWebView window);

#endif