package core

import (
	"github.com/go-gourd/gourd/config"
	"github.com/go-gourd/gourd/logger"
)

// InitLogger 初始化日志工具
func InitLogger() {

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
