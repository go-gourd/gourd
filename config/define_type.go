package config

type AppConfig struct {
	Name  string //应用名称
	Debug bool   //调试模式
}

type HttpConfig struct {
	Host   string // 监听域名、IP
	Port   uint32 // 监听端口
	Public string //静态资源目录
}

type LogConfig struct {
	Level   string //日志记录级别
	LogFile string //日志文件
	Console bool   //是否开启控制台输出
}
