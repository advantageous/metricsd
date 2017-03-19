package nodetool

import (
	"strings"
	c "github.com/cloudurable/metricsd/common"
)

func Gettimeout(nodetoolCommand string) ([]c.Metric, error) {
	var metrics = []c.Metric{}

	for _, timeouttype := range []string{"read", "range", "write", "counterwrite", "cascontention", "truncate", "streamingsocket", "misc"} {
		output, err := c.ExecCommand(nodetoolCommand, NtFunc_gettimeout, timeouttype)
		if err != nil {
			return nil, err
		}

		// Current timeout for type aaaaa: 1000 ms
		lines := strings.Split(output, c.NEWLINE)
		temp := strings.Fields(lines[0])
		ix := len(temp) - 2
		metrics = append(metrics, *c.NewMetricIntString(c.MT_MILLIS, temp[ix], "ntTo" + c.UpFirst(timeouttype), c.PROVIDER_NODETOOL))
	}

	return metrics, nil
}

