package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/ByBullet/stroxy/logger"
	"github.com/ByBullet/stroxy/util"

	"go.uber.org/zap"
)

type ServerNode struct {
	Description   string
	DomainPrefix  string
	KeyFile       string
	PublicCrtFile string
	MaxProcess    int
}

/*
解析和处理配置文件
*/
type ConfigGroup struct {
	ServerPort        int
	ServerNodes       map[string]ServerNode
	DefaultServerNode string //默认选中的节点
	CurNode           ServerNode
}

var productConfigGroup *ConfigGroup

/*
Init 初始化config模块
为了适配开发环境和发布环境，配置文件可以放在两个位置，
分别是 .../stroxy/server/resources/config.json
.../stroxy/resources/config.json
*/
func Init() {
	configFilePath := filepath.Join(util.GetResourcesPath("server"), "config.json")
	var configFile *os.File
	var err error
	if configFile, err = os.Open(configFilePath); err != nil {
		logger.PROD().Error("配置文件读取异常", zap.Error(err))
		return
	}
	defer configFile.Close()

	productConfigGroup = new(ConfigGroup)
	decode := json.NewDecoder(configFile)
	decode.Decode(productConfigGroup)

	//读取命令行参数
	if len(os.Args) > 1 {
		if _, ok := productConfigGroup.ServerNodes[os.Args[1]]; ok {
			productConfigGroup.DefaultServerNode = os.Args[1]
		}
	}

	productConfigGroup.CurNode = productConfigGroup.ServerNodes[productConfigGroup.DefaultServerNode]
	configContent, _ := json.MarshalIndent(productConfigGroup, "", "  ")
	logger.PROD().Sugar().Infof("配置文件加载完成初始化完成 %s", string(configContent))

}

/*
读取当前配置
*/
func CONF() *ConfigGroup {
	return productConfigGroup
}
