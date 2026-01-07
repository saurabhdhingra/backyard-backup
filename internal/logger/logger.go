package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func Init(logLevel, logFile string) error {
	var out *os.File
	var err error

	if logFile != "" {
		if err := os.MkdirAll(filepath.Dir(logFile), 0755); err != nil {
			return err
		}
		out, err = os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
	} else {
		out = os.Stdout
	}

	InfoLogger = log.New(out, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(out, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

func Info(format string, v ...interface{}) {
	if InfoLogger != nil {
		InfoLogger.Printf(format, v...)
	} else {
		// Fallback
		fmt.Printf("INFO: "+format+"\n", v...)
	}
}

func Error(format string, v ...interface{}) {
	if ErrorLogger != nil {
		ErrorLogger.Printf(format, v...)
	} else {
		fmt.Printf("ERROR: "+format+"\n", v...)
	}
}
