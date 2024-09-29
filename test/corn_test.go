package test

import (
    "fmt"
    "github.com/go-gourd/gourd/cron"
    "log/slog"
    "testing"
    "time"
)

func TestCorn(t *testing.T) {

    _ = cron.Add("* * * * *", func() {
        slog.Info("定时任务示例" + time.Now().Format(time.TimeOnly))
    })

    list := cron.Entries()
    fmt.Println(len(list))

    // 等待2分钟，让定时任务执行
    time.Sleep(2 * time.Minute)

    // 停止定时任务
    cron.Stop()

}
