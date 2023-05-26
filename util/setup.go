package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"strings"
	"stroxy/logger"

	"go.uber.org/zap"
)

// windows 设置代理脚本内容
const settingProxyBat = `
reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyEnable /t REG_DWORD /d 1 /f
reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyServer /d %s /f
reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyOverride /t REG_SZ /d %s /f
`

// windows 取消代理设置脚本内容
const unSettingProxyBat = `
reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyEnable /t REG_DWORD /d 0 /f 
reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyServer /d "" /f
reg delete "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyOverride /f
`

// onMacProxy 开启Mac端代理
func onMacProxy(host, port string) error {
	services, err := getServices()
	if err != nil {
		return err
	}

	for _, service := range services {
		httpProxy := exec.Command("networksetup", "-setwebproxy", service, host, port)
		httpsProxy := exec.Command("networksetup", "-setsecurewebproxy", service, host, port)

		err = httpProxy.Run()
		if err != nil {
			break
		}
		err = httpsProxy.Run()
		if err != nil {
			break
		}
	}

	if err != nil {
		return err
	}

	return nil
}

// offMacProxy 关闭Mac端代理
func offMacProxy() error {
	services, err := getServices()
	if err != nil {
		return err
	}

	for _, service := range services {
		httpProxy := exec.Command("networksetup", "-setwebproxystate", service, "off")
		httpsProxy := exec.Command("networksetup", "-setsecurewebproxystate", service, "off")

		err = httpProxy.Run()
		if err != nil {
			break
		}
		err = httpsProxy.Run()
		if err != nil {
			break
		}
	}

	if err != nil {
		return err
	}

	return nil
}

// getServices mac端获取网卡服务名称
func getServices() ([]string, error) {
	cmd := exec.Command("networksetup", "-listallnetworkservices")

	res, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	services := strings.Split(string(res), "\n")
	//第一个元素不是网卡名称，最后一个元素为空
	return services[1 : len(services)-1], nil
}

// windows 设置代理
func onWinProxy(port int, ignore string) error {
	content := fmt.Sprintf(settingProxyBat, fmt.Sprintf("127.0.0.1:%d", port), ignore)
	return runWinBatScript(content)
}

// windows 取消代理
func offWinProxy() error {
	return runWinBatScript(unSettingProxyBat)
}

/*
runWinBatScript 执行windows的bat脚本
在临时目录下创建名为stroxy.bat的脚本，然后执行，执行完成后就删除脚本
*/
func runWinBatScript(content string) error {
	s := path.Join(os.TempDir(), "stroxy.bat")
	f, err := os.Create(s)
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(content))
	if err != nil {
		return err
	}

	err = f.Close()
	if err != nil {
		return err
	}

	c := exec.Command(s)
	err = c.Run()
	if err != nil {
		return err
	}

	err = os.Remove(s)
	if err != nil {
		return err
	}
	return nil
}

// SettingProxy 设置代理
// port:代理的的口号
// ignore: 不走代理的地址用分号隔开 127.0.0.1;localhost;192.168.*
func SettingProxy(port int, ignore string) bool {
	addr := fmt.Sprintf("127.0.0.1:%d", port)
	switch runtime.GOOS {
	case "windows":
		err := onWinProxy(port, ignore)
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
		err := offWinProxy()
		if err != nil {
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
		// TODO  linux下只支持 终端的自动代理设置
	}
	logger.PROD().Info("成功取消系统代理")
	return true
}
