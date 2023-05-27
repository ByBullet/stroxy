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
)

type ProxySetting interface {
	Setting(int, string) bool
	Unsetting() bool
}

type WinProxySetting struct{}

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

func (win WinProxySetting) Setting(port int, ignore string) bool {
	content := fmt.Sprintf(settingProxyBat, fmt.Sprintf("127.0.0.1:%d", port), ignore)
	if err := runWinBatScript(content); err != nil {
		logger.PROD().Sugar().Errorln(err)
		return false
	}
	return true
}

func (win WinProxySetting) Unsetting() bool {
	if err := runWinBatScript(unSettingProxyBat); err != nil {
		logger.PROD().Sugar().Errorln(err)
		return false
	}
	return true
}

type MacProxySetting struct{}

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

func (mac MacProxySetting) Setting(port int, ignore string) bool {
	err := onMacProxy("127.0.0.1", strconv.Itoa(port))
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (mac MacProxySetting) Unsetting() bool {
	err := offMacProxy()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

type LinuxProxySetting struct{}

func (linux LinuxProxySetting) Setting(port int, ignore string) bool {
	return false
}

func (linux LinuxProxySetting) Unsetting() bool {
	return false
}

func getProxySetting(osName string) ProxySetting {
	switch osName {
	case "windows":
		return new(WinProxySetting)
	case "linux":
		return new(LinuxProxySetting)
	case "darwin":
		return new(MacProxySetting)
	}
	return nil
}

// SettingProxy 设置代理
// port:代理的的口号
// ignore: 不走代理的地址用分号隔开 127.0.0.1;localhost;192.168.*
func SettingProxy(port int, ignore string) bool {
	setting := getProxySetting(runtime.GOOS)
	return setting.Setting(port, ignore)
}

// UnsettingProxy
// 取消系统代理设置
func UnsettingProxy() bool {
	setting := getProxySetting(runtime.GOOS)
	return setting.Unsetting()
}
