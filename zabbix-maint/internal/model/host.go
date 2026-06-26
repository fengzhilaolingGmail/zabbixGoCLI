package model

// HostCloneReq 主机克隆请求
type HostCloneReq struct {
	SourceHostID  string   // 源主机ID
	NewHostName   string   // 新主机名
	NewHostIP     string   // 新主机IP（可选，覆盖接口IP）
	NewGroupIDs   []string // 新主机组（可选，继承则留空）
	CloneItems    bool     // 是否克隆监控项
	CloneTriggers bool     // 是否克隆触发器
	CloneGraphs   bool     // 是否克隆图形
}

// HostInterface 主机接口
type HostInterface struct {
	Type    int
	Main    int
	UseIP   int
	IP      string
	DNS     string
	Port    string
	Details map[string]interface{}
}

// HostMacro 主机宏
type HostMacro struct {
	Macro  string
	Value  string
	Type   int
	Description string
}

// UnifiedHost 统一主机模型
type UnifiedHost struct {
	ID         string
	Host       string
	Name       string
	Interfaces []HostInterface
	GroupIDs   []string
	GroupNames []string
	Templates  []string
	Macros     []HostMacro
	Inventory  map[string]string
	Status     int // 0=enabled, 1=disabled
}
