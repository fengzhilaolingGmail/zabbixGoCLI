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
	Name     string     `yaml:"name"`
	Endpoint string     `yaml:"endpoint"`
	Version  string     `yaml:"version"`
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

// LoadAll 加载全部配置
func LoadAll() (*Config, error) {
	path := ConfigFilePath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{}, nil
		}
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config failed: %w", err)
	}
	return &cfg, nil
}

// Load 加载指定实例配置
func Load(instanceName string) (*InstanceConfig, error) {
	cfg, err := LoadAll()
	if err != nil {
		return nil, err
	}

	for i := range cfg.Instances {
		if cfg.Instances[i].Name == instanceName {
			inst := cfg.Instances[i]
			return &inst, nil
		}
	}

	if instanceName == "" {
		defaultName := cfg.Global.DefaultInstance
		if defaultName != "" {
			for i := range cfg.Instances {
				if cfg.Instances[i].Name == defaultName {
					inst := cfg.Instances[i]
					return &inst, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("instance not found: %s", instanceName)
}

// Save 保存配置
func Save(cfg *Config) error {
	path := ConfigFilePath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}
