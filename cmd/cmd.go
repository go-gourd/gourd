package cmd

import (
	"errors"
	"fmt"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/core"
	"github.com/go-gourd/gourd/ghttp"
	"github.com/go-gourd/gourd/logger"
	"os"
	"path"
	"runtime"
	"strings"
)

type Commend struct {
	Name    string
	Explain string
	Handler func(args []string)
}

var cmdList = make(map[string]Commend)

// Add 添加命令行
func Add(cmd Commend) {
	cmdList[cmd.Name] = cmd
}

// Exec 执行命令行（由框架完成此操作）
func Exec(name string, args []string) error {

	//首先判断是否内置命令
	ok := coreCmdExec(args)
	if !ok {
		return nil
	}

	if _, ok := cmdList[name]; !ok {
		return errors.New("command `" + name + "` does not exist")
	}

	//执行匹配到的命令
	cmdList[name].Handler(args)
	return nil
}

func ConsoleParse() {

	// 取出运行参数
	args := os.Args

	// 获取当前可执行文件名称
	filenameWithSuffix := path.Base(strings.Replace(args[0], "\\", "/", -1))
	fileName := strings.TrimSuffix(filenameWithSuffix, path.Ext(filenameWithSuffix))

	if len(args) == 1 {
		fmt.Println(fmt.Sprintf(core.NoCmdHelp, fileName))
		os.Exit(0)
	}

	// 解析执行命令行
	err := Exec(args[1], args)
	if err != nil {
		fmt.Println(fmt.Sprintf(core.UndefinedCmdHelp, args[1], fileName))
		os.Exit(0)
	}
}

// 系统内置命令行处理
func coreCmdExec(args []string) bool {

	// 内部命令
	switch args[1] {
	case "restart":
		args = []string{args[0], args[1], "-d"}
		fallthrough
	case "start":
		//守护进程
		if len(args) >= 3 && args[2] == "-d" {
			if runtime.GOOS == "windows" {
				//守护进程模式暂不支持windows
				fmt.Println("[Warn] The daemon does not support Windows.")
			} else {
				//守护进程，成功后会终止当前应用
				core.DaemonRun()
				logger.Info("Daemon Running...")
				fmt.Println("[Info] Daemon Running...")
				os.Exit(0)
			}
		}

		// 开启Http监听服务
		if config.GetHttpConfig().Enable {
			go ghttp.RunHttpServer()
		}
		return false
	case "stop":
		//停止后台进程
		core.StopDaemonProcess()
		logger.Info("Daemon Process Stopped.")
		fmt.Println("[Info] Daemon Process Stopped.")
		os.Exit(0)
		return false
	}

	return true
}
