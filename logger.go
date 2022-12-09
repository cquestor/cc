package cc

import (
	"log"
	"os"
)

var (
	infoLogger  = log.New(os.Stdout, "[\033[1;34mINFO\033[0m]", log.LstdFlags|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "[\033[1;31mERROR\033[0m]", log.LstdFlags|log.Lshortfile)
	dbLogger    = log.New(os.Stdout, "[\033[1;33mSQL\033[0m]", log.LstdFlags|log.Lshortfile)
)

var (
	Info   = infoLogger.Println
	Infof  = infoLogger.Printf
	Error  = errorLogger.Println
	Errorf = errorLogger.Printf
	dbLog  = dbLogger.Println
)
