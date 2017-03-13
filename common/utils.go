package common

import (
	l "github.com/advantageous/go-logback/logging"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"math"
)

func GetLogger(debug bool, name string, flag string) l.Logger {
	return EnsureLogger(nil, debug, name, flag)
}

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

func StringArraysEqual(sa1 []string, sa2 []string) bool {
	saLen := len(sa1)
	if (saLen != len(sa2)) {
		return false
	}

	for i := 0; i < saLen; i++ {
		if (sa1[i] != sa2[i]) {
			return false
		}
	}

	return true
}

func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

func DurationToString(dur time.Duration) string {
	return strconv.FormatInt(int64(dur), 10)
}

func Jstr(name string, v string, last bool) string {
	if last {
		return QUOTE + name + QUOTE_COLON_SPACE_QUOTE + v + QUOTE
	}
	return QUOTE + name + QUOTE_COLON_SPACE_QUOTE + v + QUOTE_COMMA_SPACE
}

func Jdur(name string, v time.Duration, last bool) string {
	if last {
		return QUOTE + name + QUOTE_COLON_SPACE + DurationToString(v)
	}
	return QUOTE + name + QUOTE_COLON_SPACE + DurationToString(v) + COMMA_SPACE
}

func Jbool(name string, v bool, last bool) string {
	if last {
		return QUOTE + name + QUOTE_COLON_SPACE + BoolToString(v)
	}
	return QUOTE + name + QUOTE_COLON_SPACE + BoolToString(v) + COMMA_SPACE
}

func Junquoted(name string, v string, last bool) string {
	if last {
		return QUOTE + name + QUOTE_COLON_SPACE + v
	}
	return QUOTE + name + QUOTE_COLON_SPACE + v + COMMA_SPACE
}

func Jstrarr(name string, v []string, last bool) string {
	if v == nil || len(v) == 0 {
		return Junquoted(name, "[]", last)
	}

	temp := EMPTY
	lastStr := COMMA
	if last {
		lastStr = EMPTY
	}

	lastIndex := len(v) - 1
	for i := 0; i < lastIndex; i++ {
		temp = temp + QUOTE + v[i] + QUOTE_COMMA_SPACE
	}
	temp = temp + QUOTE + v[lastIndex] + QUOTE

	return QUOTE + name + QUOTE_COLON_SPACE + OPEN_BRACE + temp + CLOSE_BRACE + lastStr
}

func ArrayToString(a []string) string {
	result := OPEN_BRACE
	for _, s := range a {
		if (result == OPEN_BRACE) {
			result = result + QUOTE + s + QUOTE
		} else {
			result = result + COMMA + SPACE + QUOTE + s + QUOTE
		}
	}
	return result + CLOSE_BRACE
}

func Round(f float64) int64 {
	t := math.Trunc(f)
	x := math.Trunc( (f - t) * 100 )
	if x < 50 {
		return int64(t)
	}
	return  int64(t) + 1
}

func Percent(top float64, bot float64) float64 {
	return top * 100 / bot
}
