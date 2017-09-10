package justlog

// LogLevel represents the importance of the log.
type LogLevel int

const (
	// Debugging level.
	DebugLevel LogLevel = iota
	// Informative level.
	InfoLevel
	// Warning level.
	WarnLevel
	// Error level.
	ErrorLevel
	// Fatal level.
	FatalLevel
)

var (
	levelNames = []string{
		"DEBUG",
		"INFO",
		"WARN",
		"ERROR",
		"FATAL",
	}

	AllLevels = []LogLevel{
		DebugLevel,
		InfoLevel,
		WarnLevel,
		ErrorLevel,
		FatalLevel,
	}
)

// LeveledLogging is an interface that any object in the logging process can implement.
// It is implemented when the object has to be applied only to certain levels of messages.
type LeveledLogging interface {
	GetLevels() []LogLevel
}
