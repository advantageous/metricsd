package metric

type MetricType byte
type MetricIntervalType byte

type MetricValue int64
type MetricIntervalValue int64

const (
	EMPTY = ""
	SPACE = " "
	NEWLINE = "\n"
	UNDER = "_"
)

const (
	DEFAULT_LABEL = "Default"
	CONFIG_LABEL  = "Config"
)

const (
	PROVIDER_CPU  = "cpu"
	PROVIDER_DISK = "disk"
	PROVIDER_RAM  = "ram"
	PROVIDER_FREE = "free"
	PROVIDER_NODE = "node"
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

type Metric interface {
	GetProvider() string
	GetType() MetricType
	GetValue() MetricValue
	GetName() string
}

type MetricContext interface {
	GetEnv() string
	GetNameSpace() string
	GetRole() string
	SendId() bool
}

type MetricsGatherer interface {
	GetMetrics() ([]Metric, error)
	//GetIdentity(string)
}

type MetricsRepeater interface {
	ProcessMetrics(context MetricContext, metrics []Metric) error
}

type metric struct {
	metricType	MetricType
	value		MetricValue
	name		string
	provider	string
}

func (m metric) GetType() MetricType {
	return m.metricType
}

func (m metric) GetValue() MetricValue {
	return m.value
}

func (m metric) GetProvider() string {
	return m.provider
}

func (m metric) GetName() string {
	return m.name
}
