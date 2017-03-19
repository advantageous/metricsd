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
		lr.logger.Println(c.ObjectToString(&m))
	}
	return nil
}

func (lr LogMetricsRepeater) RepeatForContext() bool { return false; }
func (lr LogMetricsRepeater) RepeatForNoIdContext() bool { return true; }

func NewLogMetricsRepeater() *LogMetricsRepeater {
	return &LogMetricsRepeater{ lg.NewSimpleLogger("log-repeater")	}
}
