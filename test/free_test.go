package test

import (
	g "github.com/cloudurable/metricsd/gatherer"
	c "github.com/cloudurable/metricsd/common"
	"testing"
)

func TestFree(test *testing.T) {
	logger := GetTestLogger(test, "free")
	config := c.Config{ Debug: false, FreeCommand: "free"}
	StandardTest(test, logger, g.NewFreeMetricGatherer(nil, &config))
}
