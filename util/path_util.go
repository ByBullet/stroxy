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
	if strings.Index(dir, getTmpDir()) == 0 {
		return getCurrentAbPathByCaller()
	}
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	return dir[:strings.LastIndex(dir, fmt.Sprintf("/%s/", modName))+len(modName)+1]
}

// getTmpDir 获取系统临时目录，兼容go run
func getTmpDir() string {
	switch runtime.GOOS {
	case "linux":
		return "/tmp"
	case "windows":
		dir := os.Getenv("TEMP")
		if dir == "" {
			dir = os.Getenv("TMP")
		}
		res, _ := filepath.EvalSymlinks(dir)
		return res
	case "darwin":
		return "/tmp"
	}
	return ""
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

// isTest 是否是在单元测试模式下运行
func isTest() bool {
	for _, v := range os.Args {
		v = strings.ToLower(v)
		if strings.Contains(v, "-test") {
			return true
		}
	}
	return false
}
