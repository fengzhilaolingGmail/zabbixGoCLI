package cli

import (
	"context"
	"fmt"
)

// ConfigHandler 配置管理命令处理器
type ConfigHandler struct{}

// NewConfigHandler 创建新的配置管理处理器
func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{}
}

// Add 添加实例配置
func (h *ConfigHandler) Add(ctx context.Context, name, url, user, pass string) error {
	// TODO: implement config add
	return fmt.Errorf("config add not implemented")
}

// List 列出所有实例配置
func (h *ConfigHandler) List(ctx context.Context) error {
	// TODO: implement config list
	return fmt.Errorf("config list not implemented")
}

// Remove 删除实例配置
func (h *ConfigHandler) Remove(ctx context.Context, name string) error {
	// TODO: implement config remove
	return fmt.Errorf("config remove not implemented")
}

// Test 测试实例连接
func (h *ConfigHandler) Test(ctx context.Context, name string) error {
	// TODO: implement config test
	return fmt.Errorf("config test not implemented")
}
