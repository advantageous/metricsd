package common

const (
	EMPTY             = ""
	SPACE             = " "
	NEWLINE           = "\n"
	UNDER             = "_"
	DOT               = "."
	QUOTE             = "\""
	COMMA             = ","
	QUOTE_COLON       = "\" : "
	QUOTE_COLON_QUOTE = "\" : \""
	QUOTE_COMMA       = "\","
	OPEN_BRACE       = "["
	CLOSE_BRACE       = "]"
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
	COUNT MetricType = iota
	LEVEL
	LEVEL_PERCENT
	TIMING_MS
	SIZE_B
	SIZE_MB
	NO_UNIT
)