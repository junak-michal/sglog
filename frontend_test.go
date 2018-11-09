package sglog

import (
	"testing"
)

type mockRecvBackend struct {
	received []*LogEntry
}

func (mrb *mockRecvBackend) Log(entry *LogEntry) error {
	mrb.received = append(mrb.received, entry)
	return nil
}

func (mrb *mockRecvBackend) entryWithLvlReceived(lvl Level) bool {
	for _, entry := range mrb.received {
		if entry.Level == lvl {
			return true
		}
	}
	return false
}

type lvlSet struct {
	lvl          Level
	dbgExpected  bool
	infoExpected bool
	warnExpected bool
	backend      *mockRecvBackend
}

func newLvlSet(lvl Level, dbgExpected, infoExpected, warnExpected bool) (result *lvlSet) {
	result = new(lvlSet)
	result.lvl = lvl
	result.dbgExpected = dbgExpected
	result.infoExpected = infoExpected
	result.warnExpected = warnExpected
	result.backend = new(mockRecvBackend)
	return
}

const (
	dbgTestMsg  = "A DEBUG MESSAGE"
	infoTestMsg = "AN INFO MESSAGE"
	warnTestMsg = "A WARNING MESSAGE"
)

func (this *lvlSet) setUpLogger(logger *Logger) {
	UseBackend(this.backend)
	logger.SetLevel(this.lvl)
}

func (this *lvlSet) logMessages(logger *Logger) {
	logger.Debug(dbgTestMsg)
	logger.Info(infoTestMsg)
	logger.Warning(warnTestMsg)
}

func (this *lvlSet) verify() (result bool) {
	result = this.dbgExpected == this.backend.entryWithLvlReceived(Debug) &&
		this.infoExpected == this.backend.entryWithLvlReceived(Info) &&
		this.warnExpected == this.backend.entryWithLvlReceived(Warning)
	return
}

var setLevelTestData = []*lvlSet{
	newLvlSet(Debug, true, true, true),
	newLvlSet(Info, false, true, true),
	newLvlSet(Warning, false, false, true),
}

func TestSetLevel(t *testing.T) {
	logger := GetLogger("TestSetLevel_logger")
	for _, inData := range setLevelTestData {
		inData.setUpLogger(logger)
		inData.logMessages(logger)
		if !inData.verify() {
			t.Errorf("Messages not logged as expected. Test  input: %v", inData)
		}
	}
}

func TestFileName(t *testing.T) {
	backend := new(mockRecvBackend)
	UseBackend(backend)
	logger := GetLogger("TestFileName_logger")
	logger.Debug("Log message %d with %d arguments, %s", 1, 3, "BOO!")
	logger.Info("Log message without arguments.")
	logger.Warning("%s", "Log message as an argument.")
	if len(backend.received) != 3 {
		t.Errorf("Received too many log entries.")
	}
	for _, entry := range backend.received {
		if entry.File != "frontend_test.go" {
			t.Errorf("Received log entry with wrong file name.")
		}
	}
}

var testPkgPathData = []struct {
	fullPkgPath     string
	expectedPkgPath string
}{
	{fullPkgPath: "bitbucket.org/myteamname/project/subpkg", expectedPkgPath: "project/subpkg"},
	{fullPkgPath: "main", expectedPkgPath: "main"},
	{fullPkgPath: "cool/pkg", expectedPkgPath: "cool/pkg"},
	{fullPkgPath: "a", expectedPkgPath: "a"},
	{fullPkgPath: "/mypkg", expectedPkgPath: "mypkg"},
	{fullPkgPath: "myproj/mypkg/", expectedPkgPath: "mypkg/"},
}

func TestPkgPath(t *testing.T) {
	backend := new(mockRecvBackend)
	UseBackend(backend)
	for _, input := range testPkgPathData {
		logger := GetLogger(input.fullPkgPath)
		logger.Debug("Doesn't matter what this message contains.")
		// This should never evaluate to -1
		lastEntry := backend.received[len(backend.received)-1]
		if lastEntry.PkgPath != input.expectedPkgPath {
			t.Errorf("Expected pkgPath '%s' but reveived '%s'.", input.expectedPkgPath, lastEntry.PkgPath)
		}
	}
}
