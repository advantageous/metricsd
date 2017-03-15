package nodetool

import (
	"strings"
	c "github.com/cloudurable/metricsd/common"
)

func Netstats(nodetoolCommand string) ([]c.Metric, error) {
	output, err := c.ExecCommand(nodetoolCommand, NodetoolFunction_netstats)
	if err != nil {
		return nil, err
	}

	// -- sample netstats output --
	// [0] // Mode: NORMAL
	// [1] // Not sending any streams.
	// [2] // Read Repair Statistics:
	// [3] // Attempted: 0
	// [4] // Mismatch (Blocking): 0
	// [5] // Mismatch (Background): 0
	// [6] // Pool Name                    Active   Pending      Completed   Dropped
	// [7] // Large messages                  n/a         0              0         0
	// [8] // Small messages                  n/a         0              2         0
	// [9] // Gossip messages                 n/a         0              0         0

	lines := strings.Split(output, c.NEWLINE)

	var metrics = appendNsMode([]c.Metric{}, lines[0])

	metrics = appendNsReadRepair(metrics, lines[3], 1, "ntNsRrAttempted")
	metrics = appendNsReadRepair(metrics, lines[4], 2, "ntNsRrBlocking")
	metrics = appendNsReadRepair(metrics, lines[4], 2, "ntNsRrBackground")

	metrics = appendNsPool(metrics, lines[7], "ntNsPoolLargeMsgs")
	metrics = appendNsPool(metrics, lines[8], "ntNsPoolSmallMsgs")
	metrics = appendNsPool(metrics, lines[9], "ntNsPoolGossipMsgs")

	return metrics, nil
}

func appendNsMode(metrics []c.Metric, line string) []c.Metric {
	value := value_mode_other
	switch strings.ToLower(c.SplitGetFieldByIndex(line, 1)) {
	case "starting":		value = value_mode_starting
	case "normal":			value = value_mode_normal
	case "joining":			value = value_mode_joining
	case "leaving":			value = value_mode_leaving
	case "decommissioned": 	value = value_mode_decommissioned
	case "moving":			value = value_mode_moving
	case "draining":		value = value_mode_draining
	case "drained":			value = value_mode_drained
	}
	return append(metrics, c.Metric{c.MT_NO_UNIT, c.MetricValue(value), c.EMPTY, "ntNsMode", c.PROVIDER_NODETOOL})
}

func appendNsReadRepair(metrics []c.Metric, line string, columnIndex int, name string) []c.Metric {
	metricValue := c.MetricValue( c.ToInt64(c.SplitGetFieldByIndex(line, columnIndex), c.VALUE_ERROR) )
	return append(metrics, c.Metric{c.MT_COUNT, metricValue, c.EMPTY, name, c.PROVIDER_NODETOOL})
}

func appendNsPool(metrics []c.Metric, line string, prefix string) []c.Metric {
	valuesOnly := strings.Fields(line)
	metrics = append(metrics, c.Metric{c.MT_COUNT, c.StrToMetricValue(valuesOnly[2]), c.EMPTY, prefix + "Active", c.PROVIDER_NODETOOL})
	metrics = append(metrics, c.Metric{c.MT_COUNT, c.StrToMetricValue(valuesOnly[3]), c.EMPTY, prefix + "Pending", c.PROVIDER_NODETOOL})
	metrics = append(metrics, c.Metric{c.MT_COUNT, c.StrToMetricValue(valuesOnly[4]), c.EMPTY, prefix + "Completed", c.PROVIDER_NODETOOL})
	return append(metrics, c.Metric{c.MT_COUNT, c.StrToMetricValue(valuesOnly[5]), c.EMPTY, prefix + "Dropped", c.PROVIDER_NODETOOL})
}

