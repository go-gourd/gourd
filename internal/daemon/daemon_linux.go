//go:build linux

package daemon

import (
	"bufio"
	"fmt"
	"github.com/sevlyar/go-daemon"
	"log/slog"
	"os"
	"strconv"
	"syscall"
)

// Run 守护进程模式运行
func Run() {

	pidFile := "./daemon.pid"

	//先判断进程是否已存在
	pid := getPid(pidFile)
	if pid > 0 {
		// pid存在，提示是否结束之前进程
		fmt.Print("Pid file exists, whether to close the previous process first? [y(default)/n]")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		v := input.Text()
		if v != "n" {
			killPid(pid)
			err := os.Remove(pidFile)
			if err != nil {
				return
			}
		}
	}

	ctx := &daemon.Context{
		//PidFileName: tempDir + "/daemon.pid",
		//PidFilePerm: 0644,
		LogFileName: "./daemon.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        []string{os.Args[0], os.Args[1]},
	}

	d, err := ctx.Reborn()
	if err != nil {
		slog.Error("Unable to run: ", err)
	}
	if d != nil {
		_ = os.WriteFile(pidFile, []byte(strconv.Itoa(d.Pid)), 0666) //写入文件(字节数组)
		_ = ctx.Release()
	}
}

// 检查是否存在pid
func getPid(pidFile string) (pid int) {
	file, err := os.Open(pidFile)
	if err != nil {
		return
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		i, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return
		}
		return i
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}
	return
}

// 关闭进程
func killPid(pid int) {
	err := syscall.Kill(pid, 0)
	if err == nil {
		// 进程存在
		_ = syscall.Kill(pid, syscall.SIGINT)
	}
}

// StopDaemonProcess 结束指定进程
func StopDaemonProcess() {
	pidFile := "./daemon.pid"

	pid := getPid(pidFile)
	if pid > 0 {
		// pid存在
		killPid(pid)
		_ = os.Remove(pidFile)
	} else {
		fmt.Println("daemon process pid not exist")
	}
}
