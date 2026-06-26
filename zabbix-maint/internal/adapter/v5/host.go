package v5

import (
	"context"
	"fmt"

	"zabbix-maint/pkg/zabbix"
)

// HostFullClone 主机全克隆 (V5 手动实现)
func (a *V5Adapter) HostFullClone(ctx context.Context, req zabbix.HostCloneReq) (string, error) {
	// 1. 获取源主机完整配置
	var hosts []map[string]interface{}
	err := a.client.Call(ctx, "host.get", map[string]interface{}{
		"hostids":               req.SourceHostID,
		"output":                "extend",
		"selectInterfaces":      "extend",
		"selectItems":           "extend",
		"selectTriggers":        "extend",
		"selectGraphs":          "extend",
		"selectMacros":          "extend",
		"selectInventory":       "extend",
		"selectParentTemplates": "extend",
		"selectGroups":          "extend",
	}, &hosts)
	if err != nil || len(hosts) == 0 {
		return "", fmt.Errorf("get source host failed: %w", err)
	}
	src := hosts[0]

	// 2. 构造新主机（移除 ID，修改名称）
	newHost := deepCopyMap(src)
	delete(newHost, "hostid")
	delete(newHost, "items")
	delete(newHost, "triggers")
	delete(newHost, "graphs")
	newHost["host"] = req.NewHostName
	newHost["name"] = req.NewHostName

	// 3. 创建主机
	var created map[string]interface{}
	err = a.client.Call(ctx, "host.create", newHost, &created)
	if err != nil {
		return "", fmt.Errorf("create host failed: %w", err)
	}
	newHostID := created["hostids"].([]interface{})[0].(string)

	// 4. 复制 Items
	if req.CloneItems {
		if items, ok := src["items"].([]interface{}); ok && len(items) > 0 {
			if err := a.cloneItems(ctx, newHostID, items); err != nil {
				// 事务补偿：删除已创建的主机
				a.client.Call(ctx, "host.delete", []string{newHostID}, nil)
				return "", fmt.Errorf("clone items failed: %w", err)
			}
		}
	}

	// 5. TODO: 复制 Triggers, Graphs, Interfaces...

	return newHostID, nil
}

func (a *V5Adapter) cloneItems(ctx context.Context, hostID string, items []interface{}) error {
	// TODO: implement items cloning
	return nil
}

// HostList 查询主机列表
func (a *V5Adapter) HostList(ctx context.Context, filter string) ([]zabbix.UnifiedHost, error) {
	// TODO: implement V5 host.get
	return nil, fmt.Errorf("not implemented")
}

// HostDetail 查询主机详情
func (a *V5Adapter) HostDetail(ctx context.Context, hostID string) (*zabbix.UnifiedHost, error) {
	// TODO: implement V5 host.get
	return nil, fmt.Errorf("not implemented")
}
