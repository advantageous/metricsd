package main

import (
	"flag"
	l "github.com/advantageous/go-logback/logging"
	m "github.com/cloudurable/metricsd/metric"
	r "github.com/cloudurable/metricsd/repeater"
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

	if (config.Debug) {
		logger.Println("CPU gathering is on?", config.CpuGather)
		logger.Println("Disk space gathering is on?", config.DiskGather)
		logger.Println("Free memory gathering is on?", config.FreeGather)
	}

	count := 0
	if (config.CpuGather) {  count++ }
	if (config.DiskGather) { count++ }
	if (config.FreeGather) { count++ }

	gatherers := make([]m.MetricsGatherer, count)

	index := 0;
	if (config.CpuGather) {
		gatherers[index] = m.NewCPUMetricsGatherer(nil, config)
		index++
	}

	if (config.DiskGather) {
		gatherers[index] = m.NewDiskMetricsGatherer(nil, config)
		index++
	}

	if (config.FreeGather) {
		gatherers[index] = m.NewFreeMetricGatherer(nil, config)
		index++
	}

	m.RunWorker(gatherers, repeaters, nil, config.TimePeriodSeconds * time.Second,
		config.ReadConfigSeconds * time.Second, config.Debug, *configFile)
}
