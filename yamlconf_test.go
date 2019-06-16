package sglog

import (
	"gopkg.in/yaml.v2"
	"reflect"
	"testing"
)

func TestUnmarshallOK(t *testing.T) {
	for _, testIn := range okExamples {
		var conf confDoc
		err := yaml.UnmarshalStrict([]byte(testIn.input), &conf)
		if err != nil {
			t.Error("Unexpected error when unmarshalling.", err)
		}
		if !reflect.DeepEqual(conf, testIn.expected) {
			t.Errorf("Not unmarshalled as expected.\nExpected: %v\nGot: %v\n", testIn.expected, conf)
		}
	}
}

func TestUnmarshallFail(t *testing.T) {
	for _, testIn := range failExamples {
		var conf confDoc
		err := yaml.UnmarshalStrict([]byte(testIn), &conf)
		if err == nil {
			t.Error("Erroneous config accepted:", testIn)
		}
	}
}

func TestToLevel(t *testing.T) {
	for _, testIn := range toLevelInputs {
		actualLvl := testIn.confStr.toLevel()
		if actualLvl != testIn.expectedLvl {
			t.Errorf("Expected level %v, got %v.", testIn.expectedLvl, actualLvl)
		}
	}
}

var (
	exampleOK1 = "" +
		"default: Debug\n" +
		"loggers:\n" +
		"  bitbucket.org/project/package: Warning\n" +
		"  github.com/project/package: Warning\n"
	exampleOK2 = "" +
		"default: Warning\n" +
		"loggers:\n" +
		"  bitbucket.org/project/package: Info\n" +
		"  justaname: Warning\n"
	exampleOK3 = "" +
		"version: 1\n" +
		"default: Warning\n" +
		"loggers:\n" +
		"  dir1/dir2/dir3: Info\n" +
		"  Debug: Debug\n" +
		"  justaname: Warning\n"
	noDefaultOK = "" +
		"loggers:\n" +
		"  bitbucket.org/project/package: Warning\n"
	defaultOnlyOK = "" +
		"version: 1\n" +
		"default: Warning\n"
	okExamples = []struct {
		input    string
		expected confDoc
	}{
		{exampleOK1, confDoc{Version: 0, DefaultConf: "Debug", LoggerConfs: map[string]levelConf{
			"bitbucket.org/project/package": "Warning",
			"github.com/project/package":    "Warning",
		}}},
		{exampleOK2, confDoc{Version: 0, DefaultConf: "Warning", LoggerConfs: map[string]levelConf{
			"bitbucket.org/project/package": "Info",
			"justaname":                     "Warning",
		}}},
		{exampleOK3, confDoc{Version: 1, DefaultConf: "Warning", LoggerConfs: map[string]levelConf{
			"dir1/dir2/dir3": "Info",
			"Debug":          "Debug",
			"justaname":      "Warning",
		}}},
		{noDefaultOK, confDoc{Version: 0, DefaultConf: "", LoggerConfs: map[string]levelConf{
			"bitbucket.org/project/package": "Warning",
		}}},
		{defaultOnlyOK, confDoc{Version: 1, DefaultConf: "Warning"}},
	}
)

var (
	failVersionNAN = "" +
		"version: NotANumber\n" +
		"default: Debug\n" +
		"loggers:\n" +
		"  bitbucket.org/project/package: Warning\n"
	failLoggersString = "" +
		"version: 1\n" +
		"default: Debug\n" +
		"loggers: Debug\n"
	failExamples = []string{
		failVersionNAN, failLoggersString,
	}
)

var toLevelInputs = []struct {
	confStr     levelConf
	expectedLvl Level
}{
	{"DeBUG", Debug},
	{"Warning", Warning},
	{"INFO", Info},
	{"TYPO", Debug},
	{"Warn", Debug},  // "Warn" is also a typo.
	{"Error", Debug}, // "Error" level doesn't exist.
	{"Really anything", Debug},
}
