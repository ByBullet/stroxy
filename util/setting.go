package util

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

type ProxySetting interface {
	Setting(int, string) bool
	Unsetting() bool
}

type WinProxySetting struct{}

func (win WinProxySetting) Setting(port int, ignore string) bool {
	c := exec.Command(GetFilePath(PathSettingProxy), fmt.Sprintf("127.0.0.1:%d", port), ignore)
	err := c.Run()
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

func (win WinProxySetting) Unsetting() bool {
	c := exec.Command(GetFilePath(PathUnsettingProxy))
	if err := c.Run(); err != nil {
		log.Println(err)
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
