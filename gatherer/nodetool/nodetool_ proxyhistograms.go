package nodetool

import (
	"strings"
	c "github.com/cloudurable/metricsd/common"
)

func ProxyHistograms(nodetoolCommand string) ([]c.Metric, error) {
	output, err := c.ExecCommand(nodetoolCommand, NtFunc_proxyhistograms)
	if err != nil {
		return nil, err
	}

	//proxy histograms
	//Percentile       Read Latency      Write Latency      Range Latency   CAS Read Latency  CAS Write Latency View Write Latency
	//                     (micros)           (micros)           (micros)           (micros)           (micros)           (micros)
	//50%                      0.00               0.00               0.00               0.00               0.00               0.00
	//75%                      0.00               0.00               0.00               0.00               0.00               0.00
	//95%                      0.00               0.00               0.00               0.00               0.00               0.00
	//98%                      0.00               0.00               0.00               0.00               0.00               0.00
	//99%                      0.00               0.00               0.00               0.00               0.00               0.00
	//Min                      0.00               0.00               0.00               0.00               0.00               0.00
	//Max                      0.00               0.00               0.00               0.00               0.00               0.00
	//<blank line>

	var metrics = []c.Metric{}

	lines := strings.Split(output, c.NEWLINE)
	for index, line := range lines {
		if index > 2 && line != c.EMPTY {
			valuesOnly := strings.Fields(line)
			prefix := "ntPh" + valuesOnly[0]
			metrics = append(metrics, *c.NewMetricIntString(c.MT_MICROS, valuesOnly[1], prefix + "ReadLatency", c.PROVIDER_NODETOOL))
			metrics = append(metrics, *c.NewMetricIntString(c.MT_MICROS, valuesOnly[2], prefix + "WriteLatency", c.PROVIDER_NODETOOL))
			metrics = append(metrics, *c.NewMetricIntString(c.MT_MICROS, valuesOnly[3], prefix + "RangeLatency", c.PROVIDER_NODETOOL))
			metrics = append(metrics, *c.NewMetricIntString(c.MT_MICROS, valuesOnly[4], prefix + "CASReadLatency", c.PROVIDER_NODETOOL))
			metrics = append(metrics, *c.NewMetricIntString(c.MT_MICROS, valuesOnly[5], prefix + "CASWriteLatency", c.PROVIDER_NODETOOL))
			metrics = append(metrics, *c.NewMetricIntString(c.MT_MICROS, valuesOnly[6], prefix + "ViewWriteLatency", c.PROVIDER_NODETOOL))
		}
	}

	return metrics, nil
}

