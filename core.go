package gourd

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/core"
	"github.com/go-gourd/gourd/event"
	"github.com/go-gourd/gourd/log"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strconv"
	"strings"
	"time"
)

func (app *App) Create() *gin.Engine {

	if app.Conf.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()

	app.ginEngine = r

	return r
}

func (app *App) Run() {
	//触发Init事件
	event.OnEvent("_init", nil)

	go app.runHttpServer()

	//触发Start事件
	event.OnEvent("_start", nil)

	// 守护进程 -等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...", zap.Skip())

	//触发停止事件
	event.OnEvent("_start", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if e := app.httpServer.Shutdown(ctx); e != nil {
		log.Error("Server Shutdown:" + e.Error())
	}
	log.Info("Server exiting", zap.Skip())
}

func (app *App) runHttpServer() {

	httpConf := config.GetHttpConfig()

	//默认端口
	if httpConf.Port == 0 {
		httpConf.Port = 8080
	}

	listen := httpConf.Host + ":" + strconv.Itoa(int(httpConf.Port))

	log.Info("Started http server. "+listen, zap.Skip())

	app.httpServer = &http.Server{
		Addr:    listen,
		Handler: app.ginEngine,
	}

	// 服务连接
	if err := app.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error(err.Error())
	}
}

func consoleParse() {

	args := os.Args

	if len(args) == 1 {
		filenameWithSuffix := path.Base(strings.Replace(args[0], "\\", "/", -1))
		filenameOnly := strings.TrimSuffix(filenameWithSuffix, path.Ext(filenameWithSuffix))
		fmt.Println(fmt.Sprintf(core.NoCmdHelp, filenameOnly))
		os.Exit(0)
	}

}
