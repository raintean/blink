package blink

import "C"
import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/lxn/win"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"syscall"
)

//数据hash > icon handle的缓存映射
var iconCache = sync.Map{}

//从二进制数组中加载icon
//TODO:目前是先把ico二进制数据存到本地(common.TempPath),再使用winapi的LoadImage加载图标,因为我没有找到直接从内存中加载ico文件的方法
func LoadIconFromBytes(iconData []byte) (iconHandle win.HANDLE, err error) {
	//计算数据的hash
	bh := md5.Sum(iconData)
	dataHash := hex.EncodeToString(bh[:])

	//先判断缓存里面有没有
	if handle, isExist := iconCache.Load(dataHash); isExist {
		return handle.(win.HANDLE), nil
	}

	//缓存中没有,则释放到本地目录
	iconFilePath := filepath.Join(TempPath, "icon_"+dataHash)
	if _, err := os.Stat(iconFilePath); os.IsNotExist(err) {
		if err := ioutil.WriteFile(iconFilePath, iconData, 0644); err != nil {
			return 0, errors.New("无法创建临时icon文件: " + err.Error())
		}
	}

	//从文件中加载
	handle, err := LoadIconFromFile(iconFilePath)
	if err != nil {
		return 0, err
	}
	//存入缓存
	iconCache.Store(dataHash, handle)
	//返回结果
	return handle, nil
}

//从文件中加载icon
//TODO:目前这个方法只能从ico文件中加载图标,以后添加png/jpg等文件转码(实时)成ico文件的功能,这样就能直接设定常用的文件格式了
func LoadIconFromFile(iconFilePath string) (iconHandle win.HANDLE, err error) {
	iconFilePathW, err := syscall.UTF16PtrFromString(iconFilePath)
	if err != nil {
		return
	}
	iconHandle = win.LoadImage(
		0,
		iconFilePathW,
		win.IMAGE_ICON,
		0,
		0,
		win.LR_LOADFROMFILE,
	)
	if iconHandle == 0 {
		return 0, errors.New("加载图标文件失败," + iconFilePath)
	}
	return
}

//设置窗口图标
func (view *WebView) SetWindowIcon(iconHandle win.HANDLE) error {
	if iconHandle == 0 {
		return errors.New("icon实例非法")
	}
	win.SendMessage(view.handle, win.WM_SETICON, 0, uintptr(iconHandle))
	win.SendMessage(view.handle, win.WM_SETICON, 1, uintptr(iconHandle))
	//TODO:这里没有获取last error了
	return nil
}

//设置窗口图标(从图标文件中). 快捷方法
func (view *WebView) SetWindowIconFromFile(iconFilePath string) error {
	iconHandle, err := LoadIconFromFile(iconFilePath)
	if err != nil {
		return err
	}
	return view.SetWindowIcon(iconHandle)
}

//设置窗口图标(从图标二进制数据中). 快捷方法
func (view *WebView) SetWindowIconFromBytes(iconData []byte) error {
	iconHandle, err := LoadIconFromBytes(iconData)
	if err != nil {
		return err
	}
	return view.SetWindowIcon(iconHandle)
}
