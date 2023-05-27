package main

import (
	"io"
	"net"
	"strings"
	"stroxy/logger"
	"time"

	"github.com/juju/ratelimit"
	"go.uber.org/zap"
	"golang.org/x/net/websocket"
)

/*
获取目标host
*/
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

/*
获取请求方法，CONNCET是http是启用隧道协议，别直接转发
*/
func judgeMethod(head string) string {
	split := strings.Split(head, " ")
	return split[0]
}

/*
处理https代理请求
*/
func HttpsHandler(ws *websocket.Conn) {
	defer ws.Close()
	buffer := make([]byte, 1024)
	len, err := ws.Read(buffer)
	buffer = buffer[:len]
	if err != nil {
		ws.Close()
		return
	}
	var defaultPort string
	if judgeMethod(string(buffer)) == "CONNECT" {
		defaultPort = "443"
		dst := getHttpDstHost(string(buffer), defaultPort)
		if conn, err := net.Dial("tcp", dst); err == nil {
			ws.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
			go io.Copy(ws, conn)
			io.Copy(conn, ws)
		}
	} else {
		defaultPort = "80"
		dst := getHttpDstHost(string(buffer), defaultPort)
		if conn, err := net.Dial("tcp", dst); err == nil {
			conn.Write(buffer)
			go io.Copy(ws, conn)
			io.Copy(conn, ws)
		}
	}
}

/*
测试代理服务器是否可用
*/
func serverTest(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

/*
数据转发处理器
数据转发前发目标服务器的host:port,例如【baidu.com:80】
收到此处理器返回的ok表示握手成功
*/
var dataTransferHandler = websocket.Handler(func(c *websocket.Conn) {
	b := make([]byte, 1024)
	n, err := c.Read(b)
	if err != nil {
		return
	}
	logger.PROD().Info("处理请求", zap.String("地址", string(b[:n])))

	targetConn, err := net.Dial("tcp", string(b[:n]))

	if err != nil {
		return
	}
	defer targetConn.Close()
	c.Write([]byte("ok"))
	bucket := ratelimit.NewBucket(time.Second/(1024*1024*1.5), 1024*1024)
	go io.Copy(ratelimit.Writer(c, bucket), targetConn)
	io.Copy(targetConn, c)
})
