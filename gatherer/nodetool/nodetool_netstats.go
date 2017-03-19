package nodetool

import (
	"strings"
	c "github.com/cloudurable/metricsd/common"
)

func Netstats(nodetoolCommand string) ([]c.Metric, error) {
	output, err := c.ExecCommand(nodetoolCommand, NtFunc_netstats)
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
	codeStr := strings.ToLower(c.SplitGetFieldByIndex(line, 1))
	code := value_mode_other
	switch codeStr {
	case "starting":		code = value_mode_starting
	case "normal":			code = value_mode_normal
	case "joining":			code = value_mode_joining
	case "leaving":			code = value_mode_leaving
	case "decommissioned": 	code = value_mode_decommissioned
	case "moving":			code = value_mode_moving
	case "draining":		code = value_mode_draining
	case "drained":			code = value_mode_drained
	}
	return append(metrics, *c.NewMetricStringCode(c.MT_NONE, codeStr, code, "ntNsMode", c.PROVIDER_NODETOOL))
}

func appendNsReadRepair(metrics []c.Metric, line string, columnIndex int, name string) []c.Metric {
	return append(metrics, *c.NewMetricIntString(c.MT_COUNT, c.SplitGetFieldByIndex(line, columnIndex), name, c.PROVIDER_NODETOOL))
}

func appendNsPool(metrics []c.Metric, line string, prefix string) []c.Metric {
	valuesOnly := strings.Fields(line)
	metrics = append(metrics, *c.NewMetricIntString(c.MT_COUNT, valuesOnly[2], prefix + "Active", c.PROVIDER_NODETOOL))
	metrics = append(metrics, *c.NewMetricIntString(c.MT_COUNT, valuesOnly[3], prefix + "Pending", c.PROVIDER_NODETOOL))
	metrics = append(metrics, *c.NewMetricIntString(c.MT_COUNT, valuesOnly[4], prefix + "Completed", c.PROVIDER_NODETOOL))
	return append(metrics, *c.NewMetricIntString(c.MT_COUNT, valuesOnly[5], prefix + "Dropped", c.PROVIDER_NODETOOL))
}

