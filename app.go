package gourd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/event"
	"github.com/go-gourd/gourd/log"
	"net/http"
	"runtime"
)

type App struct {
	Version     int
	VersionName string
	ginEngine   *gin.Engine
	httpServer  *http.Server
	Conf        *config.AppConfig
	TempDir     string
}

func (app *App) Init() {

	app.Version = 2
	app.VersionName = "0.2.0"
	app.Conf = config.GetAppConfig()
	app.TempDir = app.Conf.TempDir

	log.GetLogger()

	var logo = "   _____                     _ \n" +
		"  / ____|                   | |  Go       %s\n" +
		" | |  __  ___  _   _ _ __ __| |  Gourd    v%s (%d)\n" +
		" | | |_ |/ _ \\| | | | '__/ _` |  Gin      %s\n" +
		" | |__| | (_) | |_| | | | (_| |  Log Dir  %s\n" +
		"  \\_____|\\___/ \\__,_|_|  \\__,_|  Temp Dir %s\n" +
		"--------------------------------------------------------\n"
	fmt.Printf(logo, runtime.Version(), app.VersionName, app.Version, gin.Version, log.FilePath, app.TempDir)

	//命令行解析
	consoleParse()

	//触发Boot事件
	event.OnEvent("_boot", nil)
}
