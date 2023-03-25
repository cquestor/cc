package logger

import (
	"fmt"
	"io"
	"os"
)

// GetFileWriter 获取文件流
func GetFileWriter(filename string) (io.Writer, error) {
	if !fileExists(filename) {
		return os.Create(filename)
	}
	return os.OpenFile(filename, os.O_APPEND, 0644)
}

// fileExists 判断文件是否存在
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

// setColor 设置输出颜色
func setColor(color TypeColor, v string) string {
	prefix := fmt.Sprintf("\033[%dm", color)
	suffix := "\033[0m"
	return prefix + v + suffix
}
