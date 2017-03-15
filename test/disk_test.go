package test

import (
	g "github.com/cloudurable/metricsd/gatherer"
	c "github.com/cloudurable/metricsd/common"
	"testing"
)

func TestDisk(test *testing.T) {

	logger := GetTestLogger(test, "disk")

	config := c.Config{
		Debug: false,
		DiskCommand: "df",
		DiskFileSystems: []string{"/dev/*", "udev"},
		DiskFields: []string{"total", "used", "available", "usedpct", "availablepct", "capacitypct", "mount"},
	}

	StandardTest(test, logger, g.NewDiskMetricsGatherer(nil, &config))
}
