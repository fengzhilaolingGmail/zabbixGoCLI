package cli

import (
	"context"
	"fmt"
	"os"

	"zabbix-maint/internal/config"
)

// ConfigHandler 配置管理命令处理器
type ConfigHandler struct {
	mgr *config.Manager
}

// NewConfigHandler 创建新的配置管理处理器
func NewConfigHandler() *ConfigHandler {
	mgr, err := config.NewManager()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: config manager init failed: %v\n", err)
	}
	return &ConfigHandler{mgr: mgr}
}

// Add 添加实例配置
func (h *ConfigHandler) Add(ctx context.Context, name, url, user, pass string) error {
	if h.mgr == nil {
		return fmt.Errorf("config manager not available")
	}
	if err := h.mgr.AddInstance(ctx, name, url, user, pass, ""); err != nil {
		return err
	}
	fmt.Printf("Instance '%s' added successfully.\n", name)
	return nil
}

// List 列出所有实例配置
func (h *ConfigHandler) List(ctx context.Context) error {
	if h.mgr == nil {
		return fmt.Errorf("config manager not available")
	}
	instances, err := h.mgr.ListInstances()
	if err != nil {
		return err
	}
	if len(instances) == 0 {
		fmt.Println("No instances configured.")
		fmt.Println("Use: zbx-cli config add --name <name> --url <url> --user <user> --pass <password>")
		return nil
	}
	fmt.Println("NAME           ENDPOINT")
	fmt.Println("---------------------------------------------------")
	for _, inst := range instances {
		fmt.Printf("%-14s %s\n", inst.Name, inst.Endpoint)
	}
	return nil
}

// Remove 删除实例配置
func (h *ConfigHandler) Remove(ctx context.Context, name string) error {
	if h.mgr == nil {
		return fmt.Errorf("config manager not available")
	}
	if err := h.mgr.RemoveInstance(ctx, name); err != nil {
		return err
	}
	fmt.Printf("Instance '%s' removed successfully.\n", name)
	return nil
}

// Test 测试实例连接
func (h *ConfigHandler) Test(ctx context.Context, name string) error {
	if h.mgr == nil {
		return fmt.Errorf("config manager not available")
	}
	return h.mgr.TestInstance(ctx, name)
}
