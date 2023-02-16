package config

import (
	"github.com/BurntSushi/toml"
	"os"
)

var initState = make(map[string]bool)

var appConfig AppConfig
var logConfig LogConfig
var httpConfig HttpConfig

func GetAppConfig() *AppConfig {
	name := "app"

	//已存在 -返回
	if _, ok := initState[name]; ok {
		return &appConfig
	}

	var tomlData, err = readConfigFile(name)
	if err != nil {
		panic(err)
	}

	_, err = toml.Decode(tomlData, &appConfig)
	if err != nil {
		panic(err)
	}

	//已加载
	initState[name] = true
	return &appConfig
}

func GetLogConfig() *LogConfig {
	name := "log"

	//已存在 -返回
	if _, ok := initState[name]; ok {
		return &logConfig
	}

	var tomlData, err = readConfigFile(name)
	if err != nil {
		panic(err)
	}

	_, err = toml.Decode(tomlData, &logConfig)
	if err != nil {
		panic(err)
	}

	//已加载
	initState[name] = true
	return &logConfig
}

func GetHttpConfig() *HttpConfig {
	name := "http"

	//已存在 -返回
	if _, ok := initState[name]; ok {
		return &httpConfig
	}

	var tomlData, err = readConfigFile(name)
	if err != nil {
		panic(err)
	}

	_, err = toml.Decode(tomlData, &httpConfig)
	if err != nil {
		panic(err)
	}

	//已加载
	initState[name] = true
	return &httpConfig
}

func readConfigFile(name string) (string, error) {
	var file = "./config/" + name + ".toml"
	f, err := os.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(f), nil
}
