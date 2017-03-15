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
	RepeatForContext() bool
	RepeatForNoIdContext() bool
}

type Metric struct {
	MetricType MetricType
	Value      MetricValue
	StrValue   string
	Name       string
	Provider   string
}

func MetricJsonString(m *Metric) (string) {
	return "{" +
		Jint64("MetricType", int64(m.MetricType), false) +
		Jint64("Value", int64(m.Value), false) +
		Jstr("StrValue", m.StrValue, false) +
		Jstr("Name", m.Name, false) +
		Jstr("Provider", m.Provider, true) +
		"}"
}
