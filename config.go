package cc

import (
	"os"

	"gopkg.in/yaml.v2"
)

// DEFAULT_CONFIG_PATH 默认配置文件地址
const DEFAULT_CONFIG_PATH = "application.yaml"

// AppConfig 项目配置
type AppConfig struct {
	Port         int `yaml:"port"`
	ReadTimeout  int `yaml:"read-timeout"`
	WriteTimeout int `yaml:"write-timeout"`
	IdleTimeout  int `yaml:"idle-timeout"`
	Database     struct {
		Source       string `yaml:"source"`
		MaxOpenConns int    `yaml:"max-open-conns"`
		MaxIdleConns int    `yaml:"max-idle-conns"`
	} `yaml:"database"`
}

// NewAppConfig 构造带默认参数的项目配置
func NewAppConfig() *AppConfig {
	config := &AppConfig{
		Port:         9999,
		ReadTimeout:  5,
		WriteTimeout: 10,
		IdleTimeout:  15,
	}
	config.Database.Source = ""
	config.Database.MaxOpenConns = 10
	config.Database.MaxIdleConns = 5
	return config
}

// ParseFile 解析配置文件
func (config *AppConfig) ParseFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(content, &config)
}

// ParseContent 解析配置内容
func (config *AppConfig) ParseContent(content []byte) error {
	return yaml.Unmarshal(content, &config)
}
