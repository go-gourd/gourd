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
	VersionNum  = 100
	VersionName = "1.0.0"
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

	//触发Boot事件
	event.Trigger("_boot", nil)

	app.Conf = config.GetAppConfig()

	//临时目录
	app.TempDir = app.Conf.TempDir
	if app.TempDir == "" {
		app.TempDir = "./runtime"
	}

	//版本号
	app.Version = VersionNum
	app.VersionName = VersionName

	//版本显示(优先显示配置文件中的版本号)
	if app.Conf.Version == "" {
		app.Conf.Version = fmt.Sprintf("%s (%d)", app.VersionName, app.Version)
	}

	//创建TempDir
	err := os.MkdirAll(app.TempDir, os.ModePerm)
	if err != nil {
		log.Errorf("Mkdir \"%s\" error : %s", app.TempDir, err.Error())
		panic(err)
	}

	//输出Logo信息
	if !app.DisableLogo {
		var logo = "   _____                     _ \n" +
			"  / ____|                   | |  Go       %s\n" +
			" | |  __  ___  _   _ _ __ __| |  App      v%s (%d)\n" +
			" | | |_ |/ _ \\| | | | '__/ _` |  Public   %s\n" +
			" | |__| | (_) | |_| | | | (_| |  Temp Dir %s\n" +
			"  \\_____|\\___/ \\__,_|_|  \\__,_|  Log Dir  %s\n" +
			"--------------------------------------------------------\n"

		logFile := config.GetLogConfig().LogFile
		logDirIndex := strings.LastIndex(logFile, "/")
		fmt.Printf(
			logo, runtime.Version(), app.VersionName, app.Version,
			config.GetHttpConfig().Public, app.TempDir, logFile[:logDirIndex],
		)
	}

	// 触发Init事件
	event.Trigger("_init", nil)

}

// Run 启动应用
func (app *App) Run() {

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
