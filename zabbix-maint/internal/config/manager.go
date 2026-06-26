package config

import (
	"context"
	"fmt"
)

// Manager 配置管理器
type Manager struct {
	config *Config
}

// NewManager 创建配置管理器
func NewManager() (*Manager, error) {
	// TODO: load existing config
	return &Manager{}, nil
}

// AddInstance 添加实例
func (m *Manager) AddInstance(ctx context.Context, name, endpoint, username, password, version string) error {
	// TODO: implement add instance
	return fmt.Errorf("not implemented")
}

// RemoveInstance 删除实例
func (m *Manager) RemoveInstance(ctx context.Context, name string) error {
	// TODO: implement remove instance
	return fmt.Errorf("not implemented")
}

// ListInstances 列出所有实例
func (m *Manager) ListInstances() ([]InstanceConfig, error) {
	// TODO: implement list instances
	return nil, fmt.Errorf("not implemented")
}

// GetInstance 获取实例配置
func (m *Manager) GetInstance(name string) (*InstanceConfig, error) {
	// TODO: implement get instance
	return nil, fmt.Errorf("not implemented")
}

// TestInstance 测试实例连接
func (m *Manager) TestInstance(ctx context.Context, name string) error {
	// TODO: implement test connection
	return fmt.Errorf("not implemented")
}
