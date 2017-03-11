package common

type MetricType byte
type MetricIntervalType byte

type MetricValue int64
type MetricIntervalValue int64

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

type Metric struct {
	MetricType MetricType
	Value      MetricValue
	Name       string
	Provider   string
}
