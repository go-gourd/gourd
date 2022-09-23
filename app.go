package gourd

import (
	"github.com/gin-gonic/gin"
	"github.com/go-gourd/gourd/event"
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

	globalApp.Engine = gin.Default()

	//执行Init事件
	if globalEvent.Init != nil {
		globalEvent.Init()
	}

}

func GetGinEngine() *gin.Engine {
	if globalApp.Engine == nil {
		createApp()
	}
	return globalApp.Engine
}

func StartServer() error {

	//执行Init事件
	if globalEvent.Init != nil {
		globalEvent.Init()
	}

	//执行Init事件
	if globalEvent.Start != nil {
		globalEvent.Start()
	}

	return GetGinEngine().Run()
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
