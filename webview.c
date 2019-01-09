#include "webview.h"

wkeWebView createWebWindow(bool isTransparent, int x, int y, int width, int height)
{
    wkeWebView window = wkeCreateWebWindow(isTransparent ? WKE_WINDOW_TYPE_TRANSPARENT : WKE_WINDOW_TYPE_POPUP, NULL, x, y, width, height);
    //设置数据目录
    wkeSetLocalStorageFullPath(window, wlocalstorage);
    wkeSetCookieJarFullPath(window, wcookiejar);
    //初始化网络文件系统
    initNetFS(window);
    //初始化webview事件
    initWebViewEvent(window);
    return window;
}

HWND getWindowHandle(wkeWebView window)
{
    return wkeGetWindowHandle(window);
}

void loadURL(wkeWebView window, char *url)
{
    wkeLoadURL(window, url);
    free(url);
}

void reloadURL(wkeWebView window)
{
    wkeReload(window);
}

void setWindowTitle(wkeWebView window, char *title)
{
    wkeSetWindowTitle(window, title);
    free(title);
}

const char *getWebTitle(wkeWebView window)
{
    return wkeGetTitle(window);
}

void destroyWindow(wkeWebView window)
{
    wkeDestroyWebWindow(window);
}

void WKE_CALL_TYPE onShowDevtoolsCallback(wkeWebView window, void *param)
{
    //设置数据目录
    wkeSetLocalStorageFullPath(window, wlocalstorage);
    wkeSetCookieJarFullPath(window, wcookiejar);
    initNetFS(window);
    wkeSetWindowTitle(window, "调试工具");
    wkeResizeWindow(window, 900, 650);
    wkeMoveToCenter(window);
    wkeLoadURL(window, wkeGetURL(window));
}

void showDevTools(wkeWebView window)
{
    wkeShowDevtools(window, L"http://__devtools__/inspector.html", onShowDevtoolsCallback, NULL);
}