#include "interop.h"
#include "export.h"

//初始化交互功能
void initInterop()
{
    //绑定一个调用代理函数,一定为2个参数
    jsBindFunction("__invokeProxy", invokeProxy, 2);
}

//当JS引擎初始化完毕,这个时候需要挂载golang的Interop
void WKE_CALL_TYPE onDidCreateScriptContextCallback(wkeWebView window, void *param, wkeWebFrameHandle frameId, void *context, int extensionGroup, int worldId)
{
    //只有main frame才初始化
    if (wkeWebFrameGetMainFrame(window) == frameId)
    {
        char *jsContent = goGetInteropJS(window);
        wkeRunJS(window, jsContent);
        free(jsContent);
    }
}

//调用代理函数,会绑定为JS中的全局函数,所有JS到Golang的函数调用,都将由这个函数来代理
jsValue JS_CALL invokeProxy(jsExecState es)
{
    //检查参数个数
    if (jsArgCount(es) != 2)
    {
        return jsThrowException(es, "调用代理参数个数必须为2");
    }
    else
    {
        //取得第一个参数,为调用详情json
        const utf8 *invocationString = jsToTempString(es, jsArg(es, 0));
        //取得第二个参数,为回调函数,等golang函数完成后,该函数将被调用
        jsValue callbackFunction = jsArg(es, 1);
        //给callback添加一个引用,以免回收
        jsAddRef(es, callbackFunction);
        //调用go的分发器函数,把JS对Golang的调用分发出去
        goInvokeDispatcher(jsGetWebView(es), callbackFunction, invocationString);
        return jsUndefined();
    }
}

//回调代理函数,由Golang来调用,当Golang的分发器函数完成调用后,会调用此函数,把调用结果返回给JS
void callbackProxy(wkeWebView window, jsValue callback, char *result)
{
    //拿到WebView全局的es
    jsExecState es = wkeGlobalExec(window);
    //把调用结果组合成参数列表
    jsValue args[1] = {jsString(es, result)};
    free(result);
    //调用回调函数
    jsCall(es, callback, jsUndefined(), args, 1);
    //回调函数调用完成后,释放引用
    jsReleaseRef(es, callback);
}

//run js代理,供golang来调用
const char *runJSProxy(wkeWebView window, char *script)
{
    jsValue result = wkeRunJS(window, script);
    free(script);
    return jsToTempString(wkeGlobalExec(window), result);
}