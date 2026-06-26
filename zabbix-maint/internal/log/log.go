package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// Level 日志级别
type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Logger 日志记录器
type Logger struct {
	level  Level
	logger *log.Logger
}

var defaultLogger *Logger

func init() {
	defaultLogger = &Logger{
		level:  LevelInfo,
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

// Init 初始化日志系统
// logFile: 日志文件路径，为空时只输出到 stderr
// level: debug/info/warn/error
func Init(logFile string, level string) error {
	lvl := parseLevel(level)
	var w io.Writer
	if logFile != "" {
		dir := filepath.Dir(logFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("create log dir: %w", err)
		}
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("open log file: %w", err)
		}
		w = f
		fmt.Printf("Log file: %s\n", logFile)
	} else {
		w = os.Stderr
	}
	defaultLogger = &Logger{
		level:  lvl,
		logger: log.New(w, "", log.LstdFlags|log.Lshortfile),
	}
	return nil
}

func parseLevel(s string) Level {
	switch s {
	case "debug":
		return LevelDebug
	case "warn":
		return LevelWarn
	case "error":
		return LevelError
	default:
		return LevelInfo
	}
}

func (l *Logger) logf(level Level, prefix string, format string, v ...interface{}) {
	if l.level <= level {
		l.logger.Printf("[%s] "+format, append([]interface{}{prefix}, v...)...)
	}
}

// Debugf 调试日志
func Debugf(format string, v ...interface{}) { defaultLogger.logf(LevelDebug, "DEBUG", format, v...) }

// Infof 信息日志
func Infof(format string, v ...interface{}) { defaultLogger.logf(LevelInfo, "INFO", format, v...) }

// Warnf 警告日志
func Warnf(format string, v ...interface{}) { defaultLogger.logf(LevelWarn, "WARN", format, v...) }

// Errorf 错误日志
func Errorf(format string, v ...interface{}) { defaultLogger.logf(LevelError, "ERROR", format, v...) }

// SetLevel 动态设置日志级别
func SetLevel(level string) { defaultLogger.level = parseLevel(level) }

// GetLevel 获取当前日志级别
func GetLevel() Level { return defaultLogger.level }
