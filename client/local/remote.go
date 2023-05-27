package local

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/ByBullet/stroxy/client/config"
	"github.com/ByBullet/stroxy/logger"

	"go.uber.org/zap"
	"golang.org/x/net/websocket"
)

// 完成与代理服务器的握手,如果链接失败就返回nil
// path：链接服务器的url
// origin：http请求头的url
// target：目标服务器地址，例如【baidu.com:443】
func connectProxyServer(path, origin, target string) *websocket.Conn {
	c, _ := websocket.NewConfig(path, origin)
	c.TlsConfig = &tls.Config{InsecureSkipVerify: config.ProductConfigGroup.SkipVerify}
	conn, err := websocket.DialConfig(c)
	if err != nil {
		return nil
	}
	_, _ = conn.Write([]byte(target))
	buff := make([]byte, 2)
	n, err := conn.Read(buff)
	if err != nil || n != 2 {
		logger.PROD().Error("代理服务器连接错误", zap.Error(err))
		return nil
	}
	if string(buff) != "ok" {
		return nil
	}
	return conn
}

// 测试连接服务器
func connectProxyServerTest(path, origin string) {
	c, _ := websocket.NewConfig(path, origin)
	//设置不跳过客户端对服务器身份验证验证
	c.TlsConfig = &tls.Config{InsecureSkipVerify: true}
	conn, err := websocket.DialConfig(c)
	if err != nil {
		log.Panicln(err)
	}
	_, _ = conn.Write([]byte("connect success\n"))
	b := make([]byte, 1024)
	l, _ := conn.Read(b)
	fmt.Print(string(b[:l]))
	_ = conn.Close()
}
