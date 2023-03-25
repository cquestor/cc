package logger_test

import (
	"os"
	"testing"

	"github.com/cquestor/cc/logger"
)

func TestLogger(t *testing.T) {
	testLog := logger.NewLogger(os.Stdout, nil, "[TEST]", logger.ColorBlue)
	t.Run("print", func(t *testing.T) {
		testLog.Print("print", logger.Style(logger.ColorRed, " ???\n"))
	})
	t.Run("println", func(t *testing.T) {
		testLog.Println("println", logger.Style(logger.ColorYellow, "???"))
	})
	t.Run("printf", func(t *testing.T) {
		testLog.Printf("printf %s\n", logger.Style(logger.ColorGreen, "???"))
	})
}
