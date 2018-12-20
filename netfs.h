//网络文件系统接口头文件
#ifndef NETFS_DEFINE_H
#define NETFS_DEFINE_H

#include "stdio.h"
#include "stdlib.h"
#include "wke.h"

//初始化指定webview的网络文件系统
void initNetFS(wkeWebView window);

//url加载开始,回调
bool handleLoadUrlBegin(wkeWebView window, void *param, const char *url, wkeNetJob job);
//url加载完毕,回调
void handleLoadUrlEnd(wkeWebView window, void *param, const char *url, wkeNetJob job, void *buf, int len);

#endif