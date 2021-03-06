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

func RunWorker(gatherers []MetricsGatherer, repeaters []MetricsRepeater, logger lg.Logger, interval time.Duration,
	intervalConfigRefresh time.Duration, debug bool, configFile string) {

	if logger == nil {
		if debug {
			logger = lg.NewSimpleDebugLogger("worker")
		} else {
			logger = lg.GetSimpleLogger("MT_METRIC_WORKER_DEBUG", "worker")
		}
	}

	timer := time.NewTimer(interval)

	configTimer := time.NewTimer(intervalConfigRefresh)

	var config *Config

	if newConfig, err := LoadConfig(configFile, logger); err != nil {
		logger.Error("Error reading config", err)
	} else {
		config = newConfig
	}

	terminator := makeTerminateChannel()

	for {

		select {

		case <-terminator:
			logger.Info("Exiting")
			os.Exit(0)
			break

		case <-timer.C:
			metrics := collectMetrics(gatherers, logger)
			processMetrics(metrics, repeaters, config, logger)
			timer.Reset(interval)

		case <-configTimer.C:
			if newConfig, err := LoadConfig(configFile, logger); err != nil {
				logger.Error("Error reading config", err)
			} else {
				config = newConfig
				if debug {
					logger.Info("LOADED NEW CONFIG", "ENV", config.GetEnv(),
						"NAMESPACE", config.GetNameSpace(),
						"ROLE", config.GetRole())
				}
			}
			configTimer.Reset(intervalConfigRefresh)

		}
	}
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
