package util

import (
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
)

// GetCurrentAbPath 获取项目的跟路径
func GetCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()
	if strings.Index(dir, os.TempDir()) == 0 {
		return getCurrentAbPathByCaller()
	}
	if dir[len(dir)-1] != '/' {
		dir += "/"
	}
	return dir[:strings.LastIndex(dir, modName)+len(modName)+1]
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

// fileExists 判断所给路径文件/文件夹是否存在
func fileExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// GetResourcesPath 客户端获取resources文件夹目录
func GetResourcesPath(m string) string {
	rootPath := GetCurrentAbPath()
	resourcesDir := filepath.Join(rootPath, "resources")
	if !fileExists(resourcesDir) {
		resourcesDir = filepath.Join(rootPath, m+"/resources")
	}
	return resourcesDir
}
