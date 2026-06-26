package api

import (
	"context"
	"fmt"
	"math"
	"time"
)

// RetryPolicy 重试策略
type RetryPolicy struct {
	MaxRetries     int
	BaseDelay      time.Duration
	MaxDelay       time.Duration
	RetryableCodes []int // HTTP 状态码：502, 503, 504, 429
}

// DefaultRetryPolicy 默认重试策略
var DefaultRetryPolicy = RetryPolicy{
	MaxRetries:     3,
	BaseDelay:      500 * time.Millisecond,
	MaxDelay:       5 * time.Second,
	RetryableCodes: []int{502, 503, 504, 429},
}

// RetryableError 可重试的错误
type RetryableError struct {
	Code   int
	Reason string
}

func (e *RetryableError) Error() string {
	return fmt.Sprintf("retryable error: %s (code: %d)", e.Reason, e.Code)
}

// WithRetry 带重试执行函数
func WithRetry(ctx context.Context, policy RetryPolicy, fn func() error) error {
	var lastErr error
	for i := 0; i <= policy.MaxRetries; i++ {
		if err := fn(); err != nil {
			lastErr = err
			if !isRetryable(err, policy.RetryableCodes) {
				return err
			}
			if i < policy.MaxRetries {
				delay := calculateBackoff(i, policy.BaseDelay, policy.MaxDelay)
				time.Sleep(delay)
			}
			continue
		}
		return nil
	}
	return fmt.Errorf("max retries exceeded: %w", lastErr)
}

func isRetryable(err error, retryableCodes []int) bool {
	if _, ok := err.(*RetryableError); ok {
		return true
	}
	// TODO: check HTTP status codes
	return false
}

func calculateBackoff(attempt int, baseDelay, maxDelay time.Duration) time.Duration {
	delay := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
	if delay > maxDelay {
		delay = maxDelay
	}
	return delay
}
