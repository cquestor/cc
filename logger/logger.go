package logger

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// ILogger 日志记录器接口
type ILogger interface {
	Print(v ...any)
	Println(v ...any)
	Printf(fotmat string, v ...any)
}

// CLogger 日志记录器
type CLogger struct {
	ConsoleOut io.Writer
	FileOut    io.Writer
	Prefix     string
	Color      TypeColor
	lock       *sync.Mutex
}

// TypeColor 颜色
type TypeColor int

const (
	ColorRed TypeColor = iota + 31
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
	ColorBlack
)

// NewLogger 构造日志记录器
func NewLogger(console io.Writer, file io.Writer, prefix string, color TypeColor) *CLogger {
	return &CLogger{
		ConsoleOut: console,
		FileOut:    file,
		Prefix:     prefix,
		Color:      color,
		lock:       &sync.Mutex{},
	}
}

// Style 设置颜色字体
func Style(color TypeColor, v ...any) string {
	return setColor(color, fmt.Sprint(v...))
}

// SetConsoleOut 设置控制台输出
func (logger *CLogger) SetConsoleOut(console io.Writer) {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	logger.ConsoleOut = console
}

// SetFileOut 设置文件输出
func (logger *CLogger) SetFileOut(file io.Writer) {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	logger.FileOut = file
}

// Write 输出日志
func (logger *CLogger) Write(v []byte) {
	if logger.ConsoleOut != nil {
		logger.ConsoleOut.Write(v)
	}
	if logger.FileOut != nil {
		logger.FileOut.Write(v)
	}
}

// Print 输出
func (logger *CLogger) Print(v ...any) {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	result := fmt.Sprint(v...)
	logger.Write([]byte(logger.appendOut(result)))
}

// Println 行输出
func (logger *CLogger) Println(v ...any) {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	result := fmt.Sprintln(v...)
	logger.Write([]byte(logger.appendOut(result)))
}

// Printf 格式化输出
func (logger *CLogger) Printf(format string, v ...any) {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	result := fmt.Sprintf(format, v...)
	logger.Write([]byte(logger.appendOut(result)))
}

// appendOut 添加输出信息
func (logger *CLogger) appendOut(v string) string {
	return logger.setColor(logger.setPrefix(logger.setTime(v)))
}

// setPrefix 设置前缀
func (logger *CLogger) setPrefix(v string) string {
	return logger.Prefix + " " + v
}

// setTime 设置时间
func (logger *CLogger) setTime(v string) string {
	return time.Now().Format("2006-01-02 15:04:05") + " " + v
}

// setColor 设置颜色
func (logger *CLogger) setColor(v string) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", logger.Color, v)
}
