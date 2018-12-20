//JS和GO交互的部分
#ifndef INTEROP_DEFINE_H
#define INTEROP_DEFINE_H

#include "stdio.h"
#include "stdlib.h"
#include "wke.h"

void initInterop();

void WKE_CALL_TYPE onDidCreateScriptContextCallback(wkeWebView webView, void *param, wkeWebFrameHandle frameId, void *context, int extensionGroup, int worldId);

jsValue JS_CALL invokeProxy(jsExecState es);

void callbackProxy(wkeWebView window, jsValue callback, char *result);

const char *runJSProxy(wkeWebView window, char *script);

#endif