//网络文件系统
//拦截webview所有的请求,实现特定域的伪文件系统
#include "netfs.h"
#include "export.h"

//初始化指定webview的网络文件系统
void initNetFS(wkeWebView window)
{
    wkeOnLoadUrlBegin(window, handleLoadUrlBegin, NULL);
    wkeOnLoadUrlEnd(window, handleLoadUrlEnd, NULL);
}

//url加载开始,回调
bool handleLoadUrlBegin(wkeWebView window, void *param, const char *url, wkeNetJob job)
{
    //从golang获取网络文件系统数据
    struct goGetNetFSData_Return returnValue = goGetNetFSData(window, url);
    if (returnValue.result == 1)
    {
        //返回1,表示网络文件系统不处理
        return false;
    }

    if (returnValue.result == 0)
    {
        //设置mimetype
        wkeNetSetMIMEType(job, returnValue.mineType);
        free(returnValue.mineType);
        //设置返回的数据
        wkeNetSetData(job, returnValue.data, returnValue.length);
        free(returnValue.data);
        return true;
    }
    else
    {
        //TODO:暂时返回不处理,交由上层,因为不知道怎么返回404
        return false;
    }
}

//url加载完毕,回调
void handleLoadUrlEnd(wkeWebView window, void *param, const char *url, wkeNetJob job, void *buf, int len)
{
}