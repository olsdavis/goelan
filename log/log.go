package log

/*

Because nobody has been able to develop a proper library for logging in Go up until
the moment I was writing this comment (20:21 02/21/2017), I had to write the following
code, which is quite awful - it's true; it works though.
Maybe one day I will make one; I will keep you up to date.

*/

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/olsdavis/goelan/util"
)

var (
	initialized = false
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	WarnLogger  *log.Logger
	ErrorLogger *log.Logger
	FatalLogger *log.Logger
)

// Initializes the logger.
func Init() {
	if initialized {
		return
	}

	logFileName := fmt.Sprintf("logs/%v_Server.log", time.Now().Format("2006-01-02-_15-04-05"))

	if val, _ := util.Exists("logs/"); !val {
		os.Mkdir("logs", os.ModePerm)
	}

	if val, _ := util.Exists(logFileName); !val {
		os.Create(logFileName)
	}

	file, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not create log file:", err)
	}

	writer := io.MultiWriter(file, os.Stdout)
	DebugLogger = log.New(os.Stdout, "[DEBU] ", log.Ltime|log.Lshortfile)
	InfoLogger = log.New(writer, "[INFO] ", log.Ltime)
	WarnLogger = log.New(writer, "[WARN] ", log.Ltime)
	ErrorLogger = log.New(writer, "[ERRO] ", log.Ltime|log.Llongfile)
	FatalLogger = log.New(writer, "[FATA] ", log.Ltime)

	initialized = true
}

// Debug prints a debug message to the debugging logger.
func Debug(message ...interface{}) {
	DebugLogger.Println(message)
}

// Info prints an info message to the info logger.
func Info(message ...interface{}) {
	InfoLogger.Println(message)
}

// Warn prints a warning message to the warning logger.
func Warn(message ...interface{}) {
	WarnLogger.Println(message)
}

// Error prints an error message to the error logger.
func Error(message ...interface{}) {
	ErrorLogger.Println(message)
}

// Fatal prints a fatal message to the fatal logger.
func Fatal(message ...interface{}) {
	FatalLogger.Println(message)
}
