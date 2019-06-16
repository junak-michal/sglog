package sglog

import "strings"

// confDoc represents single YAML document that is holding logging configuration.
type confDoc struct {
	Version     int                  `yaml:"version,omitempty"`
	DefaultConf levelConf            `yaml:"default"`
	LoggerConfs map[string]levelConf `yaml:"loggers"`
}

type levelConf string

func (lvlC levelConf) toLevel() Level {
	lowStr := strings.ToLower(string(lvlC))
	switch lowStr {
	case "warning":
		return Warning
	case "info":
		return Info
	default:
		return Debug
	}
}
