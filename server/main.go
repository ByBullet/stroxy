package main

import (
	"fmt"
	"net/http"
	"path/filepath"
	"runtime"

	"github.com/ByBullet/stroxy/logger"
	"github.com/ByBullet/stroxy/server/config"
	"github.com/ByBullet/stroxy/util"

	"time"

	"go.uber.org/zap"
	"golang.org/x/net/websocket"
)

func init() {
	logger.Init()
	config.Init()

}

func createController() http.Handler {
	mux := http.NewServeMux()
	//setting static file system
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static", fileServer)

	mux.Handle("/https", websocket.Handler(HttpsHandler))
	mux.Handle("/test", websocket.Handler(serverTest))
	mux.Handle("/transter", dataTransferHandler)
	mux.HandleFunc("/a", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("Hello World!!"))
	})

	return mux
}

func createServer() {
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", config.CONF().ServerPort),
		Handler:        createController(),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	logger.PROD().Info("run server on", zap.Int("port", config.CONF().ServerPort))
	var err error
	if config.CONF().CurNode.TLS {
		contextPath := util.GetResourcesPath("server")
		publicCrtPath := filepath.Join(contextPath, config.CONF().CurNode.PublicCrtFile)
		keyPath := filepath.Join(contextPath, config.CONF().CurNode.KeyFile)
		err = s.ListenAndServeTLS(publicCrtPath, keyPath)
	} else {
		err = s.ListenAndServe()
	}
	logger.PROD().Error("服务器启动失败", zap.Error(err))
}

func main() {
	// 设置线程数量
	runtime.GOMAXPROCS(config.CONF().CurNode.MaxProcess)
	createServer()
}
