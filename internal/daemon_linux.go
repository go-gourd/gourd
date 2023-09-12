//go:build linux

package internal

import (
	"bufio"
	"fmt"
	"github.com/go-gourd/gourd/config"
	"github.com/sevlyar/go-daemon"
	"log"
	"os"
	"strconv"
	"syscall"
)

// DaemonRun 守护进程模式运行
func DaemonRun() {

	tempDir := config.GetAppConfig().Temp
	pidFile := tempDir + "/daemon.pid"

	//先判断进程是否已存在
	pid := getPid(pidFile)
	if pid > 0 {
		// pid存在，提示是否结束之前进程
		fmt.Print("Pid file exists, whether to close the previous process first? [y(default)/n]")
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		v := input.Text()
		fmt.Println(v)
		if v != "n" {
			killPid(pid)
			os.Remove(pidFile)
		}
	}

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

// 检查是否存在pid
func getPid(pidFile string) (pid int) {
	file, err := os.Open(pidFile)
	if err != nil {
		return
	}
	defer file.Close()
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
		syscall.Kill(pid, syscall.SIGINT)
	}
}

// StopDaemonProcess 结束指定进程
func StopDaemonProcess() {
	pidFile := config.GetAppConfig().TempDir + "/daemon.pid"

	pid := getPid(pidFile)
	if pid > 0 {
		// pid存在
		killPid(pid)
		os.Remove(pidFile)
	}
}
