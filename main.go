package main

import (
	"flag"
	l "github.com/advantageous/go-logback/logging"
	m "github.com/cloudurable/metricsd/metric"
	r "github.com/cloudurable/metricsd/repeater"
)

func main() {

	configFile := flag.String("config", "/etc/metricsd.conf", "metrics config")
	logger := l.NewSimpleLogger("main")

	config, err := m.LoadConfig(*configFile, logger)
	if err != nil {
		panic(err)
	}

	var repeaters = []m.MetricsRepeater{r.NewAwsCloudMetricRepeater(config)}
	repeaters = []m.MetricsRepeater{}

	m.RunWorker(repeaters, nil, config, *configFile)
}
