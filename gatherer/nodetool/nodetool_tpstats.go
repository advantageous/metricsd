package nodetool

import (
	"strings"
	c "github.com/cloudurable/metricsd/common"
)

func Tpstats(nodetoolCommand string) ([]c.Metric, error) {
	output, err := c.ExecCommand(nodetoolCommand, NtFunc_tpstats)
	if err != nil {
		return nil, err
	}

	//Pool Name                         Active   Pending      Completed   Blocked  All time blocked
	//ReadStage                              0         0              3         0                 0
	//MiscStage                              0         0              0         0                 0
	//CompactionExecutor                     0         0          51133         0                 0
	//MutationStage                          0         0              1         0                 0
	//...                                    0         0              1         0                 0
	//
	//Message type           Dropped
	//READ                         0
	//RANGE_SLICE                  0
	//_TRACE                       0
	//HINT                         0
	//MUTATION                     0
	//COUNTER_MUTATION             0
	//BATCH_STORE                  0
	//BATCH_REMOVE                 0
	//REQUEST_RESPONSE             0
	//PAGED_RANGE                  0
	//READ_REPAIR                  0

	var metrics = []c.Metric{}

	lines := strings.Split(output, c.NEWLINE)
	state := 0
	for _,line := range lines {
		if (state == 0 || state == 2) {
			state++ // skip the line
		} else if state == 1 {
			if line != c.EMPTY {
				metrics = appendTpPool(metrics, line)
			} else {
				state = 2
			}
		} else if state == 3 {
			if line != c.EMPTY {
				metrics = appendTpMessageType(metrics, line)
			}
		}
	}

	return metrics, nil
}

func appendTpPool(metrics []c.Metric, line string) []c.Metric {
	valuesOnly := strings.Fields(line)
	prefix := "ntTpPool" + valuesOnly[0]
	metrics = append(metrics, *c.NewMetricIntString(c.MT_COUNT, valuesOnly[1], prefix + "Active", c.PROVIDER_NODETOOL))
	metrics = append(metrics, *c.NewMetricIntString(c.MT_COUNT, valuesOnly[2], prefix + "Pending", c.PROVIDER_NODETOOL))
	metrics = append(metrics, *c.NewMetricIntString(c.MT_COUNT, valuesOnly[3], prefix + "Completed", c.PROVIDER_NODETOOL))
	metrics = append(metrics, *c.NewMetricIntString(c.MT_COUNT, valuesOnly[4], prefix + "Blocked", c.PROVIDER_NODETOOL))
	return append(metrics, *c.NewMetricIntString(c.MT_COUNT, valuesOnly[5], prefix + "AllTimeBlocked", c.PROVIDER_NODETOOL))
}

func appendTpMessageType(metrics []c.Metric, line string) []c.Metric {
	valuesOnly := strings.Fields(line)
	parts := strings.Split(valuesOnly[0], c.UNDER)
	name := "ntTpMsgType"
	for _,part := range parts {
		if part != c.EMPTY {
			name = name + c.UpFirst(part)
		}
	}
	return append(metrics, *c.NewMetricIntString(c.MT_COUNT, valuesOnly[1], name, c.PROVIDER_NODETOOL))
}
