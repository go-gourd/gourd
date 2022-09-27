package gourd

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/core"
	"github.com/go-gourd/gourd/event"
	"github.com/go-gourd/gourd/log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type Application struct {
	Version     int
	VersionName string
	Event       event.GourdEvent
	Engine      *gin.Engine
	Config      config.AppConfig
	Http        *http.Server
}

var globalEvent = event.GourdEvent{}

var globalApp = Application{
	Version:     100,
	VersionName: "1.0.0",
	Event:       globalEvent,
}

// GetServer 获取或者创建一个服务
func GetServer() *gin.Engine {
	if globalApp.Engine == nil {

		//执行boot事件
		if globalEvent.Boot != nil {
			globalEvent.Boot()
		}

		gin.SetMode(globalApp.Config.ReleaseMode)

		globalApp.Engine = gin.New()

		//执行Init事件
		if globalEvent.Init != nil {
			globalEvent.Init()
		}
	}
	return globalApp.Engine
}

// StartServer 启动服务
func StartServer(isDaemon bool) {

	//App配置获取
	var cfg config.AppConfig
	err := config.GetConfig("app", &cfg)
	if err != nil {
		log.Info(err.Error())
	}
	globalApp.Config = cfg

	//默认端口
	if cfg.Port == 0 {
		cfg.Port = 8080
	}

	addr := cfg.Ip + ":" + strconv.Itoa(cfg.Port)

	var logo = core.GetLogo()
	//控制台输出logo
	fmt.Printf(logo, globalApp.VersionName, addr)

	//启动服务
	go runGinHttpServer(addr)

	//执行Start事件
	if globalEvent.Start != nil {
		globalEvent.Start()
	}

	//守护进程
	if isDaemon {

		// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		log.Info("Shutdown Server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if e := globalApp.Http.Shutdown(ctx); e != nil {
			log.Error("Server Shutdown:" + e.Error())
		}
		log.Info("Server exiting")
	}
}

// 启动Gin服务
func runGinHttpServer(addr string) {

	log.Info("Start gin http server. " + addr)

	globalApp.Http = &http.Server{
		Addr:    addr,
		Handler: GetServer(),
	}

	// 服务连接
	if err := globalApp.Http.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error(err.Error())
	}
}

// RegisterEvent 注册全局系统事件
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
