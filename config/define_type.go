package config

type AppConfig struct {
	Name    string `toml:"name"`    //应用名称
	Debug   bool   `toml:"debug"`   //调试模式
	TempDir string `toml:"tempDir"` //调试模式
	Version string `toml:"version"` //版本号
}

type HttpConfig struct {
	Enable bool   `toml:"enable"` //是否启用Http服务器
	Host   string `toml:"host"`   //监听域名、IP
	Port   uint32 `toml:"port"`   //监听端口
	Public string `toml:"public"` //静态资源目录
}

type DbConfigType struct {
	Type        string `toml:"type"`        //数据库类型
	Host        string `toml:"host"`        //连接地址
	Port        int    `toml:"port"`        //端口
	User        string `toml:"user"`        //用户
	Pass        string `toml:"pass"`        //密码
	Database    string `toml:"database"`    //数据库名
	Param       string `toml:"param"`       //连接参数
	SlowLogTime int    `toml:"slowLogTime"` //慢日志阈值（毫秒）0为不开启
}

// DbConfig 适用于多个连接的配置
type DbConfig map[string]DbConfigType
