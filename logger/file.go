package logger

import (
	"fmt"
	"os"
	"path"
	"time"
)

// 往文件里面写日志相关代码

type FileLogger struct {
	Level       Loglevel
	filePath    string // 日志文件保存的路径
	fileName    string // 日志文件保存的文件名
	fileObj     *os.File
	errfileObj  *os.File
	maxFileSize int64
}

// NewFileLogger构造函数
func NewFileLogger(levelStr, fp, fn string, maxSize int64) *FileLogger {
	logLevel, err := parseLogLevel(levelStr)
	if err != nil {
		panic(err)
	}
	f1 := &FileLogger{
		Level:       logLevel,
		filePath:    fp,
		fileName:    fn,
		maxFileSize: maxSize,
	}
	err = f1.initFile() // 按照文件路径与文件名将文件打开
	if err != nil {
		panic(err)
	}
	return f1
}

// 判断文件是否需要切割
func (f *FileLogger) checkSize(file *os.File) bool {
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed,err: %v\n", err)
		return false
	}
	// 如果当前文件大小大于等于日志文件的最大值就应该返回true
	return fileInfo.Size() >= f.maxFileSize
}

// 根据指定的日志文件路径和文件名打开日志文件
func (f *FileLogger) initFile() error {
	fullFileName := path.Join(f.filePath, f.fileName)
	fileObj, err := os.OpenFile(fullFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file failed, err:%v\n", err)
		return err
	}
	errFileObj, err := os.OpenFile(fullFileName+".err", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open log file failed,err: %v\n", err)
		return err
	}
	// 日志文件都已经打开了
	f.fileObj = fileObj
	f.errfileObj = errFileObj
	return nil
}

// 判断是否需要记录该日志
func (f *FileLogger) enable(logLevel Loglevel) bool {
	return logLevel >= f.Level
}

// 切割文件
func (f *FileLogger) splitFile(file *os.File) (*os.File, error) {
	// 需要切割的日志文件
	nowStr := time.Now().Format("20060102150405000")
	// 通过传进来的file变量拿到文件名称，不能固定名称，因为有错误日志与正常日志
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Printf("get file info failed,err: %v\n", err)
		return nil, err
	}
	logName := path.Join(f.filePath, fileInfo.Name())      // 拿到当前的日志文件完整路径
	newLogName := fmt.Sprintf("%s.bak%s", logName, nowStr) // 拼接一个日志文件备份的名字
	// 1. 关闭当前的日志文件
	file.Close()
	// 2. 备份一下 rename xx.log -> xx.log.bak202104301018
	os.Rename(logName, newLogName)
	// 3. 打开一个新的日志文件
	fileObj, err := os.OpenFile(logName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("open new log file failed,err: %v\n", err)
		return nil, err
	}
	// 4. 将打开的新日志文件对象赋值给f.fileObj
	// f.fileObj = fileObj // 这里不能再去赋值了，因为不知道切的是日志文件还是错误日志文件
	return fileObj, nil
}

// 记录日志的方法
func (f *FileLogger) log(lv Loglevel, format string, a ...interface{}) {
	if f.enable(lv) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now()
		funcName, fileName, lineNo := getInfo(3)
		if f.checkSize(f.fileObj) {
			newFile, err := f.splitFile(f.fileObj) // 日志文件
			if err != nil {
				return
			}
			f.fileObj = newFile // 到这里才赋值
		}
		fmt.Fprintf(f.fileObj, "[%s] [%s] [%s:%s:%d] %s\n", now.Format("2006-01-02 15:04:05"), getLogString(lv), fileName, funcName, lineNo, msg)
		if lv >= ERROR {
			// 其实把splitFile写到checkSize更好点，errfileObj传了几遍了，不大好
			if f.checkSize(f.errfileObj) {
				newFile, err := f.splitFile(f.errfileObj) // 日志文件
				if err != nil {
					return
				}
				f.errfileObj = newFile // 到这里才赋值
			}
			// 如果要记录的日志大于等于ERROR级别，还要在err日志文件中记录一遍
			fmt.Fprintf(f.errfileObj, "[%s] [%s] [%s:%s:%d] %s\n", now.Format("2006-01-02 15:04:05"), getLogString(lv), fileName, funcName, lineNo, msg)
		}
	}
}

// Debug ...
func (f *FileLogger) Debug(format string, a ...interface{}) {
	f.log(DEBUG, format, a...)
}

// Info ...
func (f *FileLogger) Info(format string, a ...interface{}) {
	f.log(INFO, format, a...)
}

// Warning ...
func (f *FileLogger) Warning(format string, a ...interface{}) {
	f.log(WARNING, format, a...)
}

// Error ...
func (f *FileLogger) Error(format string, a ...interface{}) {
	f.log(ERROR, format, a...)
}

// Fatal ...
func (f *FileLogger) Fatal(format string, a ...interface{}) {
	f.log(FATAL, format, a...)
}

// Close ...
func (f *FileLogger) Close() {
	f.fileObj.Close()
	f.errfileObj.Close()
}
