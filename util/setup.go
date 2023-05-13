package util

import (
	"os/exec"
	"strings"
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
