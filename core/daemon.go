package core

import (
	"github.com/go-gourd/gourd/config"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"strconv"
)

// DaemonRun 守护进程模式运行
func DaemonRun() {

	tempDir := config.GetAppConfig().TempDir
	pidFile := tempDir + "/daemon.pid"

	ctx := &daemon.Context{
		//PidFileName: tempDir + "/daemon.pid",
		//PidFilePerm: 0644,
		LogFileName: tempDir + "/daemon.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{os.Args[0], os.Args[1]},
	}

	d, err := ctx.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		_ = os.WriteFile(pidFile, []byte(strconv.Itoa(d.Pid)), 0666) //写入文件(字节数组)
		_ = ctx.Release()
	}
}
