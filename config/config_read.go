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
