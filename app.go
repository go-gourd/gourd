package gourd

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/event"
	"github.com/go-gourd/gourd/ghttp"
	"github.com/go-gourd/gourd/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"runtime"
	"time"
)

type App struct {
	Version     int
	VersionName string
	Conf        *config.AppConfig
	TempDir     string
}

// Init 初始化应用
func (app *App) Init() {

	app.Version = 4
	app.VersionName = "0.2.2"
	app.Conf = config.GetAppConfig()
	app.TempDir = app.Conf.TempDir

	//创建TempDir
	err := os.MkdirAll(app.TempDir, os.ModePerm)
	if err != nil {
		logger.Errorf("Mkdir TempDir Err:%s", err.Error())
		panic(err)
	}

	//初始化日志工具
	initLogger()

	//触发Boot事件
	event.OnEvent("_boot", nil)

	var logo = "   _____                     _ \n" +
		"  / ____|                   | |  Go       %s\n" +
		" | |  __  ___  _   _ _ __ __| |  Gourd    v%s (%d)\n" +
		" | | |_ |/ _ \\| | | | '__/ _` |  Gin      %s\n" +
		" | |__| | (_) | |_| | | | (_| |  Public   %s\n" +
		"  \\_____|\\___/ \\__,_|_|  \\__,_|  Temp Dir %s\n" +
		"--------------------------------------------------------\n"
	fmt.Printf(
		logo, runtime.Version(), app.VersionName, app.Version, gin.Version,
		config.GetHttpConfig().Public, app.TempDir,
	)

}

// Run 启动应用
func (app *App) Run() {
	// 触发Init事件
	event.OnEvent("_init", nil)

	//命令行解析
	consoleParse()

	// 触发Start事件
	event.OnEvent("_start", nil)

	// 守护进程 -等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutdown Server ...", zap.Skip())

	// 触发停止事件
	event.OnEvent("_stop", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭Http服务
	ghttp.HttpServerShutdown(ctx)

	logger.Info("Server exiting", zap.Skip())
}
