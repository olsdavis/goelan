package log

import (
	"fmt"
	"time"
	"github.com/olsdavis/justlog"
)

var (
	initialized = false
	Logger *justlog.Logger
)

// Initializes the logger.
func Init() {
	if initialized {
		return
	}

	Logger = justlog.NewWithHandlers(justlog.NewConsoleHandler(justlog.AllLevels...),
		justlog.NewFileHandler(fmt.Sprintf("logs/%v_Server.log", time.Now().Format("2006-01-02-_15-04-05")),
			justlog.ErrorLevel, justlog.FatalLevel, justlog.WarnLevel, justlog.InfoLevel))
	Logger.SetFormatters(justlog.NewFormatter("[%{LEVEL}] [%{TIME}]: %{MESSAGE}", justlog.AllLevels...))

	initialized = true
}

// Debug prints a debug message to the debugging logger.
func Debug(message ...interface{}) {
	Logger.Debug(message)
}

// Info prints an info message to the info logger.
func Info(message ...interface{}) {
	Logger.Info(message)
}

// Warn prints a warning message to the warning logger.
func Warn(message ...interface{}) {
	Logger.Warn(message)
}

// Error prints an error message to the error logger.
func Error(message ...interface{}) {
	Logger.Error(message)
}

// Fatal prints a fatal message to the fatal logger.
func Fatal(message ...interface{}) {
	Logger.Fatal(message)
}
