package cron

import (
	"github.com/go-co-op/gocron"
	"time"
)

var scheduler *gocron.Scheduler

var taskNum = 0

func initScheduler() {
	if scheduler == nil {
		scheduler = gocron.NewScheduler(time.UTC)
	}
}

// AddTask 添加定时任务
func AddTask(role string, call interface{}) error {
	initScheduler()

	_, err := scheduler.Cron(role).Do(call)
	if err != nil {
		return err
	}
	taskNum++
	return nil
}

// Start 启动定时任务
func Start() {
	initScheduler()

	if taskNum > 0 {
		scheduler.StartBlocking()
	}
}
