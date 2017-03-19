package common

const (
	EMPTY                   = ""
	SPACE                   = " "
	NEWLINE                 = "\n"
	UNDER                   = "_"
	DOT                     = "."
	QUOTE                   = "\""
	COMMA                   = ","
	COMMA_SPACE             = ", "
	QUOTE_COLON_SPACE       = "\": "
	QUOTE_COLON_SPACE_QUOTE = "\": \""
	QUOTE_COMMA_SPACE       = "\", "
	OPEN_BRACE              = "["
	CLOSE_BRACE             = "]"
	COLON                   = ":"
)

const (
	DEFAULT_LABEL = "Default"
	CONFIG_LABEL  = "Config"
)

const (
	PROVIDER_CPU      = "cpu"
	PROVIDER_DISK     = "disk"
	PROVIDER_FREE     = "free"
	PROVIDER_NODETOOL = "nodetool"

	PROVIDER_RAM      = "ram"
)

const (
	REPEATER_AWS    = "aws"
	REPEATER_LOGGER = "logger"
	REPEATER_CONSOLE = "console"
)

const (
	FLAG_CPU      = "MT_CPU_DEBUG"
	FLAG_DISK     = "MT_DISK_DEBUG"
	FLAG_FREE     = "MT_FREE_DEBUG"
	FLAG_NODETOOL = "MT_NODETOOL_DEBUG"
)

type MetricValueSource int8
const (
	MVS_INT MetricValueSource = iota
	MVS_FLOAT
	MVS_STRING
)

func (mvs *MetricValueSource) Name() string {
	switch *mvs {
	case MVS_INT: return "int"
	case MVS_FLOAT: return "float"
	case MVS_STRING: return "str"
	}
	return EMPTY
}

type MetricType int8
const (
	MT_COUNT MetricType = iota
	MT_PERCENT
	MT_MICROS
	MT_MILLIS
	MT_SIZE_BYTE
	MT_SIZE_KB
	MT_SIZE_MB
	MT_SIZE_GB
	MT_SIZE_TB
	MT_NONE
)

func (mt *MetricType) Name() string {
	switch *mt {
	case MT_COUNT:     return "Count"
	case MT_PERCENT:   return "Percent"
	case MT_MICROS:    return "Microseconds"
	case MT_MILLIS:    return "Milliseconds"
	case MT_SIZE_BYTE: return "Byte"
	case MT_SIZE_MB:   return "Megabytes"
	case MT_SIZE_KB:   return "Kilobytes"
	case MT_NONE:   return "None"
	}
	return EMPTY
}

/* Cloudwatch for reference
const (
	StandardUnitSeconds = "Seconds"
	StandardUnitMicroseconds = "Microseconds"
	StandardUnitMilliseconds = "Milliseconds"
	StandardUnitBytes = "Bytes"
	StandardUnitKilobytes = "Kilobytes"
	StandardUnitMegabytes = "Megabytes"
	StandardUnitGigabytes = "Gigabytes"
	StandardUnitTerabytes = "Terabytes"
	StandardUnitBits = "Bits"
	StandardUnitKilobits = "Kilobits"
	StandardUnitMegabits = "Megabits"
	StandardUnitGigabits = "Gigabits"
	StandardUnitTerabits = "Terabits"
	StandardUnitPercent = "Percent"
	StandardUnitCount = "Count"
	StandardUnitBytesSecond = "Bytes/Second"
	StandardUnitKilobytesSecond = "Kilobytes/Second"
	StandardUnitMegabytesSecond = "Megabytes/Second"
	StandardUnitGigabytesSecond = "Gigabytes/Second"
	StandardUnitTerabytesSecond = "Terabytes/Second"
	StandardUnitBitsSecond = "Bits/Second"
	StandardUnitKilobitsSecond = "Kilobits/Second"
	StandardUnitMegabitsSecond = "Megabits/Second"
	StandardUnitGigabitsSecond = "Gigabits/Second"
	StandardUnitTerabitsSecond = "Terabits/Second"
	StandardUnitCountSecond = "Count/Second"
	StandardUnitNone = "None"
)
*/

const (
	VALUE_N_A   int64 = -125
	VALUE_NAN   int64 = -126
	VALUE_ERROR int64 = -127
)

const (
	IN_VALUE_N_A = "n/a"
	IN_VALUE_NAN = "NaN"
)
