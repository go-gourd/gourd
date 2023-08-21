package config

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

var initState = make(map[string]bool)

var configDir = "./config"

var appConfig AppConfig
var httpConfig HttpConfig
var dbConfig DbConfig

// SetConfigDir 设置配置文件目录
func SetConfigDir(path string) {
	configDir = path
}

func GetAppConfig() *AppConfig {
	name := "app"

	//已存在 -返回
	if _, ok := initState[name]; ok {
		return &appConfig
	}

	var tomlData, err = ReadFile(name)
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(tomlData, &appConfig)
	if err != nil {
		panic(err)
	}

	//已加载
	initState[name] = true
	return &appConfig
}

func GetHttpConfig() *HttpConfig {
	name := "http"

	//已存在 -返回
	if _, ok := initState[name]; ok {
		return &httpConfig
	}

	var tomlData, err = ReadFile(name)
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(tomlData, &httpConfig)
	if err != nil {
		panic(err)
	}

	//已加载
	initState[name] = true
	return &httpConfig
}

func GetDbConfig() DbConfig {
	name := "database"

	//已存在 -返回
	if _, ok := initState[name]; ok {
		return dbConfig
	}

	var tomlData, err = ReadFile(name)
	if err != nil {
		panic(err)
	}

	err = toml.Unmarshal(tomlData, &dbConfig)
	if err != nil {
		panic(err)
	}

	//已加载
	initState[name] = true
	return dbConfig
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
