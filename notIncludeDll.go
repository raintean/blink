//不包含dll，使用独立dll

//+build notincludedll

package blink

import (
	"path/filepath"
)

//获取dll路径
func getDllPath() (string, error) {
	return filepath.Join(TempPath, "node.dll"), nil
}
