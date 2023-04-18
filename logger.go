package cc

import (
	"os"

	"github.com/cquestor/cc/logger"
)

var (
	infoLogger = logger.NewLogger(os.Stderr, nil, "[INFO]", logger.ColorBlue, logger.StyleBold)
	warnLogger = logger.NewLogger(os.Stderr, nil, "[WARN]", logger.ColorYellow, logger.StyleBold)
	errLogger  = logger.NewLogger(os.Stderr, nil, "[ERROR]", logger.ColorRed, logger.StyleBold)
)

var (
	LogInfo  = infoLogger.Println
	LogInfof = infoLogger.Printf
	LogWarn  = warnLogger.Println
	LogWarnf = warnLogger.Printf
	LogErr   = errLogger.Println
	LogErrf  = errLogger.Printf
)
