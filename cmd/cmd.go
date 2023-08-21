package cmd

import (
	"errors"
	"fmt"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/core"
	"github.com/go-gourd/gourd/event"
	"github.com/go-gourd/gourd/log"
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

var defaultCmd *Commend

// Add 添加命令行
func Add(cmd Commend) {
	cmdList[cmd.Name] = cmd
}

// SetDefault 设置默认执行命令 -无任何参数的清情况下
func SetDefault(cmd Commend) {
	defaultCmd = &cmd
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
		if defaultCmd == nil {
			fmt.Print(fmt.Sprintf(core.NoCmdHelp, fileName) + getCmdListHelps())
			os.Exit(0)
		}

		//执行默认命令
		defaultCmd.Handler(args)
		return
	}

	// 解析执行命令行
	err := Exec(args[1], args)
	if err != nil {
		fmt.Print(fmt.Sprintf(core.UndefinedCmdHelp, args[1], fileName) + getCmdListHelps())
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
				log.Info("Daemon Running...")
				fmt.Println("[Info] Daemon Running...")
				os.Exit(0)
			}
		}

		// 开启Http监听服务
		if config.GetHttpConfig().Enable {
			event.Trigger("_http_start", nil)
		}
		return false
	case "stop":
		//停止后台进程
		core.StopDaemonProcess()
		log.Info("Daemon Process Stopped.")
		fmt.Println("[Info] Daemon Process Stopped.")
		os.Exit(0)
		return false
	}

	return true
}

// 获取所有自定义命令提示
func getCmdListHelps() string {
	str := ""
	for _, cmd := range cmdList {
		str += fmt.Sprintf("  %-7s  %s\n", cmd.Name, cmd.Explain)
	}
	return str
}
