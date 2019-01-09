package blink

//#include "interop.h"
import "C"
import (
	"bytes"
	"errors"
	"fmt"
	"github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
	"reflect"
	"strings"
	"text/template"
	"time"
	"unsafe"
)

//js内容模板
var jsTemplate, _ = template.New("js").Parse(`
window.BlinkData={};
{{range $key,$value := .Data}}
window.BlinkData.{{$key}}=JSON.parse("{{$value}}", __date_parser__);
{{end}}
window.BlinkFunc={};
{{range $name,$func := .Func}}
window.BlinkFunc.{{$name}}=__blink_invoker__.bind(null,"{{$name}}");
{{end}}
`)

//js预制内容
var jsContent = `
//时间解析器
window.__date_parser__ = function (key, value) {
	if (typeof value === 'string') {
		let segments = /^(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):(\d{2})\.(\d*)Z$/.exec(value);
		if (segments)
  			return new Date(Date.UTC(+segments[1], +segments[2] - 1, +segments[3], +segments[4], +segments[5], +segments[6], +segments[7]));
	}
	return value;
};

window.__blink_invoker__ =  function (methodName) {
    let args = Array.prototype.slice.call(arguments);
    args.splice(0, 1);
    let callback = null;
    if (args.length > 0)
        if (typeof args[args.length - 1] === 'function') {
            callback = args[args.length - 1];
            args.splice(args.length - 1, 1);
        }
        
    let promise = new Promise(function(resolve, reject) {
        __invokeProxy(JSON.stringify({
        	MethodName: methodName,
        	Params: args.map(JSON.stringify)
        }),function(resultString) {
        	try{
        	    let result = JSON.parse(resultString, __date_parser__);
        	    if(result.Success){
        	        resolve(result.ReturnValue);
        	    }else{
        	        reject(new Error(result.Message));
        	    }
        	}catch (e) {
        	  reject(e);
        	}
        });
    });
    
    //判断最后一个参数
    if(callback){
        promise.then(function(returnValue){
            callback.apply(null, [undefined].concat(returnValue));
        }).catch(function(error){
            callback.apply(null, [error]);
        });
    }else{
        return promise;
    }
};

window.__blink_runjs__ = function(path) {
    try {
    	//获取path值
    	let pathSegments = path.split(".");
    	let value = window;
    	for (let pathSegment of pathSegments){
    	    //值存在
    	    if(value[pathSegment] !== undefined){
    	        value = value[pathSegment];
    	    }else
    	        throw new Error("指定的值/函数不存在:" + path);
    	}
    	
    	if (value === window) {
    	    throw new Error("不允许获取window的值,请设置path");
    	}
    	
    	if(typeof value === "function"){
    	    //如果是一个函数,则调用他
    	    let args = Array.prototype.slice.call(arguments);
    		args.splice(0, 1);
    		
    		let result = value.apply(null, args.map(it => JSON.parse(it, __date_parser__)))
    		return JSON.stringify({
    			Success: true,
    			ReturnValue: JSON.stringify(result)
    		})
    	} else {
    	    //如果是一个普通的值,将其格式化并传回
    	    return JSON.stringify({
    	    	Success: true,
    	    	ReturnValue: JSON.stringify(value)
    	    })
    	}
    }catch (e) {
    	return JSON.stringify({
    		Success: false,
    		Message: e.message
    	})
    }
}
`

func init() {
	//注册时间的json编码器和解码器,方便和js做时间交互
	jsoniter.RegisterTypeDecoderFunc("time.Time", func(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
		if !iter.ReadNil() {
			token := iter.Read()
			if reflect2.TypeOf(token).String() == "string" {
				timeValue, _ := time.Parse("2006-01-02T15:04:05.000Z07:00", token.(string))
				*((*time.Time)(ptr)) = timeValue
			}
		}
	})

	jsoniter.RegisterTypeEncoderFunc("time.Time", func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
		stream.WriteString((*(*time.Time)(ptr)).UTC().Format("2006-01-02T15:04:05.000Z07:00"))
	}, func(pointer unsafe.Pointer) bool { return false })
}

//export goInvokeDispatcher
func goInvokeDispatcher(window C.wkeWebView, callback C.jsValue, invocationString *C.char) {
	//获取调用的基本信息
	//对invocationString的取值放在go func的外部,因为go func是异步的
	//防止方法返回后MB回收这个字符串内存
	invocationJson := jsoniter.Get([]byte(C.GoString(invocationString)))
	methodName := invocationJson.Get("MethodName").ToString()

	go func() (returnValue []interface{}, err error) {
		//defer函数,用于处理返回值
		defer func() {
			result := struct {
				Success     bool
				Message     string
				ReturnValue []interface{}
			}{}

			//检查是否存在未捕获的异常
			if e := recover(); e != nil {
				result.Success = false
				result.Message = fmt.Sprint(e)
			} else {
				if err != nil {
					//如果有错误
					result.Success = false
					result.Message = err.Error()
				} else {
					result.Success = true
					result.ReturnValue = returnValue
				}
			}

			//返回值
			jsonData, _ := jsoniter.Marshal(result)
			jobQueue <- func() {
				view := getWebViewByWindow(window)
				if view != nil {
					if !view.IsDestroy {
						C.callbackProxy(window, callback, C.CString(string(jsonData)))
					}
				}
			}
		}()

		//拿到对应的view
		view := getWebViewByWindow(window)

		//查找对应的方法
		if _, exist := view.jsFunc[methodName]; !exist {
			return nil, fmt.Errorf("找不到方法%s", methodName)
		}
		method := view.jsFunc[methodName]
		v := reflect.ValueOf(method)
		t := v.Type()

		//参数集合
		var params = make([]reflect.Value, t.NumIn())

		//挨个格式化参数
		for index := 0; index < t.NumIn(); index++ {
			paramType := t.In(index)
			var paramValue reflect.Value

			//判断参数是引用还是值,取到正确的类型
			if paramType.Kind() == reflect.Ptr {
				paramValue = reflect.New(paramType.Elem())
			} else {
				paramValue = reflect.New(paramType)
			}

			//拿到参数对应的值,并通过json格式化出来
			paramJsonString := invocationJson.Get("Params", index).ToString()
			jsoniter.UnmarshalFromString(paramJsonString, paramValue.Interface())

			//判断参数是引用还是值
			if paramType.Kind() == reflect.Ptr {
				params[index] = paramValue
			} else {
				params[index] = paramValue.Elem()
			}
		}

		//调用,并获取结果
		invokeResult := v.Call(params)

		//真实的返回值,不是reflect.Value
		var realReturnValues []interface{}

		if len(invokeResult) > 0 {
			//判断最后返回值是不error
			lastReturnValue := invokeResult[len(invokeResult)-1]
			t := lastReturnValue.Type()
			//Error方法存不存在
			if m, exist := t.MethodByName("Error"); exist {
				//返回值是不是一个
				if m.Type.NumOut() == 1 {
					//返回值的类型是不是String
					if m.Type.Out(0).Kind() == reflect.String {
						//是Error无疑了
						//判断error是否为空
						if lastReturnValue.IsNil() {
							realReturnValues = make([]interface{}, len(invokeResult)-1)
							for index, value := range invokeResult[0 : len(invokeResult)-1] {
								realReturnValues[index] = value.Interface()
							}
							return realReturnValues, nil
						} else {
							//不为空
							return nil, lastReturnValue.Interface().(error)
						}
					}
				}
			}

			//不是Error,所有的返回值都需要处理
			realReturnValues = make([]interface{}, len(invokeResult))
			for index, value := range invokeResult {
				realReturnValues[index] = value.Interface()
			}
			return realReturnValues, nil
		} else {
			//没有返回值
			return nil, nil
		}
	}()
}

//export goGetInteropJS
func goGetInteropJS(window C.wkeWebView) *C.char {
	view := getWebViewByWindow(window)

	var buffer bytes.Buffer
	buffer.WriteString(jsContent)

	jsTemplate.Execute(&buffer, struct {
		Data map[string]string
		Func map[string]interface{}
	}{
		Data: view.jsData,
		Func: view.jsFunc,
	})

	return C.CString(buffer.String())
}

//注入一个方法或者数据
func (view *WebView) Inject(key string, value interface{}) {
	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Invalid || t.Kind() == reflect.Chan || t.Kind() == reflect.UnsafePointer {
		logger.Println("注入错误,不支持类型:", t.Kind())
	}

	if t.Kind() == reflect.Func {
		//如果是一个函数
		if _, exist := view.jsFunc[key]; exist {
			logger.Printf("%s方法已经存在,请检查重复\n", key)
		}
		view.jsFunc[key] = value
	} else {
		//如果是一个值
		if _, exist := view.jsData[key]; exist {
			logger.Printf("%s数据已经存在,请检查重复\n", key)
		}
		//将其变成json字符串
		jsonData, err := jsoniter.Marshal(value)
		if err != nil {
			return
		}
		view.jsData[key] = template.JSEscapeString(string(jsonData))
	}
}

//调用js中方法 or 获取js中的值
func (view *WebView) Invoke(path string, args ...interface{}) (returnValue jsoniter.Any, err error) {
	if view.IsDestroy {
		return nil, errors.New("WebView已经被销毁")
	}

	//所有的调用必须等待文档ready,且没有webview没有destroy
	select {
	case <-view.DocumentReady:
		break
	case <-view.Destroy:
		//view已经destroy
		return nil, errors.New("WebView已经被销毁")
	}

	defer func() {
		//处理未捕获的异常
		if e := recover(); e != nil {
			returnValue = nil
			err = errors.New(fmt.Sprint(e))
		}
	}()

	done := make(chan string)
	jobQueue <- func() {
		if len(args) > 0 {
			paramJsonStrings := make([]string, len(args))
			for index, value := range args {
				json, err := jsoniter.Marshal(value)
				if err != nil {
					paramJsonStrings[index] = "null"
					continue
				}
				paramJsonStrings[index] = `"` + template.JSEscapeString(string(json)) + `"`
			}
			argsList := strings.Join(paramJsonStrings, ",")
			result := C.runJSProxy(view.window, C.CString(fmt.Sprintf(`return __blink_runjs__('%s', %s);`, path, argsList)))
			done <- C.GoString(result)
		} else {
			result := C.runJSProxy(view.window, C.CString(fmt.Sprintf(`return __blink_runjs__('%s');`, path)))
			done <- C.GoString(result)
		}
	}
	resultJson := jsoniter.Get([]byte(<-done))

	if resultJson.Get("Success").ToBool() {
		return jsoniter.Get([]byte(resultJson.Get("ReturnValue").ToString())), nil
	} else {
		return nil, errors.New(resultJson.Get("Message").ToString())
	}
}
