package filedriller

import (
	"log"
	"os"
)

var (
	// WarningLogger writes warnings to a log file
	WarningLogger *log.Logger
	// InfoLogger writes info to a log file
	InfoLogger *log.Logger
	// ErrorLogger writes warnings to a log file
	ErrorLogger *log.Logger
)

// CreateLogger creates a custom logger
func CreateLogger(logFile string) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	//InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
