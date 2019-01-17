package blink

//#include "netfs.h"
import "C"
import (
	"fmt"
	"io/ioutil"
	"net/http"
	urlLib "net/url"
	"os"
	"path/filepath"
	"unsafe"

	"github.com/mushroomsir/mimetypes"
)

var netFileSystems = make(map[string]http.FileSystem)

func RegisterFileSystem(domain string, fs http.FileSystem) {
	netFileSystems[domain] = fs
}

func UnregisterFileSystem(domain string) {
	delete(netFileSystems, domain)
}

//export goGetNetFSData
func goGetNetFSData(window C.wkeWebView, url *C.char) (result C.int, mimeType *C.char, data unsafe.Pointer, length C.int) {
	//解析url
	u, err := urlLib.Parse(C.GoString(url))
	if err != nil {
		return C.int(-1), nil, nil, C.int(0)
	}

	//只响应http
	if u.Scheme != "http" {
		//无需解析
		return C.int(1), nil, nil, C.int(0)
	}

	//响应指定的域名
	if fs, isExist := netFileSystems[u.Host]; isExist {
		//获取文件
		f, err := fs.Open(u.Path)
		if err != nil {
			return C.int(-1), nil, nil, C.int(0)
		}
		defer f.Close()
		//读取二进制数据
		binData, err := ioutil.ReadAll(f)
		if err != nil {
			return C.int(-1), nil, nil, C.int(0)
		}
		//获取mimetype
		mime := mimetypes.Lookup(filepath.Ext(u.Path))
		return C.int(0), C.CString(mime), unsafe.Pointer(C.CString(string(binData))), C.int(len(binData))
	} else {
		//无需解析
		return C.int(1), nil, nil, C.int(0)
	}
}

//可供外部使用
func GetNetFSData(s string) ([]byte, error) {
	//解析url
	u, err := urlLib.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("解析(%s)失败：%s", s, err.Error())
	}

	//只响应http
	if u.Scheme != "http" {
		//无需解析
		return nil, fmt.Errorf("只能响应http：%s", s)
	}

	//响应指定的域名
	if fs, isExist := netFileSystems[u.Host]; isExist {
		//获取文件
		f, err := fs.Open(u.Path)
		if err != nil {
			return nil, fmt.Errorf("打开资源(%s)失败：%s", s, err.Error())
		}
		defer f.Close()
		//读取二进制数据
		buff, err := ioutil.ReadAll(f)
		if err != nil {
			return nil, fmt.Errorf("读取资源(%s)失败：%s", s, err.Error())
		}

		icoPath := filepath.Join(TempPath, "app.icon")
		fd, err := os.OpenFile(icoPath, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			return nil, fmt.Errorf("读取资源(%s)失败：%s", s, err.Error())
		}
		fd.Write(buff)
		fd.Close()

		return buff, nil
	} else {
		return nil, fmt.Errorf("资源(%s)不存在！", s)
	}
}
