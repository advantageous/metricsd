package metric

import (
	//"fmt"
	l "github.com/advantageous/go-logback/logging"
	//"os/exec"
	//"runtime"
	//"strings"
)

func ensureLogger(logger l.Logger, debug bool, name string, flag string) l.Logger {
	if logger == nil {
		if debug {
			logger = l.NewSimpleDebugLogger(name)
		} else {
			logger = l.GetSimpleLogger(flag, name)
		}
	}
	return logger
}
