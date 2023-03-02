package core

import (
	"github.com/go-gourd/gourd/config"
	"github.com/sevlyar/go-daemon"
	"log"
)

// DaemonRun 守护进程模式运行
func DaemonRun() {

	tempDir := config.GetAppConfig().TempDir

	cntxt := &daemon.Context{
		PidFileName: tempDir + "/daemon.pid",
		PidFilePerm: 0644,
		LogFileName: tempDir + "/daemon.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{"[go-daemon sample]"},
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return
	}
	defer cntxt.Release()
}
