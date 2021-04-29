package main

import (
	"logger/logger"
	"time"
)

// 测试我们自己写的日志库
func main() {
	// log := logger.NewLog("debug")
	log := logger.NewFileLogger("info", "./", "xss.log", 10*1024*1024)
	for {
		log.Debug("这是一条Debug日志")
		log.Info("这是一套Info日志")
		log.Warning("这是一套Warning日志")
		id := 10001
		name := "xss"
		log.Error("这是一套Error日志%v %v", id, name)
		log.Fatal("这是一套Fatal日志")
		time.Sleep(time.Second)
	}
}
