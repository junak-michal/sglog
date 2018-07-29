package sglog

import (
	"sync"
)

var (
	loggers    = make(map[string]*Logger)
	loggersMut sync.Mutex
)

func loggerFromGlobalMap(pkgPath string) (result *Logger) {
	loggersMut.Lock()
	defer loggersMut.Unlock()
	var contains bool
	result, contains = loggers[pkgPath]
	if !contains {
		result = newLogger(pkgPath)
		loggers[pkgPath] = result
	}
	return
}

var (
	// TODO: default should not be nil
	backend    Backend
	backendMut sync.Mutex
)

func getGlobalBackend() Backend {
	// TODO: try with RW lock
	backendMut.Lock()
	defer backendMut.Unlock()
	return backend
}

func setGlobalBackend(newBackend Backend) {
	if backend == nil {
		// TODO: error
		return
	}
	backendMut.Lock()
	defer backendMut.Unlock()
	backend = newBackend
}
