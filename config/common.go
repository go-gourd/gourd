package config

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

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
