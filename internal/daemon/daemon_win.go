//go:build windows || darwin

package daemon

// Run 守护进程模式运行
func Run() {
	//win平台不支持
}

// StopDaemonProcess 结束指定进程
func StopDaemonProcess() {
}
