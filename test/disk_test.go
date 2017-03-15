package test

import (
	g "github.com/cloudurable/metricsd/gatherer"
	c "github.com/cloudurable/metricsd/common"
	"testing"
)

func TestDisk(test *testing.T) {

	logger := GetTestLogger(test)

	config := c.Config{
		Debug: true,
		DiskCommand: "df",
		DiskFileSystems: []string{"/dev/*", "udev"},
		DiskFields: []string{"total", "used", "available", "usedpct", "availablepct", "capacitypct", "mount"},
	}

	gatherer := g.NewDiskMetricsGatherer(nil, &config)
	metrics, err := gatherer.GetMetrics()
	if err == nil {
		ShowTestMetrics(logger, metrics)
	}

}
