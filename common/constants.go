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
)

const (
	FLAG_CPU  = "MT_CPU_DEBUG"
	FLAG_DISK = "MT_DISK_DEBUG"
	FLAG_FREE = "MT_FREE_DEBUG"
	FLAG_NODE = "MT_NODE_DEBUG"
)

const (
	MT_COUNT   MetricType = iota
	MT_PERCENT
	MT_MILLIS
	MT_SIZE_B
	MT_SIZE_MB
	MT_SIZE_K
	MT_NO_UNIT
)

const (
	VALUE_N_A    = -125
	VALUE_NAN   = -126
	VALUE_ERROR = -127
)

const (
	IN_VALUE_N_A = "n/a"
	IN_VALUE_NAN = "NaN"
)
