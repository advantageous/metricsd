package main

import (
	"flag"
	l "github.com/advantageous/go-logback/logging"
	m "github.com/advantageous/metricsd/metric"
	r "github.com/advantageous/metricsd/repeater"
	"time"
)

func main() {

	configFile := flag.String("config", "/etc/metricsd.conf", "metrics config")
	logger := l.NewSimpleLogger("main")

	config, err := m.LoadConfig(*configFile, logger)
	if err != nil {
		panic(err)
	}

	repeaters := []m.MetricsRepeater{r.NewAwsCloudMetricRepeater(config)}

	gatherers := []m.MetricsGatherer{
		m.NewCPUMetricsGatherer(nil, config),
		m.NewDiskMetricsGatherer(nil, config),
		m.NewFreeMetricGatherer(nil, config)}

	m.RunWorker(gatherers, repeaters, nil, config.TimePeriodSeconds*time.Second)
}
