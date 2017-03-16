package repeater

import (
	c "github.com/cloudurable/metricsd/common"
	"fmt"
)

type ConsoleMetricsRepeater struct {}


func (lr ConsoleMetricsRepeater) ProcessMetrics(context c.MetricContext, metrics []c.Metric) error {
	for _, m := range metrics {
		fmt.Println(c.ObjectToString(&m))
	}
	return nil
}

func (lr ConsoleMetricsRepeater) RepeatForContext() bool { return false; }
func (lr ConsoleMetricsRepeater) RepeatForNoIdContext() bool { return true; }

func NewConsoleMetricsRepeater() *ConsoleMetricsRepeater {
	return &ConsoleMetricsRepeater{}
}
