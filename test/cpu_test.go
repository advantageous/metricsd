package test

import (
	"fmt"
	l "github.com/advantageous/go-logback/logging/test"
	m "github.com/cloudurable/metricsd/metric"
	"os"
	"testing"
)

func TestCpuCounts(z *testing.T) {

	test := l.NewTestSimpleLogger("cpu", z)

	dir, _ := os.Getwd()
	fmt.Println("DIR", dir)

	config := m.Config{ Debug: true, CpuProcStat: dir + "/test-data/proc/stat", };
	cpuG := m.NewCPUMetricsGatherer(nil, &config)
	metrics, err := cpuG.GetMetrics()

	config = m.Config{ Debug: true, CpuProcStat: dir + "/test-data/proc/stat2", };
	cpuG = m.NewCPUMetricsGatherer(nil, &config)
	metrics, err = cpuG.GetMetrics()

	if err != nil {
		test.Errorf("Error found %s %v", err.Error(), err)
	}

	if len(metrics) == 0 {
		test.Error("Empty metrics")
	}

	//metric := metrics[0]
	//
	//if metric.GetName() != "softirq" {
	//	test.Error("softirq not found")
	//}
	//
	//if metric.GetValue() != 100 {
	//	test.Errorf("softirq wrong value %d", metric.GetValue())
	//}

}
