package main

import (
	"logger/logger"
	"time"
)

// 测试我们自己写的日志库
func main() {
	log := logger.NewLog("error")
	for {
		log.Debug("这是一条Debug日志")
		log.Info("这是一套Info日志")
		log.Warning("这是一套Warning日志")
		log.Error("这是一套Error日志")
		log.Fatal("这是一套Fatal日志")
		time.Sleep(time.Second)
	}
}
