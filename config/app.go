package config

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type AppConfig struct {
	Name        string //应用名称
	ReleaseMode string //应用模式
	Ip          string
	Port        int  //Web端口
	Debug       bool //是否调试模式
}

func GetConfig(name string, cfg interface{}) error {

	var file = "./config/" + name + ".toml"
	var doc, err = ReadFile(file)
	if err != nil {
		return err
	}

	err = toml.Unmarshal(doc, cfg)
	if err != nil {
		return err
	}

	return nil
}

func ReadFile(file string) ([]byte, error) {
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return f, nil
}
