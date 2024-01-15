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
	VersionNum  = 106
	VersionName = "1.1.1"
)

type App struct {
	DisableLogo bool
}

// Init 初始化应用
func (app *App) Init() {

	//触发Boot事件
	event.Trigger("app.boot", nil)

	conf := config.GetAppConfig()

	//版本显示(优先显示配置文件中的版本号)
	versionName := ""
	if conf.Version == "" {
		versionName = fmt.Sprintf("%s (%d)", VersionName, VersionNum)
	} else {
		versionName = fmt.Sprintf("%s (%d)", conf.Version, conf.VersionNum)
	}

	//创建TempDir
	err := os.MkdirAll(conf.Temp, os.ModePerm)
	if err != nil {
		log.Errorf("Mkdir \"%s\" error : %s", conf.Temp, err.Error())
		panic(err)
	}

	//输出Logo信息
	if !app.DisableLogo {
		var logo = "   _____                     _ \n" +
			"  / ____|                   | |  Go       %s\n" +
			" | |  __  ___  _   _ _ __ __| |  App      v%s\n" +
			" | | |_ |/ _ \\| | | | '__/ _` |  Static   %s\n" +
			" | |__| | (_) | |_| | | | (_| |  Temp Dir %s\n" +
			"  \\_____|\\___/ \\__,_|_|  \\__,_|  Log Dir  %s\n" +
			"--------------------------------------------------------\n"

		logFile := config.GetLogConfig().LogFile
		logDirIndex := strings.LastIndex(logFile, "/")
		fmt.Printf(
			logo, runtime.Version(), versionName,
			config.GetHttpConfig().Static, conf.Temp, logFile[:logDirIndex],
		)
	}

	// 触发Init事件
	event.Trigger("app.init", nil)
}

// Run 启动应用
func (app *App) Run() {

	//命令行解析并执行
	cmd.ConsoleParse()

	// 触发Start事件
	event.Trigger("app.start", nil)

	// 守护进程 -等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Info("Shutdown Server ...", log.Skip())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 触发停止事件
	event.Trigger("app.stop", ctx)

	log.Info("Server exiting", log.Skip())
}
