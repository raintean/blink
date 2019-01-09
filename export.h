//golang导出的函数,将在这里定义,供C语言调用
#ifndef EXPORT_DEFINE_H
#define EXPORT_DEFINE_H

/******************netfs*****************/
//goGetNetFSData函数返回值
struct goGetNetFSData_Return
{
    int result;
    char *mineType;
    void *data;
    int length;
};

//获取网络文件系统数据, -> netfs.go
struct goGetNetFSData_Return goGetNetFSData(wkeWebView window, const char *url);
/*****************netfs end**************/

/******************interpo*****************/
//将JS对Golang的调用分发出去
void goInvokeDispatcher(wkeWebView window, jsValue callback, const utf8 *invocationString);
//获取interop js
char *goGetInteropJS(wkeWebView window);
/*****************interpo end**************/

/******************event*****************/
//window关闭时的回调
void goOnWindowDestroyCallback(wkeWebView window, void *param);
//document ready回调
void goOnDocumentReadyCallback(wkeWebView window);
//title changed回调
void goOnTitleChangedCallback(wkeWebView window, const utf8 *title);
/*****************event end**************/
#endif