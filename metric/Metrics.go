package metric

type MetricType byte
type MetricIntervalType byte

type MetricValue int64
type MetricIntervalValue int64

type ReloadResult byte

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
)

const (
	DEFAULT_LABEL = "Default"
	CONFIG_LABEL  = "Config"
)

const (
	PROVIDER_CPU      = "cpu"
	PROVIDER_DISK     = "disk"
	PROVIDER_RAM      = "ram"
	PROVIDER_FREE     = "free"
	PROVIDER_NODETOOL = "nodetool"
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

const (
	RELOAD_SUCCESS ReloadResult = iota
	RELOAD_FAILURE
	RELOAD_EJECT
)

type MetricContext interface {
	GetEnv() string
	GetNameSpace() string
	GetRole() string
	SendId() bool
}

type MetricsGatherer interface {
	GetMetrics() ([]Metric, error)
	Reload(config *Config) (ReloadResult)
}

type MetricsRepeater interface {
	ProcessMetrics(context MetricContext, metrics []Metric) error
}

type Metric struct {
	MetricType MetricType
	Value      MetricValue
	Name       string
	Provider   string
}
