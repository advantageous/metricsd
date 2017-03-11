package run

import (
	l "github.com/advantageous/go-logback/logging"
	c "github.com/cloudurable/metricsd/common"
	"os"
	"os/signal"
	"syscall"
	"time"
	"flag"
)

func makeTerminateChannel() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}


func Main() {

	configFile := flag.String("config", "/etc/metricsd.conf", "metrics config")
	logger := l.NewSimpleLogger("main")

	config, err := c.LoadConfig(*configFile, logger)
	if err != nil {
		panic(err)
	}

	RunWorker(nil, config, *configFile)
}

func RunWorker(logger l.Logger, config *c.Config,  configFile string) {
	logger = c.EnsureLogger(logger, config.Debug, "worker", "MT_METRIC_WORKER_DEBUG")

	logger.Info("Config file INIT", c.ConfigJsonString(config))
	interval, intervalConfigRefresh, debug := readRunConfig(config);

	timer := time.NewTimer(interval)
	configTimer := time.NewTimer(10) // intervalConfigRefresh)

	terminator := makeTerminateChannel()

	var gatherers = LoadGatherers(config)
	var repeaters = LoadRepeaters(config)
	var configChanged bool = false

	for {
		select {
		case <-terminator:
			logger.Info("Exiting")
			os.Exit(0)
			break // ask rick, I think this is redundant in go

		case <-timer.C:
			if configChanged {
				gatherers = LoadGatherers(config)
				repeaters = LoadRepeaters(config)
			}

			metrics := collectMetrics(gatherers, logger)
			processMetrics(metrics, repeaters, config, logger)
			timer.Reset(interval)

		case <-configTimer.C:
			if newConfig, err := c.LoadConfig(configFile, logger); err != nil {
				logger.Error("Error reading config", err)
			} else {
				configChanged = !c.ConfigEquals(config, newConfig)
				if configChanged {
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
