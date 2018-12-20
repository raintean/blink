//blink函数实现,负责和golang做交互,并包装wke调用
#include "blink.h"

typedef void(WKE_CALL_TYPE *FN_wkeInitializeEx)(const wkeSettings *settings);

void initBlink(char *dllpath, char *localstorage, char *cookiejar)
{
    //转换路径字符串类型
    size_t cSize = strlen(dllpath) + 1;
    wchar_t *wdllpath = (wchar_t *)malloc(sizeof(wchar_t) * cSize);
    mbstowcs(wdllpath, dllpath, cSize);

    cSize = strlen(localstorage) + 1;
    wlocalstorage = (wchar_t *)malloc(sizeof(wchar_t) * cSize);
    mbstowcs(wlocalstorage, localstorage, cSize);

    cSize = strlen(cookiejar) + 1;
    wcookiejar = (wchar_t *)malloc(sizeof(wchar_t) * cSize);
    mbstowcs(wcookiejar, cookiejar, cSize);

    //加载dll
    HMODULE hMod = LoadLibraryW(wdllpath);
    FN_wkeInitializeEx wkeInitializeExFunc = (FN_wkeInitializeEx)GetProcAddress(hMod, "wkeInitializeEx");
    wkeInitializeExFunc((wkeSettings *)0);
    WKE_FOR_EACH_DEFINE_FUNCTION(WKE_GET_PTR_ITERATOR0, WKE_GET_PTR_ITERATOR1, WKE_GET_PTR_ITERATOR2, WKE_GET_PTR_ITERATOR3,
                                 WKE_GET_PTR_ITERATOR4, WKE_GET_PTR_ITERATOR5, WKE_GET_PTR_ITERATOR6, WKE_GET_PTR_ITERATOR11);

    //初始化全局事件
    initGlobalEvent();

    //初始化JS与Golang的交互功能
    initInterop();

    //释放内存
    free(wdllpath);
    free(dllpath);
    free(localstorage);
    free(cookiejar);
}