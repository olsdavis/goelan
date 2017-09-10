package justlog

import "fmt"

var (
	defaultFormatter Formatter = NewFormatter("%{MESSAGE}")
)

// A logger struct. You can initialize one with New.
type Logger struct {
	formatters []Formatter
	handlers   []Handler
}

// Represents a log. Contains log's message and log's level.
type Log struct {
	message string
	level   LogLevel
}

// String returns the message of the log. (Used for the "fmt" package.)
func (log *Log) String() string {
	return log.message
}

// New creates and returns a new logger without any handler.
func New() *Logger {
	return NewWithHandlers()
}

// NewWithHandlers creates and returns a new logger with the given handlers.
// If you do not give any handler, you can also use New.
func NewWithHandlers(handlers ...Handler) *Logger {
	h := make([]Handler, len(handlers))
	for i, handler := range handlers {
		h[i] = handler
	}
	return &Logger{
		formatters: []Formatter{defaultFormatter},
		handlers: h,
	}
}

// SetFormatters sets the formatters of the current logger.
func (log *Logger) SetFormatters(formatters ...Formatter) *Logger {
	if len(formatters) >= 0 {
		log.formatters = formatters
	}
	return log
}

// Debug sends a debug message to the current Logger.
func (log *Logger) Debug(messages ...interface{}) {
	log.log(DebugLevel, messages)
}

// Info sends a informative message to the current logger.
func (log *Logger) Info(messages ...interface{}) {
	log.log(InfoLevel, messages)
}

// Warn sends a warning message to the current logger.
func (log *Logger) Warn(messages ...interface{}) {
	log.log(WarnLevel, messages)
}

// Error sends an error message to the current logger.
func (log *Logger) Error(messages ...interface{}) {
	log.log(ErrorLevel, messages)
}

// Fatal sends a fatal message to the current logger.
func (log *Logger) Fatal(messages ...interface{}) {
	log.log(FatalLevel, messages)
}

func (log *Logger) log(level LogLevel, messages []interface{}) {
	if len(messages) == 0 {
		return
	}

	var message string

	if len(log.formatters) == 0 {
		// Default formatter
		message = defaultFormatter.Format(2, level, messages)
	} else {
		// Taking the first formatter which handles the current level
		// in order to avoid conflicts.
		FORMAT: for _, formatter := range log.formatters {
			lvls := formatter.GetLevels()
			if lvls == nil {
				lvls = AllLevels
			}
			for _, lvl := range lvls {
				if lvl == level {
					message = formatter.Format(5, level, messages)
					break FORMAT
				}
			}
		}
	}

	l := Log{
		level: level,
		message: message,
	}

	if len(log.handlers) == 0 {
		// Default output
		fmt.Println(l.message)
	} else {
		for _, handler := range log.handlers {
			lvls := handler.GetLevels()
			if lvls == nil {
				lvls = AllLevels
			}
			for _, lvl := range lvls {
				if lvl == level {
					handler.Handle(l)
					break // take the first handler that matches the level
				}
			}
		}
	}
}
