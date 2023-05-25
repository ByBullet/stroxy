package util

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"stroxy/logger"

	"github.com/nightlyone/lockfile"
	"go.uber.org/zap"
)

// SettingProxy 设置代理
// port:代理的的口号
// ignore: 不走代理的地址
func SettingProxy(port int, ignore string) bool {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	switch runtime.GOOS {
	case "windows":
		c := exec.Command(GetFilePath(PathSettingProxy), addr, ignore)
		err := c.Run()
		if err != nil {
			log.Println(err)
			return false
		}
	case "darwin":
		err := onMacProxy("127.0.0.1", strconv.Itoa(port))
		if err != nil {
			log.Println(err)
			return false
		}
	//TODO: linux
	case "linux":
	default:
		logger.PROD().Sugar().Errorf("unsupported platform: %s", runtime.GOOS)
		return false
	}
	logger.PROD().Info("成功设置系统代理", zap.String("代理地址", addr))
	return true
}

// UnsettingProxy
// 取消系统代理设置
func UnsettingProxy() bool {
	switch runtime.GOOS {
	case "windows":
		c := exec.Command(GetFilePath(PathUnsettingProxy))
		if err := c.Run(); err != nil {
			log.Println(err)
			return false
		}
	case "darwin":
		err := offMacProxy()
		if err != nil {
			log.Println(err)
			return false
		}
	case "linux":
	default:

	}
	logger.PROD().Info("成功取消系统代理")
	return true
}

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
