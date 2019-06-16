package sglog

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// currentConfigVersion must match the "version" field in the YAML config, otherwise
// the config won't be loaded.
// The number should be increased every time config document changes its structure.
const currentConfigVersion = 1

// ConfigFile holds path to the sglog YAML config. The default value is sglog.yaml, but for many users
// it will be better idea to change the value during compilation with the -X flag.
// Example: go build -ldflags "-X github.com/junak-michal/sglog.ConfigFile=/etc/myservice/logging.yaml".
var ConfigFile = "sglog.yaml"

func init() {
	confErr := loadConfig(ConfigFile)
	if confErr != nil {
		internalLog.Warning("Unable to load config file, proceeding with defaults. Error: %v", confErr)
	}
}

func loadConfig(filename string) (err error) {
	var confYAML confDoc
	confYAML, err = loadYAMLFile(filename)
	if err != nil {
		return
	}
	if confYAML.Version != currentConfigVersion {
		return fmt.Errorf("wrong YAML config version, expected %d, got %d", currentConfigVersion, confYAML.Version)
	}
	applyYAMLConfig(confYAML)
	return
}

func loadYAMLFile(filename string) (confYAML confDoc, err error) {
	var confBytes []byte
	confBytes, err = ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	err = yaml.UnmarshalStrict(confBytes, &confYAML)
	return
}

func applyYAMLConfig(confYAML confDoc) {
	if confYAML.DefaultConf != "" {
		defaultLogLevel = confYAML.DefaultConf.toLevel()
	}
	for name, lvl := range confYAML.LoggerConfs {
		logger := loggerFromGlobalMap(name)
		logger.level = lvl.toLevel()
	}
}
