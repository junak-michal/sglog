package sglog

import (
	"sync"
	"errors"
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
	defaultBackend Backend = new(StderrBackend)
	backend                = defaultBackend
	backendMut     sync.Mutex
)

func getGlobalBackend() Backend {
	// TODO: try with RW lock
	backendMut.Lock()
	defer backendMut.Unlock()
	return backend
}

func setGlobalBackend(newBackend Backend) (err error) {
	if backend == nil {
		err = errors.New("cannot use nil logging backend")
	} else {
		backendMut.Lock()
		defer backendMut.Unlock()
		backend = newBackend
	}
	return
}
