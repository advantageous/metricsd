package common

import (
	l "github.com/advantageous/go-logback/logging"
	"os/exec"
	"strconv"
	"strings"
	"time"
	"math"
	"fmt"
)

// ========================================================================================================================
// LOGGING HELPERS
// ========================================================================================================================
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

// ========================================================================================================================
// EXEC HELPERS
// ========================================================================================================================
func ExecCommand(name string, arg ...string) (string, error) {
	if out, err := exec.Command(name, arg...).Output(); err != nil {
		return EMPTY, err
	} else {
		return string(out), nil
	}
}

// ========================================================================================================================
// DEBUG HELPERS
// ========================================================================================================================
func Dump(logger l.Logger, arr []string, label string) {
	for _,s := range arr {
		logger.Debug(label + " -->" + s + "<--")
	}

}

// ========================================================================================================================
// STRING TO NUMBER CONVERSIONS
// ========================================================================================================================
func ToInt64(i string, dflt int64) int64 {
	i64, err := strconv.ParseInt(i, 10, 0)
	if err != nil {
		return dflt
	}
	return i64
}

func ToFloat64(f string, dflt float64) float64 {
	f64, err := strconv.ParseFloat(f, 64)
	if err != nil {
		return dflt
	}
	return f64
}

// ========================================================================================================================
// OBJECT TO STRING CONVERSIONS
// ========================================================================================================================
func BoolToString(b bool) string {
	return strconv.FormatBool(b)
}

func DurationToString(dur time.Duration) string {
	return strconv.FormatInt(int64(dur), 10)
}

func ByteToString(b byte) string {
	return strconv.FormatInt(int64(b), 10)
}

func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func Float64ToStringPrecise(f float64, prec int) string {
	return strconv.FormatFloat(f, 'f', prec, 64)
}

// ========================================================================================================================
// JSON STUFF
// ========================================================================================================================
func Jstr(name string, v string, last bool) string {
	if last {
		return QUOTE + name + QUOTE_COLON_SPACE_QUOTE + v + QUOTE
	}
	return QUOTE + name + QUOTE_COLON_SPACE_QUOTE + v + QUOTE_COMMA_SPACE
}

func Jbyte(name string, v byte, last bool) string {
	return jnum(name, ByteToString(v), last)
}

func Jdur(name string, v time.Duration, last bool) string {
	return jnum(name, DurationToString(v), last)
}

func Jint64(name string, v int64, last bool) string {
	return jnum(name, Int64ToString(v), last)
}

func Jfloat64(name string, v float64, last bool) string {
	return jnum(name, Float64ToString(v), last)
}

func Jfloat64Precise(name string, v float64, prec int, last bool) string {
	return jnum(name, Float64ToStringPrecise(v, prec), last)
}

func jnum(name string, numStr string, last bool) string {
	if last {
		return QUOTE + name + QUOTE_COLON_SPACE + numStr
	}
	return QUOTE + name + QUOTE_COLON_SPACE + numStr + COMMA_SPACE
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

// ========================================================================================================================
// MATH
// ========================================================================================================================
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

// ========================================================================================================================
// STRING READING
// ========================================================================================================================
func GetLastIndex(a []string) int {
	if a != nil && len(a) > 0 {
		return len(a) - 1
	}
	return -1
}

func GetLastField(a []string) string {
	if a != nil && len(a) > 0 {
		return a[len(a) - 1]
	}
	return EMPTY
}

func GetFieldByIndex(a []string, columnIndex int) string {
	if a != nil && len(a) > columnIndex {
		return a[columnIndex]
	}
	return EMPTY
}

func SplitGetFieldByIndex(text string, columnIndex int) string {
	return GetFieldByIndex(strings.Fields(text), columnIndex)
}

func SplitGetLastField(text string) string {
	return GetLastField(strings.Fields(text))
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

// ========================================================================================================================
// STRING MANIPULATION
// ========================================================================================================================
func UpFirst(s string) string {
	return strings.ToUpper(s[0:1]) + s[1:]
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

func ObjectToString(object interface{}) string {
	s := fmt.Sprintf("%#v", object)
	open := strings.Index(s, "{")
	dot := strings.Index(s, DOT)
	if dot != -1 && dot < open {
		s = s[dot+1:]
	}
	s = strings.Replace(s, "[]string", "[]", -1)
	s = strings.Replace(s, "(nil)", EMPTY, -1)
	return s
}

func ToSizeMetricType(size string) MetricType {
	switch strings.ToUpper(size) {
	case "BYTE":  return MT_SIZE_BYTE
	case "BYTES": return MT_SIZE_BYTE
	case "KB":    return MT_SIZE_KB
	case "MB":    return MT_SIZE_MB
	case "GB":    return MT_SIZE_GB
	case "TB":    return MT_SIZE_TB
	}
	return MT_NONE
}