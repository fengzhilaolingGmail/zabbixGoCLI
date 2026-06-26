package config

import (
	"context"
	"fmt"

	"zabbix-maint/internal/api"
)

// Manager 配置管理器
type Manager struct {
	cfg *Config
}

// NewManager 创建配置管理器
func NewManager() (*Manager, error) {
	cfg, err := LoadAll()
	if err != nil {
		return nil, err
	}
	return &Manager{cfg: cfg}, nil
}

// AddInstance 添加实例
func (m *Manager) AddInstance(ctx context.Context, name, endpoint, username, password, version string) error {
	for _, inst := range m.cfg.Instances {
		if inst.Name == name {
			return fmt.Errorf("instance already exists: %s", name)
		}
	}

	m.cfg.Instances = append(m.cfg.Instances, InstanceConfig{
		Name:     name,
		Endpoint: endpoint,
		Version:  version,
		Auth:     AuthConfig{Username: username, Password: password},
	})

	if m.cfg.Global.DefaultInstance == "" {
		m.cfg.Global.DefaultInstance = name
	}

	return Save(m.cfg)
}

// RemoveInstance 删除实例
func (m *Manager) RemoveInstance(ctx context.Context, name string) error {
	var found bool
	var newInstances []InstanceConfig
	for _, inst := range m.cfg.Instances {
		if inst.Name == name {
			found = true
			continue
		}
		newInstances = append(newInstances, inst)
	}
	if !found {
		return fmt.Errorf("instance not found: %s", name)
	}
	m.cfg.Instances = newInstances
	if m.cfg.Global.DefaultInstance == name {
		m.cfg.Global.DefaultInstance = ""
		if len(newInstances) > 0 {
			m.cfg.Global.DefaultInstance = newInstances[0].Name
		}
	}
	return Save(m.cfg)
}

// ListInstances 列出所有实例
func (m *Manager) ListInstances() ([]InstanceConfig, error) {
	return m.cfg.Instances, nil
}

// GetInstance 获取实例配置
func (m *Manager) GetInstance(name string) (*InstanceConfig, error) {
	for i := range m.cfg.Instances {
		if m.cfg.Instances[i].Name == name {
			inst := m.cfg.Instances[i]
			return &inst, nil
		}
	}
	return nil, fmt.Errorf("instance not found: %s", name)
}

// TestInstance 测试实例连接（调用 apiinfo.version，无需认证）
func (m *Manager) TestInstance(ctx context.Context, name string) error {
	inst, err := m.GetInstance(name)
	if err != nil {
		return err
	}
	fmt.Printf("Testing connection to %s (%s)...\n", inst.Name, inst.Endpoint)

	client := api.NewJSONRPCClient(inst.Endpoint)
	var version string
	if err := client.Call(ctx, "apiinfo.version", nil, &version); err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	fmt.Printf("[OK] Connected! Zabbix API version: %s\n", version)

	// 尝试登录验证凭据
	fmt.Printf("Testing credentials...\n")
	authMgr := api.NewAuthManager(inst.Auth.Username, inst.Auth.Password, client)
	if err := authMgr.Refresh(ctx); err != nil {
		fmt.Printf("[WARN] Login failed: %v\n", err)
		fmt.Printf("[INFO] API endpoint is reachable, but credentials may be incorrect.\n")
		return nil
	}
	fmt.Printf("[OK] Login succeeded!\n")
	return nil
}
