package metric

import (
	l "github.com/advantageous/go-logback/logging"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func EnsureLogger(logger l.Logger, debug bool, name string, flag string) l.Logger {
	if logger == nil {
		if debug {
			logger = l.NewSimpleDebugLogger(name)
		} else {
			logger = l.GetSimpleLogger(flag, name)
		}
	}
	return logger
}

func ExecCommand(name string, arg ...string) (string, error) {
	if out, err := exec.Command(name, arg...).Output(); err != nil {
		return EMPTY, err
	} else {
		return string(out), nil
	}
}

func Dump(logger l.Logger, arr []string, label string) {
	for _,s := range arr {
		logger.Debug(label + " -->" + s + "<--")
	}

}

func ToInt64(i string, dflt int64) int64 {
	i64, err := strconv.ParseInt(i, 10, 0)
	if err != nil {
		return dflt
	}
	return i64
}

func FieldByIndex(text string, columnIndex int) string {
	temp := strings.Fields(text)
	if len(temp) > columnIndex {
		return temp[columnIndex]
	}
	return EMPTY
}

func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

func DurationToString(dur time.Duration) string {
	return strconv.FormatInt(int64(dur), 10)
}

func Jstr(name string, v string, last bool) string {
	if last {
		return QUOTE + name + QUOTE_COLON_QUOTE + v + QUOTE
	}
	return QUOTE + name + QUOTE_COLON_QUOTE + v + QUOTE_COMMA
}

func Jdur(name string, v time.Duration, last bool) string {
	if last {
		return QUOTE + name + QUOTE_COLON + DurationToString(v)
	}
	return QUOTE + name + QUOTE_COLON + DurationToString(v) + COMMA
}

func Jbool(name string, v bool, last bool) string {
	if last {
		return QUOTE + name + QUOTE_COLON + BoolToString(v)
	}
	return QUOTE + name + QUOTE_COLON + BoolToString(v) + COMMA
}
