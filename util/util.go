package util

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"stroxy/logger"
	"time"

	"github.com/nightlyone/lockfile"
)

var processLock lockfile.Lockfile

const lockFile = "lockfile"

func LockProcess() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	filePath := filepath.Join(filepath.Dir(ex), lockFile)
	processLock, _ = lockfile.New(filePath)
	err = processLock.TryLock()
	if err != nil {
		logger.PROD().Error("========进程已启动，别重复启动========")
		os.Exit(1)
	}
}

// UnlockProcess 释放进程单例锁
func UnlockProcess() {
	_ = processLock.Unlock()
}

// GetMacAddress 读取mac地址
func GetMacAddress() string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}
	return interfaces[0].HardwareAddr.String()
}

// OpenBrowser 根据操作系统类型使用系统默认浏览器打开主页
func OpenBrowser(url string) error {
	var err error

	switch goos := runtime.GOOS; goos {
	case "darwin":
		err = exec.Command("open", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}

// ThreeExp 三目运算符
// exp 条件表达式
func ThreeExp(exp bool, a, b any) any {
	if exp {
		return a
	}
	return b
}

// DelayTime 获取延迟时间
// address 格式 localhost:22
// 返回时间单位:毫秒ms
func DelayTime(address string) (time.Duration, error) {
	timeout := time.Duration(2 * time.Second)
	t1 := time.Now()
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	return time.Since(t1), nil
}
