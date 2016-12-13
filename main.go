package main

import (
	m "github.com/advantageous/metrics/metric"
	r "github.com/advantageous/metrics/repeater"
	l "github.com/advantageous/metrics/logger"
	"time"
)

func main() {


	logger := l.NewSimpleLogger("main")

	logger.Println("Starting up")


	gatherers := []m.MetricsGatherer{m.NewCPUMetricsGatherer( nil  ), m.NewDiskMetricsGatherer(nil), m.NewFreeMetricGatherer(nil)}
	repeaters := []m.MetricsRepeater{r.NewLogMetricsRepeater()}

	m.RunWorker(gatherers, repeaters, nil, time.Second * 10)
}
