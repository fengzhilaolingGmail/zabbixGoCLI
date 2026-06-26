package cli

import (
	"context"
	"fmt"
)

// OneShotHandler 一次性命令处理器
type OneShotHandler struct {
	// TODO: add fields
}

// NewOneShotHandler 创建新的一次性命令处理器
func NewOneShotHandler() *OneShotHandler {
	return &OneShotHandler{}
}

// Execute 执行一次性命令
func (h *OneShotHandler) Execute(ctx context.Context, args []string) error {
	// TODO: implement one-shot command parsing and execution
	return fmt.Errorf("one-shot mode not implemented")
}
