package main

import (
	"flag"
	l "github.com/advantageous/metricsd/logger"
	m "github.com/advantageous/metricsd/metric"
	r "github.com/advantageous/metricsd/repeater"
	"time"
)

func main() {

	configFile := flag.String("host", "/etc/metricsd.conf", "metrics config")
	logger := l.NewSimpleLogger("main")

	gatherers := []m.MetricsGatherer{m.NewCPUMetricsGatherer(nil),
		m.NewDiskMetricsGatherer(nil),
		m.NewFreeMetricGatherer(nil)}

	config, err := m.LoadConfig(*configFile, logger)
	if err != nil {
		panic(err)
	}

	repeaters := []m.MetricsRepeater{r.NewAwsCloudMetricRepeater(config)}

	m.RunWorker(gatherers, repeaters, nil, config.TimePeriodSeconds * time.Second)
}
