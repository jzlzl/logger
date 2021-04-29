package logger

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

// 往终端写日志相关内容

type Loglevel uint16

const (
	UNKNOWN Loglevel = iota
	TRACE
	DEBUG
	INFO
	WARNING
	ERROR
	FATAL
)

func parseLogLevel(s string) (Loglevel, error) {
	s = strings.ToLower(s)
	switch s {
	case "debug":
		return DEBUG, nil
	case "trace":
		return TRACE, nil
	case "info":
		return INFO, nil
	case "warning":
		return WARNING, nil
	case "error":
		return ERROR, nil
	case "fatal":
		return FATAL, nil
	default:
		err := errors.New("无效的日志级别")
		return UNKNOWN, err
	}
}

// Logger日志结构体
type Logger struct {
	Level Loglevel
}

// NewLog构造函数
func NewLog(levelstr string) Logger {
	level, err := parseLogLevel(levelstr)
	if err != nil {
		panic(err)
	}
	return Logger{
		Level: level,
	}
}

func (l Logger) enable(logLevel Loglevel) bool {
	return logLevel >= l.Level
}

func (l Logger) Debug(msg string) {
	if l.enable(DEBUG) {
		now := time.Now()
		fmt.Printf("[%s] [DEBUG] %s\n", now.Format("2006-01-02 15:04:05"), msg)
	}
}

func (l Logger) Info(msg string) {
	if l.enable(INFO) {
		now := time.Now()
		fmt.Printf("[%s] [INFO] %s\n", now.Format("2006-01-02 15:04:05"), msg)
	}
}

func (l Logger) Warning(msg string) {
	if l.enable(WARNING) {
		now := time.Now()
		fmt.Printf("[%s] [WARNING] %s\n", now.Format("2006-01-02 15:04:05"), msg)
	}
}

func (l Logger) Error(msg string) {
	if l.enable(ERROR) {
		now := time.Now()
		fmt.Printf("[%s] [ERROR] %s\n", now.Format("2006-01-02 15:04:05"), msg)
	}
}

func (l Logger) Fatal(msg string) {
	if l.enable(FATAL) {
		now := time.Now()
		fmt.Printf("[%s] [FATAL] %s\n", now.Format("2006-01-02 15:04:05"), msg)
	}
}
