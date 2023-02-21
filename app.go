package gourd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/event"
	"github.com/go-gourd/gourd/logger"
	"runtime"
)

type App struct {
	Version     int
	VersionName string
	Conf        *config.AppConfig
	TempDir     string
}

func (app *App) Init() {

	app.Version = 2
	app.VersionName = "0.2.0"
	app.Conf = config.GetAppConfig()
	app.TempDir = app.Conf.TempDir

	//初始化日志工具
	initLogger()

	var logo = "   _____                     _ \n" +
		"  / ____|                   | |  Go       %s\n" +
		" | |  __  ___  _   _ _ __ __| |  Gourd    v%s (%d)\n" +
		" | | |_ |/ _ \\| | | | '__/ _` |  Gin      %s\n" +
		" | |__| | (_) | |_| | | | (_| |  Log Dir  %s\n" +
		"  \\_____|\\___/ \\__,_|_|  \\__,_|  Temp Dir %s\n" +
		"--------------------------------------------------------\n"
	fmt.Printf(
		logo, runtime.Version(), app.VersionName, app.Version, gin.Version,
		config.GetLogConfig().LogFile, app.TempDir,
	)

	//命令行解析
	consoleParse()

	//触发Boot事件
	event.OnEvent("_boot", nil)
}

func initLogger() {

	conf := config.GetLogConfig()

	c := logger.New()
	c.SetDivision("time")     // 设置归档方式，"time"时间归档 "size" 文件大小归档
	c.SetTimeUnit(logger.Day) // 时间归档 可以设置切割单位
	c.SetEncoding("json")     // 输出格式 "json" 或者 "console"

	c.SetInfoFile(conf.LogFile) // 设置info级别日志
	if conf.LogErrorFile != "" {
		c.SetErrorFile(conf.LogErrorFile) // 设置warn级别日志
	}

	c.InitLogger()
}
