package main

// 适配器模式

// 定义目标接口
type Logger interface {
	WriteLog(message string)
}

// 定义源接口
type LoggerWriter interface {
	Write(message string)
}

// 定义适配器
type LoggerAdapter struct {
	writer LoggerWriter
}

// 实现目标接口
func (l *LoggerAdapter) WriteLog(message string) {
	// 调用源接口的方法来写入日志
	l.writer.Write(message)
}

type ExistingLogWriter struct{}

// 实现源接口
func (l *ExistingLogWriter) Write(message string) {
	// 已有日志库的写入逻辑
	println(message)
}

func main() {
	// 适配器模式
	l := &LoggerAdapter{writer: &ExistingLogWriter{}}
	l.WriteLog("Hello, World!")
}
