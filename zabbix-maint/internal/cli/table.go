package cli

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TableRenderer 表格渲染器
type TableRenderer interface {
	Render(headers []string, rows [][]string)
	RenderWithIndex(headers []string, rows [][]string) // 带序号
	RenderJSON(data interface{})                       // JSON 原始输出
}

// DefaultTableRenderer 默认表格渲染器
type DefaultTableRenderer struct{}

// NewTableRenderer 创建新的 TableRenderer
func NewTableRenderer() TableRenderer {
	return &DefaultTableRenderer{}
}

func (t *DefaultTableRenderer) Render(headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}
	// Print headers
	fmt.Println(strings.Join(headers, "\t"))
	fmt.Println(strings.Repeat("-", len(strings.Join(headers, "\t"))))
	for _, row := range rows {
		fmt.Println(strings.Join(row, "\t"))
	}
}

func (t *DefaultTableRenderer) RenderWithIndex(headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}
	newHeaders := append([]string{"#"}, headers...)
	fmt.Println(strings.Join(newHeaders, "\t"))
	fmt.Println(strings.Repeat("-", len(strings.Join(newHeaders, "\t"))))
	for i, row := range rows {
		newRow := append([]string{fmt.Sprintf("%d", i+1)}, row...)
		fmt.Println(strings.Join(newRow, "\t"))
	}
}

func (t *DefaultTableRenderer) RenderJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}
	fmt.Println(string(b))
}
