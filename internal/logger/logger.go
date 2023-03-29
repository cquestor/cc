package logger

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// TODO LogFatal LogFatalf os.Exit(1)

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

const (
	StyleNormal    TypeStyle = 0
	StyleBold      TypeStyle = 1
	StyleItalic    TypeStyle = 3
	StyleUnderline TypeStyle = 4
	StyleInverse   TypeStyle = 7
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
	Style      TypeStyle
	lock       *sync.Mutex
}

// TypeColor 字体颜色
type TypeColor int

// TypeStyle 字体样式
type TypeStyle int

var spinners = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

var spinner = getSpinner()

// NewLogger 构造日志记录器
func NewLogger(console io.Writer, file io.Writer, prefix string, color TypeColor, style ...TypeStyle) *CLogger {
	if len(style) < 1 {
		style = append(style, StyleNormal)
	}
	return &CLogger{
		ConsoleOut: console,
		FileOut:    file,
		Prefix:     prefix,
		Color:      color,
		Style:      style[0],
		lock:       &sync.Mutex{},
	}
}

// Style 设置颜色字体
func Style(color TypeColor, style TypeStyle, v ...any) string {
	return setStyle(color, style, fmt.Sprint(v...))
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
	result := logger.appendOut(fmt.Sprint(v...))
	logger.Write([]byte(setStyle(logger.Color, logger.Style, result)))
}

// Println 行输出
func (logger *CLogger) Println(v ...any) {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	result := logger.appendOut(fmt.Sprintln(v...))
	logger.Write([]byte(setStyle(logger.Color, logger.Style, result)))
}

// Printf 格式化输出
func (logger *CLogger) Printf(format string, v ...any) {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	result := logger.appendOut(fmt.Sprintf(format, v...))
	logger.Write([]byte(setStyle(logger.Color, logger.Style, result)))
}

// Spin 加载动画
func (logger *CLogger) Spin(color TypeColor, style TypeStyle, message string) {
	logger.lock.Lock()
	defer logger.lock.Unlock()
	logger.Write([]byte(fmt.Sprintf("\r\033[%d;%dm%s %s\033[0m", style, color, spinner(), message)))
}

// appendOut 添加输出信息
func (logger *CLogger) appendOut(v string) string {
	return logger.setPrefix(logger.setTime(v))
}

// setPrefix 设置前缀
func (logger *CLogger) setPrefix(v string) string {
	return logger.Prefix + " " + v
}

// setTime 设置时间
func (logger *CLogger) setTime(v string) string {
	return time.Now().Format("2006-01-02 15:04:05") + " " + v
}
