package v7

import (
	"context"
	"fmt"

	"zabbix-maint/internal/api"
	"zabbix-maint/internal/model"
)

// HostFullClone 主机全克隆 (V7: 原生 host.clone)
func (a *V7Adapter) HostFullClone(ctx context.Context, req model.HostCloneReq) (string, error) {
	params := map[string]interface{}{
		"hostid": req.SourceHostID,
		"host": map[string]interface{}{
			"host": req.NewHostName,
			"name": req.NewHostName,
		},
	}
	var result map[string]interface{}
	err := a.client.Call(ctx, "host.clone", params, &result)
	if err != nil {
		return "", err
	}
	return result["hostids"].([]interface{})[0].(string), nil
}

// HostList 查询主机列表
func (a *V7Adapter) HostList(ctx context.Context, filter string) ([]model.UnifiedHost, error) {
	// TODO: implement V7 host.get
	return nil, fmt.Errorf("not implemented")
}

// HostDetail 查询主机详情
func (a *V7Adapter) HostDetail(ctx context.Context, hostID string) (*model.UnifiedHost, error) {
	// TODO: implement V7 host.get
	return nil, fmt.Errorf("not implemented")
}
