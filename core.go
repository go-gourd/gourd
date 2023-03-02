package gourd

import (
	"fmt"
	"github.com/go-gourd/gourd/cmd"
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/core"
	"github.com/go-gourd/gourd/logger"
	"os"
	"path"
	"runtime"
	"strings"
)

// 解析并运行命令行
func consoleParse() {

	// 取出运行参数
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
		//守护进程
		if len(args) >= 3 && args[2] == "-d" {
			if runtime.GOOS == "windows" {
				//守护进程模式暂不支持windows
				fmt.Println("Warn: The daemon does not support Windows.")
			} else {

				//守护进程，成功后会终止当前应用
				core.DaemonRun()

			}
		}

		//暂未实现命令行接管，继续往下执行即可
		return
	}

	// 自定义命令
	if cmd.ExecCmd(args[1], args) != nil {
		fmt.Println(fmt.Sprintf(core.UndefineCmddHelp, args[1], fileName))
	}
	os.Exit(0)

}

func initLogger() {

	conf := config.GetLogConfig()

	c := logger.New()
	c.SetDivision("time")     // 设置归档方式，"time"时间归档 "size" 文件大小归档
	c.SetTimeUnit(logger.Day) // 时间归档 可以设置切割单位
	c.SetEncoding("json")     // 输出格式 "json" 或者 "console"

	if !conf.Console {
		c.CloseConsoleDisplay()
	}

	c.SetInfoFile(conf.LogFile) // 设置info级别日志文件
	if conf.LogErrorFile != "" {
		c.SetErrorFile(conf.LogErrorFile) // 设置warn级别日志文件
	}

	// 设置最低记录级别
	c.SetMinLevel(logger.ParseLevel(conf.Level))

	c.InitLogger()
}
