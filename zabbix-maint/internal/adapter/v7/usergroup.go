/*
 * @Author: fengzhilaoling
 * @Date: 2026-06-26 10:32:53
 * @LastEditors: fengzhilaoling
 * @LastEditTime: 2026-06-26 11:13:11
 * @FilePath: \zabbixGoCLI\zabbix-maint\internal\adapter\v7\usergroup.go
 * @Description: 文件详情
 */
package v7

import (
	"context"
	"fmt"

	"zabbix-maint/pkg/zabbix"
)

// UserGroupCreate 创建用户组
func (a *V7Adapter) UserGroupCreate(ctx context.Context, req zabbix.UserGroupCreateReq) (string, error) {
	// TODO: implement V7 usergroup.create
	return "", fmt.Errorf("not implemented")
}

// UserGroupDisable 禁用用户组 (V7: userdirectoryid=0 或移除权限)
func (a *V7Adapter) UserGroupDisable(ctx context.Context, groupID string) error {
	// TODO: implement V7 usergroup.update with userdirectoryid=0
	return fmt.Errorf("not implemented")
}

// UserGroupEnable 启用用户组 (V7: 恢复权限配置)
func (a *V7Adapter) UserGroupEnable(ctx context.Context, groupID string) error {
	// TODO: implement V7 usergroup.update restore rights
	return fmt.Errorf("not implemented")
}

// UserGroupClear 清空用户组
func (a *V7Adapter) UserGroupClear(ctx context.Context, groupID string) error {
	// TODO: implement V7 usergroup.update (清空 userids)
	return fmt.Errorf("not implemented")
}

// UserGroupDelete 删除用户组
func (a *V7Adapter) UserGroupDelete(ctx context.Context, groupID string) error {
	// TODO: implement V7 usergroup.delete
	return fmt.Errorf("not implemented")
}

// UserGroupList 查询用户组列表
func (a *V7Adapter) UserGroupList(ctx context.Context, filter string) ([]zabbix.UnifiedUserGroup, error) {
	// TODO: implement V7 usergroup.get
	return nil, fmt.Errorf("not implemented")
}
