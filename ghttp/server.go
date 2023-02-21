package ghttp

import (
	"context"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/logger"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

var httpServer *http.Server

// RunHttpServer 启动Http监听服务
func RunHttpServer() {

	httpConf := config.GetHttpConfig()

	//默认端口
	if httpConf.Port == 0 {
		httpConf.Port = 8080
	}

	listen := httpConf.Host + ":" + strconv.Itoa(int(httpConf.Port))

	logger.Info("Started http server. "+listen, zap.Skip())

	httpServer = &http.Server{
		Addr:    listen,
		Handler: GetEngine(),
	}

	// 服务连接
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err.Error())
		panic(err)
	}
}

// GetHttpServer 获取http.Server实例
func GetHttpServer() *http.Server {
	return httpServer
}

func HttpServerShutdown(ctx context.Context) {
	if httpServer != nil {
		if e := httpServer.Shutdown(ctx); e != nil {
			logger.Error("HttpServer Shutdown:" + e.Error())
		}
	}
}
