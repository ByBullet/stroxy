package local

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/ByBullet/stroxy/client/config"
	"github.com/ByBullet/stroxy/logger"

	"go.uber.org/zap"
)

type Listener struct {
	localPort int
	server    net.Listener // 此字段为nil意味着没有开启代理
	Pac       bool         //true：智能代理  ， false： 全局代理
}

func NewLocalListener(port int) (result *Listener) {
	result = new(Listener)
	result.localPort = port
	result.Pac = true
	return
}

// 获取一个未被占用的端口
func getFreePort() int {
	// 创建一个监听器，地址为 ":0"，表示随机分配一个未被占用的端口
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		logger.PROD().Error("操作系统没用未绑定端口", zap.Error(err))
		return 0
	}

	// 获取随机分配的端口
	port := l.Addr().(*net.TCPAddr).Port

	// 关闭监听器
	err = l.Close()
	if err != nil {
		logger.PROD().Error("关闭失败", zap.Error(err))
		return 0
	}

	return port
}

// Listen
// 启动监听
func (l *Listener) Listen() {
	if l.server != nil {
		logger.PROD().Error("代理程序已运行，勿重复启动")
		return
	}
	for l.server == nil {
		conn, err := net.Listen("tcp", fmt.Sprintf(":%d", l.localPort))
		if err != nil {
			logger.PROD().Error("tcp访问监听失败", zap.Int("端口", l.localPort), zap.Error(err))
			newPort := getFreePort()
			l.localPort = newPort
			config.ProductConfigGroup.SetConfigArg(config.ArgPort, strconv.Itoa(newPort))
			continue
		}
		l.server = conn
	}

	go l.listenProcess()
}

func IoCopy(src io.ReadCloser, dst io.WriteCloser) {
	b := make([]byte, 1024*1024/2)
	for {
		n, err := src.Read(b)
		if err != nil {
			_ = src.Close()
			_ = dst.Close()
			return
		}
		_, _ = dst.Write(b[:n])
	}

}

// 数据交换
func dataExchange(localConn, targetConn io.ReadWriteCloser) {
	go IoCopy(targetConn, localConn)
	IoCopy(localConn, targetConn)
}

// 接收要代理的连接请求
func (l *Listener) listenProcess() {
	logger.PROD().Info("代理服务器启动", zap.Int("端口", l.localPort))
	for {
		c, err := l.server.Accept()
		if err != nil {
			logger.PROD().Info("代理服务器监听关闭", zap.Error(err))
			return
		}
		go handleProcess(c, l.Pac)
	}
}

// 处理请求后的握手，权限认证等处理
func handleProcess(c net.Conn, pac bool) {
	buff := make([]byte, 500)
	size, err := c.Read(buff)
	if err != nil {
		return
	}
	buff1 := buff[:size]
	data1 := string(buff1)
	method := strings.Split(data1, " ")[0]
	var targetConn io.ReadWriteCloser
	if strings.ToUpper(method) == "CONNECT" {
		destAddr := getHttpDstHost(data1, "443")
		host := strings.Split(destAddr, ":")[0]
		if pac && localIpLimit.Check(host) { //不走代理，直接连接
			logger.PROD().Debug("本地直连", zap.String("协议", "https"), zap.String("地址", destAddr))
			tempConn, err := net.Dial("tcp", destAddr)
			if err != nil {
				_ = c.Close()
				return
			}
			_, _ = c.Write([]byte("HTTP/1.0 200 Connection established\r\n\r\n"))
			targetConn = tempConn
		} else { //走代理
			logger.PROD().Debug("代理连接", zap.String("协议", "https"), zap.String("地址", destAddr))
			tempConn := connectProxyServer(config.ProductConfigGroup.ProxyUrl, config.ProductConfigGroup.ProxyOrigin, destAddr)
			if tempConn == nil {
				_ = c.Close()
				return
			}
			_, _ = c.Write([]byte("HTTP/1.0 200 Connection established\r\n\r\n"))
			targetConn = tempConn
		}
	} else {
		destAddr := getHttpDstHost(data1, "80")
		host := strings.Split(destAddr, ":")[0]
		if pac && localIpLimit.Check(host) { //不走代理
			logger.PROD().Debug("本地直连", zap.String("协议", "http"), zap.String("地址", destAddr))
			tempConn, err := net.Dial("tcp", destAddr)
			if err != nil {
				logger.PROD().Error("本地直连目标访问失败", zap.String("协议", "http"), zap.String("地址", destAddr), zap.Error(err))

				_ = c.Close()
				return
			}
			_, _ = tempConn.Write(buff1)
			targetConn = tempConn
		} else { //走代理
			logger.PROD().Debug("代理连结", zap.String("协议", "http"), zap.String("地址", destAddr))
			tempConn := connectProxyServer(config.ProductConfigGroup.ProxyUrl, config.ProductConfigGroup.ProxyOrigin, destAddr)
			if tempConn == nil {
				logger.PROD().Error("代理服务器目标访问失败", zap.String("协议", "http"), zap.String("地址", destAddr), zap.Error(err))
				_ = c.Close()
				return
			}
			_, _ = tempConn.Write(buff1)
			targetConn = tempConn
		}
	}
	dataExchange(c, targetConn)
}

// Stop
// 关闭服务器
func (l *Listener) Stop() {
	if l.server == nil {
		return
	}
	_ = l.server.Close()
	l.server = nil
}

// 获取目标host
func getHttpDstHost(head, defPort string) string {
	headLines := strings.Split(head, "\r\n")
	for i, v := range headLines {
		if i > 0 {
			kvs := strings.Split(v, ": ")
			if strings.ToUpper(kvs[0]) == "HOST" {
				if kvs[1] != "" && !strings.Contains(kvs[1], ":") {
					return kvs[1] + ":" + defPort
				}
				return kvs[1]
			}
		}
	}
	return ""
}

// SetLocalPort
// 在服务器关闭的时候修改端口
func (l *Listener) SetLocalPort(port int) {
	if l.server == nil {
		l.localPort = port
	}
}
