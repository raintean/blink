package blink

//#include "netfs.h"
import "C"
import (
	"github.com/mushroomsir/mimetypes"
	"io/ioutil"
	"net/http"
	urlLib "net/url"
	"path/filepath"
	"unsafe"
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
