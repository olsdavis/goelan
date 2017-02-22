package log

/*

Because nobody has been able to develop a proper library for logging in Go up until
the moment I am writing this comment (20:21 02/21/2017), I had to write the following
code, which is quite awful - it's true; although, it works.
Maybe one day I will make one; I will keep you up to date.

*/

import (
	"../util"
	"fmt"
	"io"
	"log"
	"os"
	"time"
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
		fmt.Println(os.Stderr, "Could not create log file:", err)
	}

	writer := io.MultiWriter(file, os.Stdout)
	DebugLogger = log.New(os.Stdout, "[DEBU] ", log.Ltime|log.Lshortfile)
	InfoLogger = log.New(writer, "[INFO] ", log.Ltime)
	WarnLogger = log.New(writer, "[WARN] ", log.Ltime)
	ErrorLogger = log.New(writer, "[ERRO] ", log.Ltime|log.Llongfile)
	FatalLogger = log.New(writer, "[FATA] ", log.Ltime)

	initialized = true
}

func Debug(message ...interface{}) {
	DebugLogger.Println(message)
}

func Info(message ...interface{}) {
	InfoLogger.Println(message)
}

func Warn(message ...interface{}) {
	WarnLogger.Println(message)
}

func Error(message ...interface{}) {
	ErrorLogger.Println(message)
}

func Fatal(message ...interface{}) {
	FatalLogger.Println(message)
}
