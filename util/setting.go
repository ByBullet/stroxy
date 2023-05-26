package util

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
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
