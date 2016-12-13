package metric

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	lg "github.com/advantageous/metricsd/logger"

)

func makeTerminateChannel() <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	return ch
}

func RunWorker(gatherers []MetricsGatherer, repeaters []MetricsRepeater, logger lg.Logger, interval time.Duration) {

	if logger == nil {
		logger = lg.GetSimpleLogger("MT_METRIC_WORKER_DEBUG", "worker")
	}

	timer := time.NewTimer(interval)

	terminator := makeTerminateChannel()

	for {

		select {

		case <-terminator:
			logger.Info("Exiting")
			os.Exit(0)
			break

		case <-timer.C:
			metrics := collectMetrics(gatherers, logger)
			processMetrics(metrics, repeaters, logger)
			timer.Reset(interval)

		}
	}
}
func processMetrics(metrics []Metric, repeaters []MetricsRepeater, logger lg.Logger) {
	for _,r := range repeaters {
		if err:=r.ProcessMetrics(metrics); err !=nil {
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