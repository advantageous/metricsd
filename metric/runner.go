package metric

import (
	lg "github.com/advantageous/go-logback/logging"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func makeTerminateChannel() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}

func RunWorker(repeaters []MetricsRepeater, logger lg.Logger, config *Config,  configFile string) {
	logger = EnsureLogger(logger, config.Debug, "worker", "MT_METRIC_WORKER_DEBUG")

	logger.Info("Config file INIT", ConfigJsonString(config))
	interval, intervalConfigRefresh, debug := readRunConfig(config);

	timer := time.NewTimer(interval)
	configTimer := time.NewTimer(10) // intervalConfigRefresh)

	terminator := makeTerminateChannel()

	var gatherers = LoadGatherers(config)
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
			}

			metrics := collectMetrics(gatherers, logger)
			processMetrics(metrics, repeaters, config, logger)
			timer.Reset(interval)

		case <-configTimer.C:
			if newConfig, err := LoadConfig(configFile, logger); err != nil {
				logger.Error("Error reading config", err)
			} else {
				configChanged = !ConfigEquals(config, newConfig)
				if configChanged {
					config = newConfig
					interval, intervalConfigRefresh, debug = readRunConfig(config);
					logger.Info("Config file CHANGED", ConfigJsonString(config))
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

func readRunConfig(config * Config) (time.Duration, time.Duration, bool){
	return 	config.TimePeriodSeconds * time.Second,
	 		config.ReadConfigSeconds * time.Second,
			config.Debug
}

func processMetrics(metrics []Metric, repeaters []MetricsRepeater, context *Config, logger lg.Logger) {
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

func collectMetrics(gatherers []MetricsGatherer, logger lg.Logger) []Metric {

	metrics := []Metric{}

	for _, g := range gatherers {
		m, err := g.GetMetrics()
		if err != nil {
			logger.PrintError("Problem getting metrics from gatherer", err)
		}
		metrics = append(metrics, m...)
	}

	return metrics
}
