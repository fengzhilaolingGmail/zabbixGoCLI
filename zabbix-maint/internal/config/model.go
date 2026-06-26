package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ConfigFileName 默认配置文件名
const ConfigFileName = "config.yaml"

// ConfigDir 配置目录
func ConfigDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".zbx-cli"
	}
	return filepath.Join(home, ".zbx-cli")
}

// ConfigFilePath 配置文件完整路径
func ConfigFilePath() string {
	return filepath.Join(ConfigDir(), ConfigFileName)
}

// GlobalConfig 全局配置
type GlobalConfig struct {
	Timeout         string `yaml:"timeout"`
	LogLevel        string `yaml:"log_level"`
	DefaultInstance string `yaml:"default_instance"`
}

// RetryConfig 重试配置
type RetryConfig struct {
	MaxRetries int    `yaml:"max_retries"`
	BaseDelay  string `yaml:"base_delay"`
}

// InstanceConfig 实例配置
type InstanceConfig struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Version  string `yaml:"version"`
	Auth     AuthConfig `yaml:"auth"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Config 总配置
type Config struct {
	Instances []InstanceConfig `yaml:"instances"`
	Global    GlobalConfig     `yaml:"global"`
	Retry     RetryConfig      `yaml:"retry"`
}

// Load 加载配置
func Load(instanceName string) (*InstanceConfig, error) {
	// TODO: implement config loading
	return nil, fmt.Errorf("config loading not implemented")
}

// Save 保存配置
func Save(cfg *Config) error {
	// TODO: implement config saving
	return fmt.Errorf("config saving not implemented")
}
