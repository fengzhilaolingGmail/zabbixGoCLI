package cli

import (
	"fmt"
	"strings"
)

// Option 选择器选项
type Option struct {
	ID       string
	Label    string
	Detail   string
	Disabled bool // 不可选
}

// PromptReader 交互式输入读取器
type PromptReader interface {
	// 基础输入
	String(prompt string, required bool) (string, error)
	Password(prompt string) (string, error)
	Int(prompt string, min, max int) (int, error)
	Confirm(prompt string, defaultYes bool) (bool, error)

	// 选择器
	Select(prompt string, options []Option) (string, error)           // 单选
	MultiSelect(prompt string, options []Option) ([]string, error)    // 多选
	TableSelect(prompt string, headers []string, rows [][]string) (int, error) // 表格选择

	// 搜索
	SearchSelect(prompt string, searcher func(query string) ([]Option, error)) (string, error)
}

// DefaultPromptReader 默认提示读取器实现
type DefaultPromptReader struct{}

// NewPromptReader 创建新的 PromptReader
func NewPromptReader() PromptReader {
	return &DefaultPromptReader{}
}

func (r *DefaultPromptReader) String(prompt string, required bool) (string, error) {
	fmt.Printf("→ %s: ", prompt)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return "", err
	}
	if required && strings.TrimSpace(input) == "" {
		return "", fmt.Errorf("input is required")
	}
	return input, nil
}

func (r *DefaultPromptReader) Password(prompt string) (string, error) {
	fmt.Printf("→ %s: ", prompt)
	// TODO: use term package for hidden input
	var input string
	_, err := fmt.Scanln(&input)
	return input, err
}

func (r *DefaultPromptReader) Int(prompt string, min, max int) (int, error) {
	fmt.Printf("→ %s [%d-%d]: ", prompt, min, max)
	var val int
	_, err := fmt.Scanln(&val)
	if err != nil || val < min || val > max {
		return 0, fmt.Errorf("invalid input")
	}
	return val, nil
}

func (r *DefaultPromptReader) Confirm(prompt string, defaultYes bool) (bool, error) {
	defaultStr := "Y/n"
	if !defaultYes {
		defaultStr = "y/N"
	}
	fmt.Printf("⚠️  %s [%s]: ", prompt, defaultStr)
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return defaultYes, nil
	}
	input = strings.ToLower(strings.TrimSpace(input))
	if input == "" {
		return defaultYes, nil
	}
	return input == "y" || input == "yes", nil
}

func (r *DefaultPromptReader) Select(prompt string, options []Option) (string, error) {
	fmt.Printf("→ %s:\n", prompt)
	for i, opt := range options {
		marker := "  "
		if opt.Disabled {
			marker = "  [禁用] "
		}
		fmt.Printf("  %s%d. %s\n", marker, i+1, opt.Label)
	}
	fmt.Print("  请选择: ")
	var choice int
	_, err := fmt.Scanln(&choice)
	if err != nil || choice < 1 || choice > len(options) {
		return "", fmt.Errorf("invalid selection")
	}
	return options[choice-1].ID, nil
}

func (r *DefaultPromptReader) MultiSelect(prompt string, options []Option) ([]string, error) {
	fmt.Printf("→ %s (多选, 逗号分隔):\n", prompt)
	for i, opt := range options {
		fmt.Printf("  %d. %s\n", i+1, opt.Label)
	}
	fmt.Print("  请选择: ")
	var input string
	_, err := fmt.Scanln(&input)
	if err != nil {
		return nil, err
	}
	// TODO: parse comma-separated indices
	return nil, fmt.Errorf("not implemented")
}

func (r *DefaultPromptReader) TableSelect(prompt string, headers []string, rows [][]string) (int, error) {
	// TODO: implement table selection
	return 0, fmt.Errorf("not implemented")
}

func (r *DefaultPromptReader) SearchSelect(prompt string, searcher func(query string) ([]Option, error)) (string, error) {
	fmt.Printf("→ %s (支持搜索): ", prompt)
	var query string
	_, err := fmt.Scanln(&query)
	if err != nil {
		return "", err
	}
	options, err := searcher(query)
	if err != nil {
		return "", err
	}
	return r.Select("请选择", options)
}
