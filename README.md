# blink
使用html来编写golang的GUI程序(only windows), 基于[miniblink开源库](https://github.com/weolar/miniblink49)  

## Demo
[Demo项目地址](https://github.com/raintean/blink-demo)

## 特性
---
- [x] 一个可执行文件, miniblink的dll被嵌入库中
- [x] 生成的可执行文件灰常小,15M左右,upx后 12M左右
- [x] 支持无缝golang和浏览器页面js的交互 (Date类型都做了处理), 并支持异步调用golang中的方法(Promise), 支持异常获取.
- [x] 嵌入开发者工具(bdebug构建tags开启)
- [x] 支持虚拟文件系统, 基于golang的http.FileSystem, 意味着go-bindata出的资源可以直接嵌入程序, 无需开启额外的http服务
- [x] 添加了部分简单的接口(最大化,最小化,无任务栏图标等)
- [x] 设置窗口图标(参见icon.go文件)
- [ ] 支持文件拖拽
- [ ] 自定义dll,而不是使用内嵌的dll(防止更新不及时)
- [ ] golang调用js方法时的异步.
- [ ] dll的内存加载, 尝试过基于MemoryModule的方案, 没有成功, 目前是释放dll到临时目录, 再加载.
- [ ] 还有很多...

## 安装
```bash
go get github.com/raintean/blink
```

## 示例
```go
package main

import (
	"github.com/raintean/blink"
	"github.com/elazarl/go-bindata-assetfs"
	"log"
)

func main() {
	//设置调试模式
	blink.SetDebugMode(true)

	//初始化blink模块
	err := blink.InitBlink()
	if err != nil {
		log.Fatal(err)
	}

	//注册虚拟网络文件系统到域名app
	blink.RegisterFileSystem("app", &assetfs.AssetFS{
		Asset:     bin.Asset,
		AssetDir:  bin.AssetDir,
		AssetInfo: bin.AssetInfo,
	})

	//新建view,加载URL
	view := blink.NewWebView(false, 1366, 920)
	//直接加载虚拟文件系统中的网页
	view.LoadURL("http://app/index.html")
	view.SetWindowTitle("Golang GUI Application")
	view.MoveToCenter()
	view.ShowWindow()
	view.ShowDevTools()

	<-make(chan bool)
	
}
```

## golang和js交互
js调用/获取golang中的方法或者值,异常可捕获
> main.go
```golang
//golang注入方法
view.Inject("GetData", func(num int) (int, error) {
	if num > 10 {
		return 0, errors.New("num不能大于10")
	} else {
		return num + 1, nil
	}
})

//golang注入值
view.Inject("Data", "a string")
```
> index.js
```javascript
await BlinkFunc.GetData(10) //-> 11
await BlinkFunc.GetData(11) //-> throw Error("num不能大于10")
BlinkData.Data // -> "a string"
```
golang调用/获取javascript中的方法或者值,异常可捕获(err变量返回)
> index.js
```javascript
window.Foo = new Date();
window.Bar = function (name) {
    return `hello ${name}`;
};
```
> main.go
```golang
value, err := view.Invoke("Foo")
value.ToXXX // -> Time(golang类型)
value, err := view.Invoke("Bar", "blink")
value.ToString() // -> "hello blink"
```
## 注意
- 网页调试工具默认不打包进可执行文件,请启用BuildTags **bdebug**, eg. `go build -tags bdebug`
- 使用本库需依赖cgo编译环境(mingw32)

## ...
再次感谢miniblink项目, 另外如果觉得本项目好用请点个星.  
欢迎PR, > o <
