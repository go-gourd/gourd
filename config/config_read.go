package config

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

// 默认配置文件目录
const defaultDir = "./config"

var (
	// 配置文件目录
	configDir = defaultDir
	// 所有配置
	c Config
)

// SetConfigDir 设置文件目录
func SetConfigDir(path string) {
	configDir = path
}

func SetAppConfig(app *AppConfig) {
	c.App = app
}

func GetAppConfig() *AppConfig {

	//已存在 -返回
	if c.App != nil {
		return c.App
	}

	// 初始化配置默认值
	c.App = &AppConfig{
		Name:    "gourd",
		Debug:   false,
		Runtime: "./runtime",
		Temp:    "./runtime/temp",
	}

	// 读取配置文件
	err := Unmarshal("app", c.App)
	if err != nil {
		panic(err)
	}

	return c.App
}

func SetLogConfig(log *LogConfig) {
	c.Log = log
}

func GetLogConfig() *LogConfig {

	//已存在 -返回
	if c.Log != nil {
		return c.Log
	}

	// 初始化配置默认值
	c.Log = &LogConfig{
		Console: true,
	}

	// 读取配置文件
	err := Unmarshal("log", c.Log)
	if err != nil {
		panic(err)
	}

	return c.Log
}

func SetHttpConfig(http *HttpConfig) {
	c.Http = http
}

func GetHttpConfig() *HttpConfig {

	//已存在 -返回
	if c.Http != nil {
		return c.Http
	}

	// 初始化配置默认值
	c.Http = &HttpConfig{
		Enable: false,
		Host:   "0.0.0.0",
		Port:   8080,
		Static: "",
	}

	// 读取配置文件
	err := Unmarshal("http", c.Http)
	if err != nil {
		panic(err)
	}

	return c.Http
}

func SetDatabaseConfig(database *DatabaseConfig) {
	c.Database = database
}

func GetDbConfig() DatabaseConfig {

	//已存在 -返回
	if c.Database != nil {
		return *c.Database
	}

	// 初始化配置默认值
	c.Database = &DatabaseConfig{}

	// 读取配置文件
	err := Unmarshal("database", c.Database)
	if err != nil {
		panic(err)
	}

	return *c.Database
}

// Unmarshal 读取自定义配置文件
func Unmarshal(name string, v any) error {
	var tomlData, err = ReadFile(name)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(tomlData, v)
	if err != nil {
		return err
	}
	return nil
}

func ReadFile(name string) ([]byte, error) {
	var file = configDir + "/" + name + ".toml"
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return f, nil
}
