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
			logger = l.NewSimpleDebugLogger("cpu")
		} else {
			logger = l.GetSimpleLogger("MT_CPU_DEBUG", "cpu")
		}
	}
	return logger
}
