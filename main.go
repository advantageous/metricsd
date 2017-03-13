package main

import (
	l "github.com/advantageous/go-logback/logging"
	c "github.com/cloudurable/metricsd/common"
	r "github.com/cloudurable/metricsd/repeater"
	g "github.com/cloudurable/metricsd/gatherer"
	"os"
	"os/signal"
	"syscall"
	"time"
	"flag"
)

func main() {

	// load the config file
	configFile := flag.String("config", "/etc/metricsd.conf", "metrics config")

	logger := l.NewSimpleLogger("main-init")
	config, err := c.LoadConfig(*configFile, logger)
	if err != nil {
		panic(err)
	}

	logger = c.GetLogger(config.Debug, "main", "MT_MAIN_DEBUG")
	logger.Info("Config file INIT", c.ConfigJsonString(config))

	// begin the work
	interval, intervalConfigRefresh, debug := readRunConfig(config);

	timer := time.NewTimer(interval)
	configTimer := time.NewTimer(intervalConfigRefresh)

	terminator := makeTerminateChannel()

	var gatherers []c.MetricsGatherer
	var repeaters []c.MetricsRepeater
	var load bool = true

	for {
		select {
		case <-terminator:
			logger.Info("Exiting")
			os.Exit(0)

		case <-timer.C:
			if load {
				gatherers = g.LoadGatherers(config)
				repeaters = r.LoadRepeaters(config)
			}
			metrics := collectMetrics(gatherers, logger)
			processMetrics(metrics, repeaters, config, logger)
			timer.Reset(interval)

		case <-configTimer.C:
			if newConfig, err := c.LoadConfig(*configFile, logger); err != nil {
				logger.Error("Error reading config", err)
			} else {
				load = !c.ConfigEquals(config, newConfig)
				if load {
					config = newConfig
					interval, intervalConfigRefresh, debug = readRunConfig(config);
					logger.Info("Config file CHANGED", c.ConfigJsonString(config))
				} else {
					if debug {
						logger.Info("Config file SAME")
					}
				}
			}
			configTimer.Reset(intervalConfigRefresh)
		}
	}
}

func makeTerminateChannel() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}

func readRunConfig(config *c.Config) (time.Duration, time.Duration, bool){
	return 	config.TimePeriodSeconds * time.Second,
		config.ReadConfigSeconds * time.Second,
		config.Debug
}

func processMetrics(metrics []c.Metric, repeaters []c.MetricsRepeater, context *c.Config, logger l.Logger) {
	for _, r := range repeaters {
		if err := r.ProcessMetrics(context, metrics); err != nil {
			logger.PrintError("Repeater failed", err)
		}
	}

	noIdContext := context.GetNoIdContext()

	for _, r := range repeaters {
		if err := r.ProcessMetrics(noIdContext, metrics); err != nil {
			logger.PrintError("Repeater failed", err)
		}
	}
}

func collectMetrics(gatherers []c.MetricsGatherer, logger l.Logger) []c.Metric {

	metrics := []c.Metric{}

	for _, g := range gatherers {
		m, err := g.GetMetrics()
		if err != nil {
			logger.PrintError("Problem getting metrics from gatherer", err)
		}
		metrics = append(metrics, m...)
	}

	return metrics
}
