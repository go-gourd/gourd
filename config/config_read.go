package config

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

var (
	// 配置文件目录
	configDir = "./config"

	// 配置缓存
	appConfig  *AppConfig
	logConfig  *LogConfig
	httpConfig *HttpConfig
	dbConfig   *DbConfig
)

// SetConfigDir 设置配置文件目录
func SetConfigDir(path string) {
	configDir = path
}

func GetAppConfig() *AppConfig {

	//已存在 -返回
	if appConfig != nil {
		return appConfig
	}

	// 初始化配置默认值
	appConfig = &AppConfig{
		Name:    "gourd",
		Debug:   false,
		TempDir: "./runtime",
		Version: "1.0.0",
	}

	tomlData, err := ReadFile("app")
	if err != nil {
		return appConfig
	}

	// 配置文件存在，解析配置文件
	err = toml.Unmarshal(tomlData, &appConfig)
	if err != nil {
		panic(err)
	}

	return appConfig
}

func GetLogConfig() *LogConfig {

	//已存在 -返回
	if logConfig != nil {
		return logConfig
	}

	// 初始化配置默认值
	logConfig = &LogConfig{
		Console: true,
	}

	var tomlData, err = ReadFile("log")
	if err != nil {
		return logConfig
	}

	err = toml.Unmarshal(tomlData, &logConfig)
	if err != nil {
		panic(err)
	}

	return logConfig
}

func GetHttpConfig() *HttpConfig {

	//已存在 -返回
	if httpConfig != nil {
		return httpConfig
	}

	// 初始化配置默认值
	httpConfig = &HttpConfig{
		Enable: false,
		Host:   "0.0.0.0",
		Port:   8080,
		Public: "",
	}

	var tomlData, err = ReadFile("http")
	if err != nil {
		return httpConfig
	}

	err = toml.Unmarshal(tomlData, &httpConfig)
	if err != nil {
		panic(err)
	}

	return httpConfig
}

func GetDbConfig() DbConfig {

	//已存在 -返回
	if dbConfig != nil {
		return *dbConfig
	}

	// 初始化配置默认值
	dbConfig = &DbConfig{}

	var tomlData, err = ReadFile("database")
	if err != nil {
		return *dbConfig
	}

	err = toml.Unmarshal(tomlData, dbConfig)
	if err != nil {
		panic(err)
	}

	return *dbConfig
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
