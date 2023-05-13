package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"stroxy/logger"
	"stroxy/util"
)

const (
	// ProxyModelAll or ProxyModelAuto 代理模式
	ProxyModelAll  = "all"
	ProxyModelAuto = "auto"

	ArgPort          = "LocalPort"
	ArgDefaultServer = "DefaultServer"
)

type ProxyServer struct {
	Name, DomainPrefix string
	Port               int
}

// Group
// 解析和处理配置文件
type Group struct {
	ProxyServerAddrTemplate string
	ProxyOriginTemplate     string
	ProxyServerAddr         string //代理服务器连接url
	ProxyOrigin             string //代理服务器连接头中的origin
	ProxyServerMap          map[string]ProxyServer
	DefaultServer           string
	LocalPort               int    //本地监听端口
	ProxyModel              string //代理模式
	SkipVerify              bool   //跳过tls双向认证
	IpDate                  string //ip.txt更新时间
	WebsiteAddr             string //官网地址
	Auth                    string //身份信息
	IgnoreAddress           []string
}

// 提供配置文件的自动打开和关闭，实际的文件操作通过回调f()实现
func fileOperate(f func(file *os.File)) {
	if file, err := os.Create(util.GetFilePath(util.PathConfig)); err == nil {
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
	cg.ProxyServerAddr = fmt.Sprintf(cg.ProxyServerAddrTemplate, defaultServer.DomainPrefix, defaultServer.Port)
	cg.ProxyOrigin = fmt.Sprintf(cg.ProxyOriginTemplate, defaultServer.DomainPrefix, defaultServer.Port)
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

func Init() {
	if f, err := os.Open(util.GetFilePath(util.PathConfig)); err == nil {
		defer f.Close()
		ProductConfigGroup = new(Group)
		decode := json.NewDecoder(f)
		_ = decode.Decode(ProductConfigGroup)
		ProductConfigGroup.buildServerAddr()
		logger.PROD().Info("config 模块初始化完成")
		configContent, _ := json.MarshalIndent(ProductConfigGroup, "", "  ")
		logger.PROD().Sugar().Debug("配置信息 %s", string(configContent))
	} else {
		log.Fatalf("err: %v\n", err)
	}
}
