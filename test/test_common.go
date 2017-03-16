package test

import (
	lg "github.com/advantageous/go-logback/logging"
	l "github.com/advantageous/go-logback/logging/test"
	c "github.com/cloudurable/metricsd/common"
	"os"
	"flag"
	"testing"
	"fmt"
)

func GetTestLogger(test *testing.T, label string) (lg.Logger) {
	return l.NewTestSimpleLogger(label + "-test", test)
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

func StandardTest(test *testing.T, gatherer c.MetricsGatherer) {
	metrics, err := gatherer.GetMetrics()
	if err == nil {
		AssertMetrics(test, metrics)
		ShowTestMetrics(metrics)
	} else {
		test.Errorf("Error found %s %v", err.Error(), err)
	}
}

func AssertMetrics(test *testing.T, metrics []c.Metric) {
	if metrics == nil {
		test.Error("Nil metrics")
	}

	if len(metrics) == 0 {
		test.Error("Empty metrics")
	}
}

func ShowTestMetrics(metrics []c.Metric) {
	if metrics != nil && len(metrics) > 0 {
		for _,m := range metrics {
			fmt.Println(c.ObjectToString(&m))
		}
	}
}