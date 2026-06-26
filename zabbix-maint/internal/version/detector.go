package version

import (
	"context"
	"fmt"
	"strings"

	"zabbix-maint/internal/api"
)

// Version 版本类型
type Version string

const (
	Version5 Version = "5.x"
	Version7 Version = "7.x"
)

func (v Version) String() string {
	return string(v)
}

// DetectVersion 自动探测 Zabbix 版本
func DetectVersion(ctx context.Context, client *api.JSONRPCClient) (Version, error) {
	var version string
	err := client.Call(ctx, "apiinfo.version", nil, &version)
	if err != nil {
		return "", err
	}

	if strings.HasPrefix(version, "5.") {
		return Version5, nil
	} else if strings.HasPrefix(version, "7.") {
		return Version7, nil
	}
	return "", fmt.Errorf("unsupported version: %s", version)
}
