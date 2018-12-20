package blink

import (
	"os"
	"path/filepath"
)

//临时目录,用于存放临时文件如:dll,cookie等
var TempPath = filepath.Join(os.TempDir(), "blink")

//是否为调试模式
var DebugMode = true
