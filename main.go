package main

import (
	"flag"
	l "github.com/advantageous/go-logback/logging"
	m "github.com/cloudurable/metricsd/metric"
	r "github.com/cloudurable/metricsd/repeater"
	"time"
	"strings"
)

func main() {

	configFile := flag.String("config", "/etc/metricsd.conf", "metrics config")
	logger := l.NewSimpleLogger("main")

	config, err := m.LoadConfig(*configFile, logger)
	if err != nil {
		panic(err)
	}

	repeaters := []m.MetricsRepeater{r.NewAwsCloudMetricRepeater(config)}

	var gatherers = []m.MetricsGatherer{}
	if (config.CpuGather) {
		gatherers = append(gatherers, m.NewCPUMetricsGatherer(nil, config))
	}

	if (config.DiskGather) {
		gatherers = append(gatherers, m.NewDiskMetricsGatherer(nil, config))
	}

	if (config.FreeGather) {
		gatherers = append(gatherers, m.NewFreeMetricGatherer(nil, config))
	}

	if (config.NodetoolGather) {
		nodetoolFunctions := strings.Split(config.NodetoolFunctions, m.SPACE)
		for _,nodeFunction := range nodetoolFunctions {
			if m.NodetoolFunctionSupported(nodeFunction) {
				gatherers = append(gatherers, m.NewNodetoolMetricGatherer(nil, config, nodeFunction))
			}

		}
	}

	m.RunWorker(gatherers, repeaters, nil, config.TimePeriodSeconds * time.Second,
		config.ReadConfigSeconds * time.Second, config.Debug, *configFile)
}
