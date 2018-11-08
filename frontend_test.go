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
	logger := GetLogger("testlogger")
	for _, inData := range setLevelTestData {
		inData.setUpLogger(logger)
		inData.logMessages(logger)
		if !inData.verify() {
			t.Errorf("Messages not logged as expected. Test  input: %v", inData)
		}
	}
}
