package metric

import (
	"testing"
	"os"
	"fmt"
	l "github.com/advantageous/metricsd/logger/test"
)
func TestCpuCounts(z *testing.T) {

	test := l.NewTestSimpleLogger("cpu", z)

	dir, _ := os.Getwd()
	fmt.Println("DIR", dir)
	cpuG := NewCPUMetricsGathererWithPath(dir +"/test-data/proc/stat", MetricInterval{30, SECONDS}, test )

	metrics,err:=cpuG.GetMetrics()

	cpuG.path = dir +"/test-data/proc/stat2"
	metrics,err=cpuG.GetMetrics()

	if err!=nil {
		test.Errorf("Error found %s %v", err.Error(), err)
	}

	if len(metrics) == 0 {
		test.Error("Empty metrics")
	}

	metric := metrics[0]

	if metric.GetName() != "softirq" {
		test.Error("softirq not found")
	}


	if metric.GetValue() != 100 {
		test.Errorf("softirq wrong value %d", metric.GetValue())
	}

}