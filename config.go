package cc

import (
	"os"

	"gopkg.in/yaml.v3"
)

type appConfig struct {
	App struct {
		Name string `yaml:"name"`
	}
	Server struct {
		Port int32 `yaml:"port"`
	}
	Database struct {
		DriverName string `yaml:"driver-name"`
		Url        string `yaml:"url"`
		Logger     bool   `yaml:"logger"`
	}
}

func (config *appConfig) Default() {
	config.Server.Port = 3000
	config.Database.Logger = false
}

func (config *appConfig) IsExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == os.ErrNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (config *appConfig) Parse(filePath string) error {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(file, &config); err != nil {
		return err
	}
	return nil
}
