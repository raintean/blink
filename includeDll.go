//包含dll，启动时释放
//+build !notincludedll

package blink

import (
	"fmt"
	"github.com/raintean/blink/internal/dll"
	"os"
	"path/filepath"
	"runtime"
)

//获取dll路径
func getDllPath() (string, error) {
	dllPath := filepath.Join(TempPath, "blink_"+runtime.GOARCH+".dll")
	//准备释放dll到临时目录
	err := os.MkdirAll(TempPath, 0644)
	if err != nil {
		return "", err
	}
	data, err := dll.Asset("blink.dll")
	if err != nil {
		return "", err
	}

	err = func() error {
		file, err := os.Create(dllPath)
		defer file.Close()
		if err != nil {
			return fmt.Errorf("无法创建dll文件,err: %s", err)
		}
		n, err := file.Write(data)
		if err != nil {
			return fmt.Errorf("无法写入dll文件,err: %s", err)
		}
		if len(data) != n {
			return fmt.Errorf("写入校验失败")
		}
		return nil
	}()
	if err != nil {
		return "", err
	}
	return dllPath, nil
}
