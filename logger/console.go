package logger

import (
	"fmt"
	"time"
)

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

func log(lv Loglevel, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	now := time.Now()
	funcName, fileName, lineNo := getInfo(3)
	fmt.Printf("[%s] [%s] [%s:%s:%d] %s\n", now.Format("2006-01-02 15:04:05"), getLogString(lv), fileName, funcName, lineNo, msg)
}

func (l Logger) Debug(format string, a ...interface{}) {
	if l.enable(DEBUG) {
		log(DEBUG, format, a...)
	}
}

func (l Logger) Info(format string, a ...interface{}) {
	if l.enable(INFO) {
		log(INFO, format, a...)
	}
}

func (l Logger) Warning(format string, a ...interface{}) {
	if l.enable(WARNING) {
		log(WARNING, format, a...)
	}
}

func (l Logger) Error(format string, a ...interface{}) {
	if l.enable(ERROR) {
		log(ERROR, format, a...)
	}
}

func (l Logger) Fatal(format string, a ...interface{}) {
	if l.enable(FATAL) {
		log(FATAL, format, a...)
	}
}
