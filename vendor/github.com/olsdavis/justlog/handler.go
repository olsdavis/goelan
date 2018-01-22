package justlog

import (
	"os"
	"fmt"
	"bytes"
	"regexp"
	"strings"
	"time"
	"runtime"
	"path/filepath"
	"strconv"
)

const (
	declareArg = '%'
	openArg = '{'
	closeArg = '}'
)

var pid string = strconv.Itoa(os.Getpid())

var simpleMessageChain appender = func(calldepth int, lvl LogLevel, messages []interface{}) string {
	return strings.Trim(fmt.Sprintf(strings.TrimSpace(strings.Repeat("%v ", len(messages))), messages...), "[]")
}

var simpleTimeChain appender = func(calldepth int, lvl LogLevel, messages []interface{}) string {
	return time.Now().Format("15:04:05")
}

var simpleDateChain appender = func(calldepth int, lvl LogLevel, messages []interface{}) string {
	return time.Now().Format("2006/01/02")
}

var simpleLevelChain appender = func(calldepth int, lvl LogLevel, messages []interface{}) string {
	return levelNames[lvl]
}

var simpleShortCallerChain appender = func(calldepth int, lvl LogLevel, messages []interface{}) string {
	file, line := getCaller(calldepth + 1)
	return fmt.Sprintf("%v:%v", filepath.Base(file), line)
}
var simpleLongCallerChain appender = func(calldepth int, lvl LogLevel, messages []interface{}) string {
	file, line := getCaller(calldepth + 1)
	return fmt.Sprintf("%v:%v", file, line)
}

var simplePidChain appender = func(calldepth int, lvl LogLevel, messages []interface{}) string {
	return pid
}

func getCaller(calldepth int) (string, int) {
	_, file, line, ok := runtime.Caller(calldepth)
	if !ok {
		file = "???"
		line = 0
	}
	return file, line
}

var argumentsRegex *regexp.Regexp = regexp.MustCompile("%\\{[A-Za-z_]+}")

type appender func(int, LogLevel, []interface{}) string

type leveled struct {
	levels []LogLevel
}

// GetLevels returns the level(s) for which this current object should be called.
func (lvl leveled) GetLevels() []LogLevel {
	return lvl.levels
}

type Formatter interface {
	// GetLevels returns the level(s) for which this current Formatter should be called.
	GetLevels() []LogLevel

	// Formats the messages of the log and returns the string.
	// The integer is the number of functions that have been called to get to the formatter (call's depth).
	// The log level is message's level.
	// The interfaces are the different parts of the message.
	Format(int, LogLevel, []interface{}) string
}

type argumentFormatter struct {
	leveled
	chain []appender
}

func (formatter argumentFormatter) Format(calldepth int, level LogLevel, messages []interface{}) string {
	buffer := new(bytes.Buffer)
	for _, c := range formatter.chain {
		if c != nil {
			txt := c(calldepth, level, messages)
			buffer.WriteString(txt)
		}
	}
	return buffer.String()
}

const (
	messageArgument = "%{MESSAGE}"
	levelArgument = "%{LEVEL}"
	timeArgument = "%{TIME}"
	dateArgument = "%{DATE}"
	shortCallerArgument = "%{SHORT_CALLER}"
	longCallerArgument = "%{LONG_CALLER}"
	pidArgument = "%{PID}"
)

// NewFormatter creates a formatter with the given format (which's arguments are replaced by values)
// and levels for which it should apply.
// The current available arguments are:
//
// - "%{MESSAGE}": will be replaced by the concatenation of the messages (or just by the message if there is only
// a singular interface which has been provided);
//
// - "%{LEVEL}": will be replaced by message's log string (eg: logger.Error("such problem wow") will format "ERROR");
//
// - "%{TIME}": will be replaced by the time of when the message is being formatted with the following format: "HH:mm:ss"
// #NoNazi ("HH" is for the hours, "mm" is for the minutes and "ss" is for the seconds);
//
// - "%{DATE}": will be replaced by the date of when the message is being formatted with the following format: "YYYY:MM:DD"
// ("YYYY" is for the year, "MM" is the month and "DD" is the day);
//
// - "%{SHORT_CALLER}": will be replaced by the name and the line of where the logger has been called (eg: the format is
// "%{SHORT_CALLER}: %{MESSAGE}" and at line 4 of my file "paprika.go" I have called logger.Info("Starfoullilai"),
// the message will be: "paprika.go:4: Starfoullilai");
//
// - "%{LONG_CALLER}": will be replaced by the full path of the file and the line of where the logger has been called
// (if we take the same example as above, but the format is LONG_CALLER instead of SHORT_CALLER, the output will then be:
// "/home/your/path/to/the/file/test.go:4: Starfoullilai");
//
// - "%{PID}": will be replaced by the process id.
//
// Known bug: add a trailing space at the end of the format, the last char is removed if you have some text after an argument
func NewFormatter(format string, levels ...LogLevel) Formatter {
	appenders := make([]appender, 0)
	freeIndices := make([]int, 0)

	opened := false
	simple := false
	currentString := ""
	for i := 0; i < len(format); i++ {
		char := format[i]
		if char == declareArg && (len(format) + 1 > i && format[i + 1] == openArg) {
			opened = true
		} else if char == closeArg {
			opened = false
			simple = true
		}

		if simple || len(format) == i + 1 {
			str := currentString
			currentString = ""
			appenders = append(appenders, func(calldepth int, lvl LogLevel, messages []interface{}) string {
				return str
			})
			appenders = append(appenders, nil)
			freeIndices = append(freeIndices, len(appenders) - 1)
			simple = false
			continue
		}

		if !opened {
			currentString += string(char)
		}
	}

	matches := argumentsRegex.FindAllString(format, -1)
	if len(matches) == 0 {
		appenders[freeIndices[0]] = simpleMessageChain
	} else {
		used := 0
		for _, match := range matches {
			switch match {
			case messageArgument:
				appenders[freeIndices[used]] = simpleMessageChain
				used++
			case levelArgument:
				appenders[freeIndices[used]] = simpleLevelChain
				used++
			case timeArgument:
				appenders[freeIndices[used]] = simpleTimeChain
				used++
			case dateArgument:
				appenders[freeIndices[used]] = simpleDateChain
				used++
			case shortCallerArgument:
				appenders[freeIndices[used]] = simpleShortCallerChain
				used++
			case longCallerArgument:
				appenders[freeIndices[used]] = simpleLongCallerChain
				used++
			case pidArgument:
				appenders[freeIndices[used]] = simplePidChain
				used++
			}
		}
	}
	return argumentFormatter{
		chain: appenders,
		leveled: leveled{levels},
	}
}

type Handler interface {
	// GetLevels returns the level(s) for which this current Handler should be called.
	GetLevels() []LogLevel

	// Handle handles the log.
	Handle(Log)
}

/* Default implementations */

/* File implementation */
type fileHandler struct {
	leveled
	path string
	file *os.File
}

// NewFileHandler creates a file Handler, which will output the logs to the given file path.
// path is the path to your file.
// levels are the levels that you want your Handler to handle.
func NewFileHandler(path string, levels ...LogLevel) Handler {
	return &fileHandler{
		leveled: leveled{levels},
		path: path,
		file: nil,
	}
}

func (handler *fileHandler) Handle(log Log) {
	// init file stream
	if handler.file == nil {
		if _, err := os.Stat(handler.path); os.IsNotExist(err) {
			err = os.MkdirAll(filepath.Dir(handler.path), os.ModePerm)
			if err != nil {
				panic(fmt.Sprintf("could not create dir %v, error %v", handler.path, err))
			}
			f, err := os.Create(handler.path)
			if err != nil {
				panic(fmt.Sprintf("could not create log file for path %v, error: %v", handler.path, err))
			}
			handler.file = f
		} else {
			f, err := os.OpenFile(handler.path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil {
				panic(fmt.Sprintf("could not open file for path %v, error: %v", handler.path, err))
			}
			handler.file = f
		}
	}

	if log.message != "" {
		_, err := handler.file.WriteString(log.message + "\n")
		if err != nil {
			panic(err)
		}
	}
}

/* Basic console implementation */
type consoleHandler struct {
	leveled
}

// NewConsoleHandler creates and returns a Handler which prints every log
// message that it has access to (according to its level(s)) to the console
func NewConsoleHandler(levels ...LogLevel) Handler {
	return &consoleHandler{leveled{levels}}
}

func (handler consoleHandler) Handle(log Log) {
	fmt.Println(log.message)
}
