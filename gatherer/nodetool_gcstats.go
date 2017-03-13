package gatherer

import (
	"strings"
	c "github.com/cloudurable/metricsd/common"
)

func gcstats(nodetoolCommand string) ([]c.Metric, error) {
	output, err := c.ExecCommand(nodetoolCommand, function_gc_stats)
	if err != nil {
		return nil, err
	}

	// -- sample gcstats output --
	//Interval (ms) Max GC Elapsed (ms)Total GC Elapsed (ms)Stdev GC Elapsed (ms)   GC Reclaimed (MB)         Collections      Direct Memory Bytes
	//3491665                   0                   0                 NaN                   0                   0                       -1

	lines := strings.Split(output, c.NEWLINE)
	values := strings.Fields(lines[1])

	var metrics = []c.Metric{}
	metrics = append(metrics, c.Metric{c.TIMING_MS, numericMetricValue(values[0]), "gcInterval", c.PROVIDER_NODETOOL})
	metrics = append(metrics, c.Metric{c.TIMING_MS, numericMetricValue(values[1]), "gcMaxElapsed", c.PROVIDER_NODETOOL})
	metrics = append(metrics, c.Metric{c.TIMING_MS, numericMetricValue(values[2]), "gcTotalElapsed", c.PROVIDER_NODETOOL})
	metrics = append(metrics, c.Metric{c.TIMING_MS, numericMetricValue(values[3]), "gcStdevElapsed", c.PROVIDER_NODETOOL})
	metrics = append(metrics, c.Metric{c.SIZE_MB, numericMetricValue(values[4]), "gcReclaimed", c.PROVIDER_NODETOOL})
	metrics = append(metrics, c.Metric{c.COUNT, numericMetricValue(values[5]), "gcCollections", c.PROVIDER_NODETOOL})
	metrics = append(metrics, c.Metric{c.SIZE_B, numericMetricValue(values[6]), "gcDirectMemoryBytes", c.PROVIDER_NODETOOL})

	return metrics, nil
}

