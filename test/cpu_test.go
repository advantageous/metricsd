package test

import (
	g "github.com/cloudurable/metricsd/gatherer"
	c "github.com/cloudurable/metricsd/common"
	"testing"
	"os"
)

func TestCpuCounts(test *testing.T) {

	logger := GetTestLogger(test, "cpu")
	dir, _ := os.Getwd()
	config := c.Config{ Debug: false, CpuProcStat: dir + "/test-data/proc/stat1", }

	cpu := g.NewCPUMetricsGatherer(nil, &config)
	StandardTest(test, logger, cpu)

	cpu.TestingChangeProcStatPath(dir + "/test-data/proc/stat2")
	StandardTest(test, logger, cpu)
}
