package gourd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gourd/gourd/core"
	"github.com/go-gourd/gourd/event"
	"io"
	"os"
)

type Application struct {
	Version     int
	VersionName string
	Event       event.GourdEvent
	Engine      *gin.Engine
}

var globalEvent = event.GourdEvent{}

var globalApp = Application{
	Version:     100,
	VersionName: "1.0.0",
	Event:       globalEvent,
}

func createApp() {

	//执行boot事件
	if globalEvent.Boot != nil {
		globalEvent.Boot()
	}

	var logo = core.GetLogo()

	//控制台输出logo
	fmt.Printf(logo, globalApp.VersionName, 8080)

	gin.SetMode(gin.ReleaseMode)

	globalApp.Engine = gin.New()

	//执行Init事件
	if globalEvent.Init != nil {
		globalEvent.Init()
	}

}

func GetServer() *gin.Engine {
	if globalApp.Engine == nil {
		createApp()
	}
	return globalApp.Engine
}

func StartServer() error {

	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	//TODO: 配置获取
	//var cfg config.AppConfig
	//err := config.GetConfig("app", &cfg)
	//if err != nil {
	//	log.Info(err.Error())
	//}
	//
	//fmt.Print("App.Port = ")
	//fmt.Println(cfg.Port)

	app := GetServer()

	//执行Start事件
	if globalEvent.Start != nil {
		globalEvent.Start()
	}

	return app.Run()
}

// RegisterEvent 注册全局事件
func RegisterEvent(name string, callback event.Handler) {

	if name == "boot" {
		// 框架启动前
		globalEvent.Boot = callback
	} else if name == "init" {
		// 框架初始化后
		globalEvent.Init = callback
	} else if name == "start" {
		// 服务启动时
		globalEvent.Start = callback
	}

}
