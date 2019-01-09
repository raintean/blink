#include "event.h"
#include "export.h"

void initGlobalEvent()
{
}

//当文档ready的时候
void WKE_CALL_TYPE onDocumentReady2Callback(wkeWebView window, void *param, wkeWebFrameHandle frameId)
{
    //只触发main frame 的 ready
    if (wkeWebFrameGetMainFrame(window) == frameId)
    {
        goOnDocumentReadyCallback(window);
    }
}

//当网页标题(title)改变的时候
void WKE_CALL_TYPE onTitleChangedCallback(wkeWebView window, void *param, const wkeString title)
{
    goOnTitleChangedCallback(window, wkeGetString(title));
}

void initWebViewEvent(wkeWebView window)
{
    //窗口被销毁
    wkeOnWindowDestroy(window, goOnWindowDestroyCallback, NULL);
    //JS引擎初始化完毕
    wkeOnDidCreateScriptContext(window, onDidCreateScriptContextCallback, NULL);
    //document ready
    wkeOnDocumentReady2(window, onDocumentReady2Callback, NULL);
    //title changed
    wkeOnTitleChanged(window, onTitleChangedCallback, NULL);
}