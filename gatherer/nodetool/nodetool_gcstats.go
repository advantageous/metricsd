package nodetool

import (
	"strings"
	c "github.com/cloudurable/metricsd/common"
)

func Gcstats(nodetoolCommand string) ([]c.Metric, error) {
	output, err := c.ExecCommand(nodetoolCommand, NodetoolFunction_gcstats)
	if err != nil {
		return nil, err
	}

	// -- sample gcstats output --
	//Interval (ms) Max GC Elapsed (ms)Total GC Elapsed (ms)Stdev GC Elapsed (ms)   GC Reclaimed (MB)         Collections      Direct Memory Bytes
	//3491665                   0                   0                 NaN                   0                   0                       -1

	lines := strings.Split(output, c.NEWLINE)
	values := strings.Fields(lines[1])

	var metrics = []c.Metric{}
	metrics = append(metrics, c.Metric{c.MT_MILLIS, c.StrToMetricValue(values[0]), c.EMPTY, "ntGcInterval", c.PROVIDER_NODETOOL})
	metrics = append(metrics, c.Metric{c.MT_MILLIS, c.StrToMetricValue(values[1]), c.EMPTY, "ntGcMaxElapsed", c.PROVIDER_NODETOOL})
	metrics = append(metrics, c.Metric{c.MT_MILLIS, c.StrToMetricValue(values[2]), c.EMPTY, "ntGcTotalElapsed", c.PROVIDER_NODETOOL})
	metrics = append(metrics, c.Metric{c.MT_MILLIS, c.StrToMetricValue(values[3]), c.EMPTY, "ntGcStdevElapsed", c.PROVIDER_NODETOOL})
	metrics = append(metrics, c.Metric{c.MT_SIZE_MB, c.StrToMetricValue(values[4]), c.EMPTY, "ntGcReclaimed", c.PROVIDER_NODETOOL})
	metrics = append(metrics, c.Metric{c.MT_COUNT, c.StrToMetricValue(values[5]), c.EMPTY, "ntGcCollections", c.PROVIDER_NODETOOL})
	metrics = append(metrics, c.Metric{c.MT_SIZE_B, c.StrToMetricValue(values[6]), c.EMPTY, "ntGcDirectMemoryBytes", c.PROVIDER_NODETOOL})

	return metrics, nil
}

