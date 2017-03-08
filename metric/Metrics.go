package metric

type MetricType byte
type MetricIntervalType byte

type MetricValue int64
type MetricIntervalValue int64

const EMPTY = ""
const SPACE = " "

const LINUX_LABEL = "Linux"
const DARWIN_LABEL = "Darwin"
const CONFIG_LABEL = "Config"
const GOOS_DARWIN = "darwin"

const (
	COUNT MetricType = iota
	LEVEL
	TIMING
	LEVEL_PERCENT
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
}

type MetricsRepeater interface {
	ProcessMetrics(context MetricContext, metrics []Metric) error
}

type metric struct {
	metricType MetricType
	value      MetricValue
	name       string
	provider   string
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
