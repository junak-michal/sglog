package sglog

import (
	"fmt"
	"os"
)

// StderrBackend represents a logging backend that sends all messages to stderr.
// It is the default and "backup" backend.
type StderrBackend struct{}

func (stdeb *StderrBackend) Log(entry *LogEntry) error {
	msg := fmt.Sprintf(stdErrLogFmt, lvlToString(entry.Level), entry.PkgPath, entry.File, entry.Line, entry.Message)
	fmt.Fprintln(os.Stderr, msg) // ignoring potential error
	return nil                   // never fails
}

const stdErrLogFmt = "%s %s\t[%s:%d]\t%s"

func lvlToString(level Level) (result string) {
	switch level {
	default:
		fallthrough
	case Debug:
		result = "DBUG"
	case Info:
		result = "INFO"
	case Warning:
		result = "WARN"
	}
	return
}
