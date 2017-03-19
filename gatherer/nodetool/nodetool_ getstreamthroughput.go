package nodetool

import (
	c "github.com/cloudurable/metricsd/common"
	"strings"
)

func GetStreamThroughput(nodetoolCommand string) ([]c.Metric, error) {
	// Current stream throughput: 200 Mb/s
	output, err := c.ExecCommand(nodetoolCommand, NtFunc_getstreamthroughput)
	if err != nil {
		return nil, err
	}

	// colonAt := strings.
	var metrics = []c.Metric{}
	for _,ntfun := range []string{"statusbackup", "statusbinary", "statusgossip", "statushandoff", "statusthrift"} {
		name := "ntStatus" + c.UpFirst(ntfun[6:])
		metrics = append(metrics, *c.NewMetricString(strings.TrimSuffix(output, c.NEWLINE), name, c.PROVIDER_NODETOOL))
	}

	for _,ntfun := range []string{"version"} {
		output, err := c.ExecCommand(nodetoolCommand, ntfun)
		if err != nil {
			return nil, err
		}
		name := "ntStatus" + c.UpFirst(ntfun)
		metrics = append(metrics, *c.NewMetricString(strings.TrimSuffix(output, c.NEWLINE), name, c.PROVIDER_NODETOOL))
	}

	return metrics, nil
}
