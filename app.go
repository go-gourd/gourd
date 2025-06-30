package gourd

import (
	"context"
	"fmt"
	"github.com/go-gourd/gourd/cmd"
	"github.com/go-gourd/gourd/event"
	"log/slog"
	"os"
	"os/signal"
	"runtime"
	"time"
)

// 版本信息
const (
	VersionNum  = 141
	VersionName = "1.4.1"
)

type App struct {
	DisableLogo  bool
	EventHandler event.Handler
	Context      context.Context
	ConfigDir    string
	cancel       context.CancelFunc
}

// Init 初始化应用
func (app *App) Init() {

	if app.Context == nil {
		app.Context = context.Background()
	}

	// 创建运行时可取消上下文
	runCtx, runCancel := context.WithCancel(app.Context)
	app.Context = runCtx
	app.cancel = runCancel

	// 初始化事件
	if app.EventHandler != nil {
		app.EventHandler(app.Context)
	}

	//触发Boot事件
	event.Trigger("app.boot", app.Context)

	//版本显示(优先显示配置文件中的版本号)
	versionName := fmt.Sprintf("%s (%d)", VersionName, VersionNum)

	//输出Logo信息
	if !app.DisableLogo {
		var logo = "   _____                     _ \n" +
			"  / ____|                   | |  Go        %s\n" +
			" | |  __  ___  _   _ _ __ __| |  Gourd     v%s\n" +
			" | | |_ |/ _ \\| | | | '__/ _` |  Platform  %s\n" +
			" | |__| | (_) | |_| | | | (_| |  Arch      %s\n" +
			"  \\_____|\\___/ \\__,_|_|  \\__,_|  Time      %s\n" +
			"------------------------------------------------------------\n"

		fmt.Printf(
			logo, runtime.Version(), versionName, runtime.GOOS,
			runtime.GOARCH, time.Now().Format("06-01-02 15:04:05"),
		)
	}

	// 触发Init事件
	event.Trigger("app.init", app.Context)
}

// Run 启动应用
func (app *App) Run() {

	//命令行解析并执行
	cmd.ConsoleParse()

	// 触发Start事件
	event.Trigger("app.start", app.Context)

	// 守护进程 -等待中断信号以用于关闭服务（设置 10 秒的超时时间）
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	signal.Stop(quit)
	slog.Info("Shutdown Server ...")

	// 通知协程停止
	app.cancel()

	// 创建 10 秒超时上下文用于退出流程
	shutdownCtx, shutdownCancel := context.WithTimeout(app.Context, 10*time.Second)
	defer shutdownCancel()

	// 触发停止事件
	event.Trigger("app.stop", shutdownCtx)

	slog.Info("Server has exited")
}
