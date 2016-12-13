package repeater

import (
	lg "github.com/advantageous/metrics/logger"
	m "github.com/advantageous/metrics/metric"

)

type LogMetricsRepeater struct {
	logger lg.Logger
}

func (lr LogMetricsRepeater)ProcessMetrics(metrics []m.Metric) error {
	for _, m := range metrics {
		lr.logger.Printf("%s %d %d", m.GetName(), m.GetType(), m.GetValue())
	}
	return nil
}

func NewLogMetricsRepeater() LogMetricsRepeater {
	logger := lg.NewSimpleLogger("log-repeater")
	return LogMetricsRepeater{logger}
}