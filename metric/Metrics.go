package metric

type MetricType byte
type MetricIntervalType byte

type MetricValue int64
type MetricIntervalValue int64

const (
	COUNT MetricType = iota
	LEVEL
	TIMING
)

type Metric interface {
	GetProvider() string
	GetType() MetricType
	GetValue() MetricValue
	GetName() string
}

type MetricsGatherer interface {
	GetMetrics() ([]Metric, error)
}

type MetricsRepeater interface {
	ProcessMetrics(metrics []Metric) error
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
