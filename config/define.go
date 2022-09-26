package config

type AppConfig struct {
	Name         string //应用名称
	ReleaseMode  string //应用模式
	Ip           string
	Port         int    //Web端口
	Debug        bool   //是否调试模式
	CodeRootPath string //代码主目录（开发项目目录）
}

type LogConfig struct {
	Level   string //日志记录级别
	LogFile string //日志文件
	Console bool   //是否开启控制台输出
}

type Database struct {
	Driver string //驱动
	Db     string
	Host   string
	Port   uint
	User   string
	Pass   string
}
