package cron

import (
	"github.com/go-gourd/gourd/event"
	"github.com/robfig/cron/v3"
)

//cron 表达式使用 6 个空格分隔的字段表示一组时间。
//Field name   | Mandatory? | Allowed values  | Allowed special characters
//----------   | ---------- | --------------  | --------------------------
//Seconds      | Yes        | 0-59            | * / , -
//Minutes      | Yes        | 0-59            | * / , -
//Hours        | Yes        | 0-23            | * / , -
//Day of month | Yes        | 1-31            | * / , - ?
//Month        | Yes        | 1-12 or JAN-DEC | * / , -
//Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?

var c *cron.Cron

// Init 初始化定时任务 -由框架自动完成
func Init() {
	if c != nil {
		return
	}
	c = cron.New()

	//系统启动
	event.Listen("_start", func(_ any) {
		// 开始执行（每个任务会在自己的 goroutine 中执行）
		c.Start()
	})
}

// Add 添加一个定时任务
func Add(spec string, cmd func()) error {
	Init()

	c.Run()
	_, err := c.AddFunc(spec, cmd)
	return err
}

// Entries 返回任务列表
func Entries() []cron.Entry {
	return c.Entries()
}

// Stop 停止定时任务运行
func Stop() {
	if c == nil {
		return
	}
	c.Stop()
}
