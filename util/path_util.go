package util

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const modName = "stroxy"

const (
	PathConfig = "/resources/config.json" // 配置文件路径
	PathIp     = "/resources/ip.txt"      // ip范围文件路径

	PathSettingProxy   = "/script/setting.bat"
	PathUnsettingProxy = "/script/unsetting.bat"
)

// GetFilePath
// 获取所在项目的相对路径
func GetFilePath(relativePath string) string {
	return fmt.Sprintf("%s%s", GetCurrentAbPath(), relativePath)
}

// GetCurrentAbPath 获取项目的跟路径
func GetCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()
	if strings.Index(dir, os.TempDir()) == 0 {
		return getCurrentAbPathByCaller()
	}
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	return dir[:strings.LastIndex(dir, fmt.Sprintf("/%s/", modName))+len(modName)+1]
}

// getCurrentAbPathByExecutable 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// getCurrentAbPathByCaller 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return path.Dir(abPath)
}
