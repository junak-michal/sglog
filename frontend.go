package sglog

import (
	"fmt"
	"path"
	"runtime"
)

// Logger represents an object that can be used for logging by single package.
// The instances of Logger should be acquired by calling GetLogger function.
type Logger struct {
	pkgPath   string
	pathStrip string
	level     Level
}

// GetLogger gets a Logger instance registered under pkgPath key. Even though the pkgPath key can be
// any string, it is strongly encouraged to use package path as the key so other applications
// can access the Logger and set its level etc.
// If there is not yet a Logger instance registered under pkgPath key, new one shall be created.
func GetLogger(pkgPath string) *Logger {
	return loggerFromGlobalMap(pkgPath)
}

// SetLevel sets a logging level of the Logger to given level.
// Messages with lower priority than logging level of the logger shall
// not be logged.
func (logger *Logger) SetLevel(level Level) {
	if level.isValid() {
		logger.level = level
	}
}

// Debug logs a message under a Debug log level.
// Arguments are handled in the manner of fmt.Printf.
func (logger *Logger) Debug(format string, a ...interface{}) {
	logger.log(Debug, format, a...)
}

// Info logs a message under an Info log level.
// Arguments are handled in the manner of fmt.Printf.
func (logger *Logger) Info(format string, a ...interface{}) {
	logger.log(Info, format, a...)
}

// Warning logs a message under a Warning log level.
// Arguments are handled in the manner of fmt.Printf.
func (logger *Logger) Warning(format string, a ...interface{}) {
	logger.log(Warning, format, a...)
}

func newLogger(pkgPath string) (result *Logger) {
	result = new(Logger)
	result.pkgPath = pkgPath
	result.pathStrip = stripPkgPath(pkgPath)
	result.level = Debug
	return
}

func stripPkgPath(pkgPath string) string {
	rest, last := path.Split(pkgPath)
	if rest == "" || rest == "/" {
		return last
	}
	nextToLast := path.Base(rest)
	return fmt.Sprintf("%s/%s", nextToLast, last)
}

// callerSkipFromLog holds number of frames that we need to ascend from the log method to the
// actual code that wanted to log a message.
const callerSkipFromLog = 2

func (logger *Logger) log(level Level, format string, a ...interface{}) {
	if level < logger.level {
		return
	}
	// If we do not recover the information then they won't be part of the LogEntry - no need to handle it
	// in any other way.
	_, filePath, line, _ := runtime.Caller(callerSkipFromLog)
	_, fileName := path.Split(filePath)
	entry := LogEntry{
		Level:   level,
		PkgPath: logger.pathStrip,
		File:    fileName,
		Line:    line,
		Message: fmt.Sprintf(format, a...),
	}
	passEntryToBackend(&entry)
}

func passEntryToBackend(entry *LogEntry) {
	backendInUse := getGlobalBackend()
	err := backendInUse.Log(entry)
	if err != nil && backendInUse != defaultBackend {
		// Ignoring possible error.
		// Default should by as safe to use as possible.
		defaultBackend.Log(entry)
	}
}

var internalLog = GetLogger("github.com/junak-michal/sglog")
