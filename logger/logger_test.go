package logger_test

import (
	"os"
	"testing"
	"time"

	"github.com/cquestor/cc/logger"
)

func TestLogger(t *testing.T) {
	testLog := logger.NewLogger(os.Stdout, nil, "[TEST]", logger.ColorBlue)
	t.Run("print", func(t *testing.T) {
		testLog.Print("print", logger.Style(logger.ColorRed, logger.StyleBold, " ???\n"))
	})
	t.Run("println", func(t *testing.T) {
		testLog.Println("println", logger.Style(logger.ColorYellow, logger.StyleItalic, "???"))
	})
	t.Run("printf", func(t *testing.T) {
		testLog.Printf("printf %s\n", logger.Style(logger.ColorGreen, logger.StyleUnderline, "???"))
	})
	t.Run("spin", func(t *testing.T) {
		index := 0
		for index < 50 {
			testLog.Spin(logger.ColorGreen, logger.StyleBold, "Loading...")
			index++
			time.Sleep(time.Millisecond * 80)
		}
	})
}
