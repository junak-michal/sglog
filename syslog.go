package sglog

import (
	"fmt"
	"log/syslog"
)

// SyslogBackend represents a logging backend that sends all messages to a
// syslog server - either local or remote.
type SyslogBackend struct {
	writer *syslog.Writer
}

// LocalSyslogServer creates a new SyslogBackend that logs messages into a local syslog server under given facility
// and uses arg[0] as a tag (process name).
// See https://tools.ietf.org/html/rfc5424#section-6.2.1 for list of facilities and their meaning.
func LocalSyslogBackend(facility syslog.Priority) (result *SyslogBackend) {
	result = new(SyslogBackend)
	result.writer, _ = syslog.New(facility, "")
	return
}

func (sba *SyslogBackend) Log(entry *LogEntry) (err error) {
	msg := fmt.Sprintf(syslogLogFmt, entry.PkgPath, entry.File, entry.Line, entry.Message)
	switch entry.Level {
	default:
		fallthrough
	case Debug:
		err = sba.writer.Debug(msg)
	case Info:
		err = sba.writer.Info(msg)
	case Warning:
		err = sba.writer.Warning(msg)
	}
	return
}

const syslogLogFmt = "%s [%s:%d] %s"
