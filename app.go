package gourd

import (
	"context"
	"fmt"
	"github.com/go-gourd/gourd/cmd"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/event"
	"github.com/go-gourd/gourd/log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"time"
)

// 版本信息
const (
	VersionNum  = 13
	VersionName = "0.3.0"
)

type App struct {
	Version     int
	VersionName string
	Conf        *config.AppConfig
	TempDir     string
	DisableLogo bool
}

// Init 初始化应用
func (app *App) Init() {

	app.Conf = config.GetAppConfig()
	app.TempDir = app.Conf.TempDir

	//版本显示
	if app.VersionName != "" {
		app.Version = VersionNum
		app.VersionName = VersionName
	}

	//创建TempDir
	err := os.MkdirAll(app.TempDir, os.ModePerm)
	if err != nil {
		log.Errorf("Mkdir TempDir Err:%s", err.Error())
		panic(err)
	}

	if !app.DisableLogo {
		var logo = "   _____                     _ \n" +
			"  / ____|                   | |  Go       %s\n" +
			" | |  __  ___  _   _ _ __ __| |  App      v%s (%d)\n" +
			" | | |_ |/ _ \\| | | | '__/ _` |  Public   %s\n" +
			" | |__| | (_) | |_| | | | (_| |  Temp Dir %s\n" +
			"  \\_____|\\___/ \\__,_|_|  \\__,_|  Log Dir  %s\n" +
			"--------------------------------------------------------\n"

		logFile := "log.Config.LogFile"
		logDirIndex := strings.LastIndex(logFile, "/")
		fmt.Printf(
			logo, runtime.Version(), app.VersionName, app.Version,
			config.GetHttpConfig().Public, app.TempDir, logFile[:logDirIndex],
		)
	}

	//触发Boot事件
	event.Trigger("_boot", nil)

}

// Run 启动应用
func (app *App) Run() {
	// 触发Init事件
	event.Trigger("_init", nil)

	//命令行解析
	cmd.ConsoleParse()

	// 触发Start事件
	event.Trigger("_start", nil)

	// 守护进程 -等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...", log.Skip())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 触发停止事件
	event.Trigger("_stop", ctx)

	log.Info("Server exiting", log.Skip())
}
