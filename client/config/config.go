package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	"github.com/ByBullet/stroxy/logger"
	"github.com/ByBullet/stroxy/util"

	"go.uber.org/zap"
)

const (
	// ProxyModelAll or ProxyModelAuto 代理模式
	ProxyModelAll  = "all"
	ProxyModelAuto = "auto"

	ArgPort          = "LocalPort"
	ArgDefaultServer = "DefaultServer"
)

type ProxyServer struct {
	Name, ProxyOrigin, ProxyUrl string
	Port                        int
}

// Group
// 解析和处理配置文件
type Group struct {
	ProxyUrl       string //代理服务器连接url
	ProxyOrigin    string //代理服务器连接头中的origin
	ProxyServerMap map[string]ProxyServer
	DefaultServer  string
	LocalPort      int    //本地监听端口
	ProxyModel     string //代理模式
	SkipVerify     bool   //跳过tls双向认证
	IpDate         string //ip.txt更新时间
	WebsiteAddr    string //官网地址
	Auth           string //身份信息
	IgnoreAddress  []string
}

// 提供配置文件的自动打开和关闭，实际的文件操作通过回调f()实现
func fileOperate(f func(file *os.File)) {
	configFilePath := filepath.Join(util.GetResourcesPath("client"), "config.json")
	if file, err := os.Create(configFilePath); err == nil {
		defer file.Close()
		f(file)
	} else {
		fmt.Printf("err: %v\n", err)
	}
}

func (cg *Group) SetConfigArg(argName, value string) {
	valueOfCG := reflect.ValueOf(cg).Elem()
	v := valueOfCG.FieldByName(argName)
	if !v.IsValid() {
		return
	}
	switch v.Type().Kind() {
	case reflect.String:
		v.SetString(value)
	case reflect.Int:
		i, _ := strconv.Atoi(value)
		v.SetInt(int64(i))
	default:
		logger.PROD().Error("数据类型错误")
	}
	cg.flushToFile()
	if argName == "DefaultServer" {
		cg.buildServerAddr()
	}
}

func (cg *Group) buildServerAddr() {
	defaultServer := cg.ProxyServerMap[cg.DefaultServer]
	cg.ProxyUrl = defaultServer.ProxyUrl
	cg.ProxyOrigin = defaultServer.ProxyOrigin
}

// 把内存中修改的配置写回到文件
func (cg *Group) flushToFile() {
	fileOperate(func(file *os.File) {
		encode := json.NewEncoder(file)
		encode.SetIndent("", "	")
		_ = encode.Encode(*cg)
	})
}

var ProductConfigGroup *Group

/*
Init 初始化config模块
为了适配开发环境和发布环境，配置文件可以放在两个位置，
分别是 .../stroxy/client/resources/config.json和.../stroxy/resources/config.json
第二种情况在发布环境
*/
func Init() {

	var configFile *os.File
	var err error
	configFilePath := filepath.Join(util.GetResourcesPath("client"), "config.json")
	if configFile, err = os.Open(configFilePath); err != nil {
		logger.PROD().Error("配置文件读取异常", zap.Error(err))
		return
	}
	defer configFile.Close()

	ProductConfigGroup = new(Group)
	decode := json.NewDecoder(configFile)
	_ = decode.Decode(ProductConfigGroup)
	ProductConfigGroup.buildServerAddr()
	logger.PROD().Info("config 模块初始化完成")
	configContent, _ := json.MarshalIndent(ProductConfigGroup, "", "  ")
	logger.PROD().Sugar().Debug("配置信息 %s", string(configContent))

}
