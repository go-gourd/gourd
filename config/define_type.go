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

type DbConfig struct {
	Host     string //连接地址
	Port     int    //端口
	User     string //用户
	Pass     string //密码
	Database string //数据库名
	Param    string //连接参数
}
