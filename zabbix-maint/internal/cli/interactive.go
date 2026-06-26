package cli

import (
	"context"
	"fmt"
	"os"

	"zabbix-maint/internal/config"
	"zabbix-maint/internal/version"
	"zabbix-maint/pkg/zabbix"
)

// RunInteractiveMode 启动交互式菜单模式
func RunInteractiveMode(ctx context.Context, instanceName string) error {
	// 1. 加载配置并连接
	cfg, err := config.Load(instanceName)
	if err != nil {
		return err
	}

	client, ver, err := connect(ctx, cfg)
	if err != nil {
		return err
	}

	// 2. 构建菜单树
	tree := BuildMenuTree()

	// 3. 主循环
	for {
		clearScreen()
		printHeader(cfg, ver)
		printMenu(tree.Current)

		choice, err := readChoice(len(tree.Current.SubMenus))
		if err != nil {
			continue
		}

		if choice == 0 {
			if tree.Current.Parent == nil {
				fmt.Println("再见!")
				return nil
			}
			tree.Current = tree.Current.Parent
			continue
		}

		selected := tree.Current.SubMenus[choice-1]

		if len(selected.SubMenus) > 0 {
			tree.Current = selected
		} else if selected.Action != nil {
			clearScreen()
			printHeader(cfg, ver)
			fmt.Printf("\n【%s】\n\n", selected.Title)

			if err := selected.Action(ctx, client, ver); err != nil {
				fmt.Printf("❌ 操作失败: %v\n", err)
			}

			fmt.Println("\n按 Enter 继续...")
			fmt.Scanln()
		}
	}
}

func connect(ctx context.Context, cfg *config.InstanceConfig) (zabbix.ZabbixOperator, version.Version, error) {
	// TODO: implement connection logic
	return nil, "", nil
}

func clearScreen() {
	// TODO: implement screen clearing
}

func printHeader(cfg *config.InstanceConfig, ver version.Version) {
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("  Zabbix 维护工具 v1.0\n")
	fmt.Printf("  当前实例: %s [Zabbix %s]\n", cfg.Name, ver)
	fmt.Printf("  连接状态: ✅ 已连接\n")
	fmt.Println("═══════════════════════════════════════════════════════════════")
}

func printMenu(node *MenuNode) {
	fmt.Printf("\n  %s:\n\n", node.Title)
	for i, sub := range node.SubMenus {
		hint := ""
		if sub.VersionHint != "" {
			hint = fmt.Sprintf(" (%s)", sub.VersionHint)
		}
		fmt.Printf("  %d. %s%s\n", i+1, sub.Title, hint)
	}
	fmt.Println("  0. 返回上级")
	fmt.Println()
}

func readChoice(max int) (int, error) {
	fmt.Printf("  请选择 [0-%d]: ", max)
	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil || choice < 0 || choice > max {
		return 0, fmt.Errorf("invalid choice")
	}
	return choice, nil
}
