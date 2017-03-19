package nodetool

import (
	c "github.com/cloudurable/metricsd/common"
	"strings"
)

func Statuses(nodetoolCommand string) ([]c.Metric, error) {
	//statusbackup                 Status of incremental backup
	//statusbinary                 Status of native transport (binary protocol)
	//statusgossip                 Status of gossip
	//statushandoff                Status of storing future hints on the current node
	//statusthrift                 Status of thrift server
	//version I DID THIS HERE INSTEAD OF HAVING ANOTHER WHOLE TOOL

	var metrics = []c.Metric{}
	for _,ntfun := range []string{"statusbackup", "statusbinary", "statusgossip", "statushandoff", "statusthrift"} {
		output, err := c.ExecCommand(nodetoolCommand, ntfun)
		if err != nil {
			return nil, err
		}
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
