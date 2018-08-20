package sglog

// Level is an "enum" of log levels.
type Level int

const (
	// Debug is the lowest log level. It should be used for debugging messages that are probably
	// too detailed for production logs.
	Debug Level = iota + 1
	// Info is a log level for messages about standard system run. These messages shouldn't flood the log
	// which might make Info a nice default log level for logs in production.
	Info
	// Warning is a log level for things that went wrong and were somehow handled by the application.
	// These messages should definitely be logged in production.
	Warning
)

func (lvl Level) isValid() bool {
	return lvl >= Debug && lvl <= Warning
}

// LogEntry represents information about single log message gathered by the frontend part of the sglog.
type LogEntry struct {
	Level   Level
	PkgPath string
	File    string
	Line    int
	Message string
}
