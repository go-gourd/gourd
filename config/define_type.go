package config

// AppConfig 应用配置
type AppConfig struct {
	Name       string `toml:"name" json:"name"`               // 应用名称
	Debug      bool   `toml:"debug" json:"debug"`             // 调试模式
	Runtime    string `toml:"runtime" json:"runtime"`         // 运行时目录
	Temp       string `toml:"temp" json:"temp"`               // 临时文件目录
	Version    string `toml:"version" json:"version"`         // 版本名
	VersionNum int    `toml:"version_num" json:"version_num"` // 版本号
}

// HttpConfig Http服务器配置
type HttpConfig struct {
	Enable bool   `toml:"enable" json:"enable"` // 是否启用Http服务器
	Host   string `toml:"host" json:"host"`     // 监听域名、IP
	Port   uint32 `toml:"port" json:"port"`     // 监听端口
	Static string `toml:"static" json:"static"` // 静态资源目录
}

// LogConfig 日志配置
type LogConfig struct {
	Level     string `toml:"level" json:"level"`           // 日志记录级别
	LogFile   string `toml:"log_file" json:"log_file"`     // 日志文件
	ErrorFile string `toml:"error_file" json:"error_file"` // 错误日志文件 -默认不独立存放
	Console   bool   `toml:"console" json:"console"`       // 是否开启控制台输出
	Encoding  string `toml:"encoding" json:"encoding"`     // 输出格式 "json" 或者 "console"
	Division  string `toml:"division" json:"division"`     // 归档方式 time、size
	TimeUnit  string `toml:"time_unit" json:"time_unit"`   // 按时间分割日志的单位 year、month、day、hour、minute
}

// DatabaseConfigType 适用于单个连接的配置
type DatabaseConfigType struct {
	Type        string `toml:"type" json:"type"`                   // 数据库类型
	Host        string `toml:"host" json:"host"`                   // 连接地址
	Port        int    `toml:"port" json:"port"`                   // 端口
	User        string `toml:"user" json:"user"`                   // 用户
	Pass        string `toml:"pass" json:"pass"`                   // 密码
	Database    string `toml:"database" json:"database"`           // 数据库名
	Param       string `toml:"param" json:"param"`                 // 更多连接参数
	SlowLogTime int    `toml:"slow_log_time" json:"slow_log_time"` // 慢日志阈值（毫秒）0为不开启
}

// DatabaseConfig 适用于多个连接的配置
type DatabaseConfig map[string]DatabaseConfigType

// Config 所有配置
type Config struct {
	App      *AppConfig      `toml:"app" json:"app"`
	Http     *HttpConfig     `toml:"http" json:"http"`
	Log      *LogConfig      `toml:"log" json:"log"`
	Database *DatabaseConfig `toml:"database" json:"database"`
}
