package service

import (
	"context"
	"fmt"

	"zabbix-maint/pkg/zabbix"
)

// BatchService 批量操作服务
type BatchService struct {
	client zabbix.ZabbixOperator
}

// NewBatchService 创建批量操作服务
func NewBatchService(client zabbix.ZabbixOperator) *BatchService {
	return &BatchService{client: client}
}

// BatchUserCreate 批量创建用户
func (s *BatchService) BatchUserCreate(ctx context.Context, users []BatchUserItem) ([]string, error) {
	// TODO: implement batch user creation
	return nil, fmt.Errorf("not implemented")
}

// BatchHostClone 批量克隆主机
func (s *BatchService) BatchHostClone(ctx context.Context, hosts []BatchHostItem) ([]string, error) {
	// TODO: implement batch host clone
	return nil, fmt.Errorf("not implemented")
}

// BatchUserItem 批量创建用户项
type BatchUserItem struct {
	Alias    string
	Name     string
	Passwd   string
	GroupIDs []string
	RoleID   string
}

// BatchHostItem 批量克隆主机项
type BatchHostItem struct {
	SourceHostID string
	NewHostName  string
	NewHostIP    string
	NewGroupIDs  []string
}
