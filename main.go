package main

import (
	m "github.com/advantageous/metricsd/metric"
	r "github.com/advantageous/metricsd/repeater"
	l "github.com/advantageous/metricsd/logger"
	"time"
)

func main() {


	logger := l.NewSimpleLogger("main")

	logger.Println("Starting up")


	gatherers := []m.MetricsGatherer{m.NewCPUMetricsGatherer( nil  ), m.NewDiskMetricsGatherer(nil), m.NewFreeMetricGatherer(nil)}
	repeaters := []m.MetricsRepeater{r.NewLogMetricsRepeater()}

	m.RunWorker(gatherers, repeaters, nil, time.Second * 10)
}
