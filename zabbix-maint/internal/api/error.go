package api

import (
	"fmt"
)

// ZabbixError Zabbix 错误
type ZabbixError struct {
	Code    string
	Message string
	Version string
	Raw     error
}

func (e *ZabbixError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// Unwrap 返回原始错误
func (e *ZabbixError) Unwrap() error {
	return e.Raw
}

// 错误码定义
const (
	ErrCodeAuthFailed          = "AUTH_001"  // 认证失败
	ErrCodeVersionMismatch     = "VER_001"   // 版本不匹配
	ErrCodeAPINotFound         = "API_001"   // API 方法不存在（版本差异）
	ErrCodeParamInvalid        = "PARAM_001" // 参数校验失败
	ErrCodeResourceExist       = "RES_001"   // 资源已存在
	ErrCodeResourceNotFound    = "RES_002"   // 资源不存在
	ErrCodeCloneFailed         = "CLONE_001" // 克隆失败
	ErrCodeUserCancel          = "USER_001"  // 用户取消操作
	ErrCodeVersionNotSupport   = "VER_002"   // 当前版本不支持该功能
)
