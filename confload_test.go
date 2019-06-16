package sglog

import (
	"fmt"
	"testing"
)

func TestYAMLExample1(t *testing.T) {
	stateBefore := getCurrentState()
	loadErr := loadConfig("confexamples/confexample1.yaml")
	if loadErr != nil {
		t.Error(loadErr)
	}
	if defaultLogLevel != Debug {
		t.Error("Wrong default log level: ", defaultLogLevel)
	}
	expected := []expectedPair{
		{"bitbucket.org/project/package", Warning},
		{"github.com/project/package", Info},
		{"justaname", Debug},
	}
	if err := checkChanges(stateBefore, expected); err != nil {
		t.Error(err)
	}
	applyState(stateBefore)
}

func TestYAMLExample2(t *testing.T) {
	stateBefore := getCurrentState()
	loadErr := loadConfig("confexamples/confexample2.yaml")
	if loadErr != nil {
		t.Error(loadErr)
	}
	if defaultLogLevel != Info {
		t.Error("Wrong default log level: ", defaultLogLevel)
	}
	expected := []expectedPair{
		{"bitbucket.org/mysw/currentlydebugging", Debug},
		{"bitbucket.org/mysw/noisyinfo", Warning},
		{"bitbucket.org/somelibrary/whynotrepeatthedefault", Info},
		{"main", Debug},
	}
	if err := checkChanges(stateBefore, expected); err != nil {
		t.Error(err)
	}
	applyState(stateBefore)
}

func checkChanges(previousState *sglogState, expectedChanges []expectedPair) error {
	currentState := getCurrentState()
	for _, pair := range expectedChanges {
		actualLvl := currentState.levelsByName[pair.name]
		if actualLvl != pair.lvl {
			return fmt.Errorf("expected level %d for package %s, but got level %d", pair.lvl, pair.name, actualLvl)
		}
		delete(currentState.levelsByName, pair.name)
	}
	for name, previousLvl := range previousState.levelsByName {
		actualLvl := currentState.levelsByName[name]
		if actualLvl != previousLvl {
			return fmt.Errorf("expected unchanged level %d for package %s, but got level %d", previousLvl, name, actualLvl)
		}
		delete(currentState.levelsByName, name)
	}
	remainder := len(currentState.levelsByName)
	if remainder != 0 {
		return fmt.Errorf("there are %d unexpected loggers", remainder)
	}
	return nil
}

func getCurrentState() *sglogState {
	currentLevels := make(map[string]Level)
	for name, logger := range loggers {
		currentLevels[name] = logger.level
	}
	return &sglogState{defaultLogLevel: defaultLogLevel, levelsByName: currentLevels}
}

func applyState(state *sglogState) {
	defaultLogLevel = state.defaultLogLevel
	for name, logger := range loggers {
		if _, contains := state.levelsByName[name]; contains {
			logger.level = state.levelsByName[name]
		} else {
			delete(loggers, name)
		}
	}
}

type sglogState struct {
	defaultLogLevel Level
	levelsByName    map[string]Level
}

type expectedPair struct {
	name string
	lvl  Level
}
