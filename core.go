package gourd

import (
	"context"
	"fmt"
	"github.com/go-gourd/gourd/cmd"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/core"
	"github.com/go-gourd/gourd/event"
	"github.com/go-gourd/gourd/ghttp"
	"github.com/go-gourd/gourd/logger"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"
)

func (app *App) Run() {
	//触发Init事件
	event.OnEvent("_init", nil)

	// 开启Http监听服务
	if config.GetHttpConfig().Enable {
		go ghttp.RunHttpServer()
	}

	//触发Start事件
	event.OnEvent("_start", nil)

	// 守护进程 -等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("Shutdown Server ...", zap.Skip())

	//触发停止事件
	event.OnEvent("_stop", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//关闭Http服务
	ghttp.HttpServerShutdown(ctx)

	logger.Info("Server exiting", zap.Skip())
}

// 解析并运行命令行
func consoleParse() {

	//取出运行参数
	args := os.Args

	// 获取当前可执行文件名称
	filenameWithSuffix := path.Base(strings.Replace(args[0], "\\", "/", -1))
	fileName := strings.TrimSuffix(filenameWithSuffix, path.Ext(filenameWithSuffix))

	if len(args) == 1 {
		fmt.Println(fmt.Sprintf(core.NoCmdHelp, fileName))
		os.Exit(0)
	}

	// 内部命令
	if args[1] == "start" {
		//TODO: 暂未实现命令行接管，继续往下执行即可
		return
	}

	// 自定义命令
	if cmd.ExecCmd(args[1], args) != nil {
		fmt.Println(fmt.Sprintf(core.UndefineCmddHelp, args[1], fileName))
	}
	os.Exit(0)

}
