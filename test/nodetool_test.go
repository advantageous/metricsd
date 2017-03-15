package test

import (
	g "github.com/cloudurable/metricsd/gatherer"
	c "github.com/cloudurable/metricsd/common"
	"testing"
)

func TestNodetool(test *testing.T) {

	logger := GetTestLogger(test)

	config := c.Config{
		Debug: true,
		// NodetoolFunctions: []string{"cfstats", "tpstats", "gcstats", "getlogginglevels", "netstats", "gettimeout"},
		NodetoolFunctions: []string{"cfstats"},
	}

	gatherers := g.NewNodetoolMetricGatherers(nil, &config)
	for _,gatherer := range gatherers {
		metrics, err := gatherer.GetMetrics()
		if err == nil {
			ShowTestMetrics(logger, metrics)
		}
	}
}
