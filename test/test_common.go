package test

import (
	lg "github.com/advantageous/go-logback/logging"
	l "github.com/advantageous/go-logback/logging/test"
	c "github.com/cloudurable/metricsd/common"
	"os"
	"flag"
	"testing"
)

func GetTestLogger(test *testing.T) (lg.Logger) {
	return l.NewTestSimpleLogger("test", test)
}

func GetTestConfig(logger lg.Logger) (*c.Config) {

	// load the config file
	configFile := flag.String("config", "/etc/metricsd.conf", "metrics config")

	config, err := c.LoadConfig(*configFile, logger)
	if err != nil {
		logger.CriticalError("Error reading config", err)
		os.Exit(1)
	}

	return config
}

func ShowTestMetrics(logger lg.Logger, metrics []c.Metric) {
	if metrics != nil && len(metrics) > 0 {
		for _,m := range metrics {
			logger.Info(c.MetricJsonString(&m))
		}
	}
}

