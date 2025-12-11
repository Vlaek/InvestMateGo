package logger

import (
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func InfoLog(format string, v ...interface{}) {
	Info.Printf(format, v...)
}

func ErrorLog(format string, v ...interface{}) {
	Error.Printf(format, v...)
}
