package logging

import (
	"log"
	"os"
)

var (
	Error *log.Logger
	Warn  *log.Logger
	Info  *log.Logger
)

func init() {
	flags := log.Ldate | log.Ltime | log.Lshortfile

	Error = log.New(os.Stdout, "FATAL: ", flags)
	Warn = log.New(os.Stdout, "WARN: ", flags)
	Info = log.New(os.Stdout, "INFO: ", flags)
}
