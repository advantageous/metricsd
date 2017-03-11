package repeater

import (
	lg "github.com/advantageous/go-logback/logging"
	c "github.com/cloudurable/metricsd/common"
)

type LogMetricsRepeater struct {
	logger lg.Logger
}

func (lr LogMetricsRepeater) ProcessMetrics(context c.MetricContext, metrics []c.Metric) error {
	for _, m := range metrics {
		lr.logger.Printf("%s %d %d", m.Name, m.MetricType, m.Value)
	}
	return nil
}

func NewLogMetricsRepeater() *LogMetricsRepeater {
	logger := lg.NewSimpleLogger("log-repeater")
	return &LogMetricsRepeater{logger}
}
