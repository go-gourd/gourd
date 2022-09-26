package utils

import (
	"os"
	"path/filepath"
)

func GetRootPath() string {
	str, _ := os.Getwd()
	return str
}

func GetRelativePath(path string) string {
	if filepath.IsAbs(path) {
		p, err := filepath.Rel(GetRootPath(), path)
		if err != nil {
			return path
		}
		return p
	} else {
		return path
	}
}

func CheckAndMkdir(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}
	if !os.IsNotExist(err) {
		return nil
	}
	err = os.Mkdir(path, os.ModePerm)
	if err == nil {
		return err
	}
	return nil
}
