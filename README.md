# sglog
The name `sglog` is an abbreviation of the phrase "Simple Go Logger". The goal of the project is to provide simple
logger that allows:
1. Per-package logging.
2. Logging under `Debug`, `Info` and `Warning` levels.
3. Switching logging backends.
4. Configuration via simple config file.

## Basic Usage
Each package should have its own `Logger`. The easiest way to do that is to have a file `log.go` inside each package.
The content of the file should be something like:
```
package greatpackage

import "github.com/junak-michal/sglog"

var log = sglog.GetLogger("github.com/project/greatpackage")
```

The important part is that each package really creates its own `Logger` by passing its name
to `sglog.GetLogger` function. That helps distinguishing log messages once they are written into a
file/syslog/whateverbackend.

Once the `log` variable is defined, functions `Debug(format string, a ...interface{})`, 
`Info(format string, a ...interface{})`, and `Warning(format string, a ...interface{})` can be used to log actual log
messages under corresponding levels.

## Log Message Content
The format of log message is defined by the backend that is used (the default backend logs messages into stderr). The content
of the message should be however similar for all backends. Each log message contains:
- log level of the message,
- package name (from `Logger` - see above),
- file name,
- line number,
- the actual log message.

## Setting Log Level
The default log level is `Debug`, that means that every `Logger` log all levels of messages - `Debug`, `Info`, and `Warning`.

That might not be desirable - especially in the production environment. There are two ways of changing log level.
1. Calling `Logger.SetLevel`. That sets log level for single `Logger` instance.
2. Using YAML config file to set not only default log level but also log level for any package
that should have a different log level than the default one.

### YAML config
Here is an example of the YAML config with comments that should explain its structure.
```
# Simple example that sets the default log level to Info and sets different log level for two packages.
# Version must be 1.
version: 1
# Every logger that is not listed in the loggers section gets the default log level.
# If omitted, Debug is the default.
default: Info
# Map of loggers (identified by a package name) and their log levels.
loggers:
  github.com/project/tooverbose: Warning
  github.com/project/problematic: Debug
```

The YAML config is loaded when the sglog package is loaded. The default path
of the YAML config is `sglog.yaml`. It might be desirable to change that default path,
which is very simple with Go. You can use the "-X" flag to change value of variable `github.com/junak-michal/sglog.ConfigFile`.
For example if the desired YAML config path is `/etc/myservice/logging.yaml`, build your program following way:
```
go build -ldflags "-X github.com/junak-michal/sglog.ConfigFile=/etc/myservice/logging.yaml"
```
