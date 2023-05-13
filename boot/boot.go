package boot

import (
	"go.uber.org/zap"
	"os"
	"strings"
	"stroxy/config"
	"stroxy/env"
	"stroxy/local"
	"stroxy/logger"
	"stroxy/util"
)

var localServer *local.Listener

func init() {
	util.LockProcess()
	env.LoadEnv()
	//初始化logger设置
	logger.Init()
	//初始化配置文件
	config.Init()
	//初始化Ip范围检测
	local.InitLimit()
	//创建本本地端口监听对象
	localServer = local.NewLocalListener(config.ProductConfigGroup.LocalPort)
}

// RunProxy
// 启动代理
func RunProxy() {
	localServer.Listen()
	build := strings.Builder{}
	for _, value := range config.ProductConfigGroup.IgnoreAddress {
		build.WriteString(value)
		build.WriteString("@")
	}
	util.SettingProxy(config.ProductConfigGroup.LocalPort, build.String())
	logger.PROD().Debug("启动代理服务器监听")
}

// StopProxy
// 关闭代理
func StopProxy() {
	util.UnsettingProxy()
	localServer.Stop()
	logger.PROD().Debug("关闭代理服务器监听")
}

// ExitSystem
// 退出程序
func ExitSystem() {
	util.UnsettingProxy()
	logger.PROD().Debug("程序退出")
	util.UnlockProcess()
	os.Exit(0)
}

func SelectProxyMode(mode string) {
	switch mode {
	case "auto":
		localServer.Pac = true
		logger.PROD().Debug("代理模式切换为智能代理")
	case "all":
		localServer.Pac = false
		logger.PROD().Debug("代理模式切换为全局代理")
	}
}

// SelectNode
// 选择节点
func SelectNode(nodeName string) {
	config.ProductConfigGroup.SetConfigArg(config.ArgDefaultServer, nodeName)
	logger.PROD().Debug("切换节点", zap.String("节点名称", nodeName))
}
