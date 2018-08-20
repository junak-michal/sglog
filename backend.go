package sglog

// Backend receives LogEntry instances from the frontend (represented by the Logger type) after they were filtered by
// log level. Backend's only responsibility is to write them somewhere - syslog, file, database, error output, etc.
type Backend interface {
	// Log writes a single LogEntry to the backend returning a non-nil error if the operation was not successful.
	Log(entry *LogEntry) error
}

// UseBackend sets given backend to be used by all loggers.
// If the backend cannot be used, the default one shall be used instead an a warning message shall be logged into it.
func UseBackend(backend Backend) {
	err := setGlobalBackend(backend)
	if err != nil {
		internalLog.Warning("Unable to use given backend, error: %v.", err)
	}
}
