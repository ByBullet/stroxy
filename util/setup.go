package util

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"stroxy/logger"

	"go.uber.org/zap"
)

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
		// TODO  linux下只支持终端的自动代理设置
	}
	logger.PROD().Info("成功取消系统代理")
	return true
}

// SettingProxy2 设置代理
// port:代理的的口号
// ignore: 不走代理的地址用分号隔开 127.0.0.1;localhost;192.168.*
func SettingProxy2(port int, ignore string) bool {
	setting := getProxySetting(runtime.GOOS)
	return setting.Setting(port, ignore)
}

// UnsettingProxy2
// 取消系统代理设置
func UnsettingProxy2() bool {
	setting := getProxySetting(runtime.GOOS)
	return setting.Unsetting()
}
