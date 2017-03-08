package metric

import (
	l "github.com/advantageous/go-logback/logging"
	"os/exec"
	"strconv"
	"strings"
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

func execCommand(name string, arg ...string) (string, error) {
	if out, err := exec.Command(name, arg...).Output(); err != nil {
		return EMPTY, err
	} else {
		return string(out), nil
	}
}

func toInt64(i string, dflt int64) int64 {
	i64, err := strconv.ParseInt(i, 10, 0)
	if err != nil {
		return dflt
	}
	return i64
}

func splitValuesOnly(text string) []string {
	var results = []string{}
	split := strings.Split(text, SPACE)
	for _,d := range split {
		if (d != EMPTY) { // because the split puts an EMPTY in the array for every extra space
			results = append(results, d)
		}
	}

	return results;
}

func parseForColumn(text string, columnIndex int) string {
	temp := splitValuesOnly(text)
	if len(temp) > columnIndex {
		return temp[columnIndex]
	}
	return EMPTY
}