//go:build windows || darwin

package daemon

// DaemonRun 守护进程模式运行
func DaemonRun() {
	//win平台不支持
}

// StopDaemonProcess 结束指定进程
func StopDaemonProcess() {
}
