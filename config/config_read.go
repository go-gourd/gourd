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
)

var c Config

// SetConfigDir 设置文件目录
func SetConfigDir(path string) {
	configDir = path
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
		Temp:    "./runtime",
		Version: "1.0.0",
	}

	tomlData, err := ReadFile("app")
	if err != nil {
		return c.App
	}

	// 配置文件存在，解析配置文件
	err = toml.Unmarshal(tomlData, c.App)
	if err != nil {
		panic(err)
	}

	return c.App
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

	var tomlData, err = ReadFile("log")
	if err != nil {
		return c.Log
	}

	err = toml.Unmarshal(tomlData, c.Log)
	if err != nil {
		panic(err)
	}

	return c.Log
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

	var tomlData, err = ReadFile("http")
	if err != nil {
		return c.Http
	}

	err = toml.Unmarshal(tomlData, c.Http)
	if err != nil {
		panic(err)
	}

	return c.Http
}

func GetDbConfig() DatabaseConfig {

	//已存在 -返回
	if c.Database != nil {
		return *c.Database
	}

	// 初始化配置默认值
	c.Database = &DatabaseConfig{}

	var tomlData, err = ReadFile("database")
	if err != nil {
		return *c.Database
	}

	err = toml.Unmarshal(tomlData, c.Database)
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
